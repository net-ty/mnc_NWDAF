package context

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sort"

	"github.com/free5gc/smf/factory"
	"github.com/free5gc/smf/logger"
)

type UEDefaultPaths struct {
	AnchorUPFs      []string // list of UPF name
	DefaultPathPool DefaultPathPool
}

type DefaultPathPool map[string]*DataPath // key: UPF name

func NewUEDefaultPaths(upi *UserPlaneInformation, topology []factory.UPLink) (*UEDefaultPaths, error) {
	defaultPathPool := make(map[string]*DataPath)
	source, err := findSourceInTopology(upi, topology)
	if err != nil {
		return nil, err
	}
	destinations, err := extractAnchorUPFForULCL(upi, source, topology)
	if err != nil {
		return nil, err
	}
	for _, destination := range destinations {
		path, err := generateDefaultDataPath(source, destination, topology)
		if err != nil {
			return nil, err
		}
		defaultPathPool[destination] = path
	}
	defautlPaths := &UEDefaultPaths{
		AnchorUPFs:      destinations,
		DefaultPathPool: defaultPathPool,
	}
	return defautlPaths, nil
}

func findSourceInTopology(upi *UserPlaneInformation, topology []factory.UPLink) (string, error) {
	sourceList := make([]string, 0)
	for key, node := range upi.AccessNetwork {
		if node.Type == UPNODE_AN {
			sourceList = append(sourceList, key)
		}
	}
	for _, anName := range sourceList {
		for _, link := range topology {
			if link.A == anName || link.B == anName {
				// if multiple gNBs exist, select one according to some criterion
				logger.InitLog.Debugf("%s is AN", anName)
				return anName, nil
			}
		}
	}
	return "", errors.New("Not found AN node in topology")
}

func extractAnchorUPFForULCL(upi *UserPlaneInformation, source string, topology []factory.UPLink) ([]string, error) {
	upList := make([]string, 0)
	visited := make(map[string]bool)
	queue := make([]string, 0)

	queue = append(queue, source)
	queued := make(map[string]bool)
	queued[source] = true

	for {
		node := queue[0]
		queue = queue[1:]
		findNewLink := false
		for _, link := range topology {
			if link.A == node {
				if !queued[link.B] {
					queue = append(queue, link.B)
					queued[link.B] = true
					findNewLink = true
				}
				if !visited[link.B] {
					findNewLink = true
				}
			}
			if link.B == node {
				if !queued[link.A] {
					queue = append(queue, link.A)
					queued[link.A] = true
					findNewLink = true
				}
				if !visited[link.A] {
					findNewLink = true
				}
			}
		}
		visited[node] = true
		if !findNewLink {
			logger.InitLog.Debugf("%s is Anchor UPF", node)
			upList = append(upList, node)
		}
		if len(queue) == 0 {
			break
		}
	}
	if len(upList) == 0 {
		return nil, errors.New("Not found Anchor UPF in topology")
	}
	sort.Strings(upList)
	return upList, nil
}

func generateDefaultDataPath(source string, destination string, topology []factory.UPLink) (*DataPath, error) {
	visited := make(map[string]bool)
	currentPath := []string{}
	allPaths := make(map[int][]string)
	count := 0
	getPathBetweenByNodeName(source, destination, topology, visited, currentPath, allPaths, &count)
	if len(allPaths) == 0 {
		return nil, fmt.Errorf("Path not exist: %s to %s", source, destination)
	}

	dataPath := NewDataPath()
	lowerBound := 0
	var parentNode *DataPathNode = nil

	// if multiple Paths exist, select one according to some criterion
	for idx, nodeName := range allPaths[0] {
		newUeNode, err := NewUEDataPathNode(nodeName)
		if err != nil {
			return nil, err
		}
		if idx == lowerBound {
			dataPath.FirstDPNode = newUeNode
		}
		if parentNode != nil {
			newUeNode.AddPrev(parentNode)
			parentNode.AddNext(newUeNode)
		}
		parentNode = newUeNode
	}
	logger.CtxLog.Tracef("New default data path (%s to %s): ", source, destination)
	logger.CtxLog.Traceln("\n" + dataPath.String() + "\n")
	return dataPath, nil
}

func getPathBetweenByNodeName(src string, dest string, links []factory.UPLink, visited map[string]bool,
	currentPath []string, allPaths map[int][]string, count *int) {
	if visited[src] {
		return
	}
	visited[src] = true
	currentPath = append(currentPath, src)
	logger.InitLog.Traceln("current path:", currentPath)
	if src == dest {
		allPaths[*count] = currentPath[1:]
		(*count)++
		logger.InitLog.Traceln("all path:", allPaths)
		visited[src] = false
		currentPath = []string{}
		return
	}
	for _, link := range links {
		// search A to B only
		if link.A == src {
			getPathBetweenByNodeName(link.B, dest, links, visited, currentPath, allPaths, count)
		}
	}
	visited[src] = false
	currentPath = []string{}
}

func createUPFListForSelectionULCL(inputList []string) (outputList []string) {
	offset := rand.Intn(len(inputList))
	return append(inputList[offset:], inputList[:offset]...)
}

func (dfp *UEDefaultPaths) SelectUPFAndAllocUEIPForULCL(upi *UserPlaneInformation, selection *UPFSelectionParams) (string, net.IP) {

	sortedUPFList := createUPFListForSelectionULCL(dfp.AnchorUPFs)

	for _, upfName := range sortedUPFList {
		logger.CtxLog.Debugf("check start UPF: %s", upfName)
		upf := upi.UPFs[upfName]

		pools := getUEIPPool(upf, selection)
		if pools == nil || len(pools) == 0 {
			continue
		}
		sortedPoolList := createPoolListForSelection(pools)
		for _, pool := range sortedPoolList {
			logger.CtxLog.Debugf("check start UEIPPool(%+v)", pool.ueSubNet)
			addr := pool.allocate()
			if addr != nil {
				logger.CtxLog.Infof("Selected UPF: %s", upfName)
				return upfName, addr
			}
			// if all addresses in pool are used, search next pool
			logger.CtxLog.Debug("check next pool")
		}
		// if all addresses in UPF are used, search next UPF
		logger.CtxLog.Debug("check next upf")
	}
	// checked all UPFs
	logger.CtxLog.Warnf("UE IP pool exhausted for DNN[%s] S-NSSAI[sst: %d sd: %s] DNAI[%s]\n", selection.Dnn,
		selection.SNssai.Sst, selection.SNssai.Sd, selection.Dnai)
	return "", nil
}

func (dfp *UEDefaultPaths) GetDefaultPath(upfName string) *DataPath {
	firstNode := dfp.DefaultPathPool[upfName].CopyFirstDPNode()
	dataPath := &DataPath{
		Activated:     false,
		IsDefaultPath: true,
		Destination:   dfp.DefaultPathPool[upfName].Destination,
		FirstDPNode:   firstNode,
	}
	return dataPath
}
