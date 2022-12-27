package context

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"reflect"
	"sort"

	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/context/pool"
	"github.com/free5gc/smf/factory"
	"github.com/free5gc/smf/logger"
)

// UserPlaneInformation store userplane topology
type UserPlaneInformation struct {
	UPNodes                   map[string]*UPNode
	UPFs                      map[string]*UPNode
	AccessNetwork             map[string]*UPNode
	UPFIPToName               map[string]string
	UPFsID                    map[string]string               // name to id
	UPFsIPtoID                map[string]string               // ip->id table, for speed optimization
	DefaultUserPlanePath      map[string][]*UPNode            // DNN to Default Path
	DefaultUserPlanePathToUPF map[string]map[string][]*UPNode // DNN and UPF to Default Path
}

type UPNodeType string

const (
	UPNODE_UPF UPNodeType = "UPF"
	UPNODE_AN  UPNodeType = "AN"
)

// UPNode represent the user plane node topology
type UPNode struct {
	Type   UPNodeType
	NodeID pfcpType.NodeID
	ANIP   net.IP
	Dnn    string
	Links  []*UPNode
	UPF    *UPF
}

// UPPath represent User Plane Sequence of this path
type UPPath []*UPNode

func AllocateUPFID() {
	UPFsID := smfContext.UserPlaneInformation.UPFsID
	UPFsIPtoID := smfContext.UserPlaneInformation.UPFsIPtoID

	for upfName, upfNode := range smfContext.UserPlaneInformation.UPFs {
		upfid := upfNode.UPF.UUID()
		upfip := upfNode.NodeID.ResolveNodeIdToIp().String()

		UPFsID[upfName] = upfid
		UPFsIPtoID[upfip] = upfid
	}
}

// NewUserPlaneInformation process the configuration then returns a new instance of UserPlaneInformation
func NewUserPlaneInformation(upTopology *factory.UserPlaneInformation) *UserPlaneInformation {
	nodePool := make(map[string]*UPNode)
	upfPool := make(map[string]*UPNode)
	anPool := make(map[string]*UPNode)
	upfIPMap := make(map[string]string)
	allUEIPPools := []*UeIPPool{}

	for name, node := range upTopology.UPNodes {
		upNode := new(UPNode)
		upNode.Type = UPNodeType(node.Type)
		switch upNode.Type {
		case UPNODE_AN:
			upNode.ANIP = net.ParseIP(node.ANIP)
			anPool[name] = upNode
		case UPNODE_UPF:
			// ParseIp() always return 16 bytes
			// so we can't use the length of return ip to separate IPv4 and IPv6
			// This is just a work around
			var ip net.IP
			if net.ParseIP(node.NodeID).To4() == nil {
				ip = net.ParseIP(node.NodeID)
			} else {
				ip = net.ParseIP(node.NodeID).To4()
			}

			switch len(ip) {
			case net.IPv4len:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
					NodeIdValue: ip,
				}
			case net.IPv6len:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeIpv6Address,
					NodeIdValue: ip,
				}
			default:
				upNode.NodeID = pfcpType.NodeID{
					NodeIdType:  pfcpType.NodeIdTypeFqdn,
					NodeIdValue: []byte(node.NodeID),
				}
			}

			upNode.UPF = NewUPF(&upNode.NodeID, node.InterfaceUpfInfoList)
			snssaiInfos := make([]SnssaiUPFInfo, 0)
			for _, snssaiInfoConfig := range node.SNssaiInfos {
				snssaiInfo := SnssaiUPFInfo{
					SNssai: SNssai{
						Sst: snssaiInfoConfig.SNssai.Sst,
						Sd:  snssaiInfoConfig.SNssai.Sd,
					},
					DnnList: make([]DnnUPFInfoItem, 0),
				}

				for _, dnnInfoConfig := range snssaiInfoConfig.DnnUpfInfoList {
					ueIPPools := make([]*UeIPPool, 0)
					for _, pool := range dnnInfoConfig.Pools {
						ueIPPool := NewUEIPPool(&pool)
						if ueIPPool == nil {
							logger.InitLog.Fatalf("invalid pools value: %+v", pool)
						} else {
							ueIPPools = append(ueIPPools, ueIPPool)
							allUEIPPools = append(allUEIPPools, ueIPPool)
						}
					}
					snssaiInfo.DnnList = append(snssaiInfo.DnnList, DnnUPFInfoItem{
						Dnn:             dnnInfoConfig.Dnn,
						DnaiList:        dnnInfoConfig.DnaiList,
						PduSessionTypes: dnnInfoConfig.PduSessionTypes,
						UeIPPools:       ueIPPools,
					})
				}
				snssaiInfos = append(snssaiInfos, snssaiInfo)
			}
			upNode.UPF.SNssaiInfos = snssaiInfos
			upfPool[name] = upNode
		default:
			logger.InitLog.Warningf("invalid UPNodeType: %s\n", upNode.Type)
		}

		nodePool[name] = upNode

		ipStr := upNode.NodeID.ResolveNodeIdToIp().String()
		upfIPMap[ipStr] = name
	}

	if isOverlap(allUEIPPools) {
		logger.InitLog.Fatalf("overlap cidr value between UPFs")
	}

	for _, link := range upTopology.Links {
		nodeA := nodePool[link.A]
		nodeB := nodePool[link.B]
		if nodeA == nil || nodeB == nil {
			logger.InitLog.Warningf("UPLink [%s] <=> [%s] not establish\n", link.A, link.B)
			continue
		}
		nodeA.Links = append(nodeA.Links, nodeB)
		nodeB.Links = append(nodeB.Links, nodeA)
	}

	userplaneInformation := &UserPlaneInformation{
		UPNodes:                   nodePool,
		UPFs:                      upfPool,
		AccessNetwork:             anPool,
		UPFIPToName:               upfIPMap,
		UPFsID:                    make(map[string]string),
		UPFsIPtoID:                make(map[string]string),
		DefaultUserPlanePath:      make(map[string][]*UPNode),
		DefaultUserPlanePathToUPF: make(map[string]map[string][]*UPNode),
	}

	return userplaneInformation
}

func NewUEIPPool(factoryPool *factory.UEIPPool) *UeIPPool {
	_, ipNet, err := net.ParseCIDR(factoryPool.Cidr)
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	minAddr, maxAddr, err := calcAddrRange(ipNet)
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	newPool, err := pool.NewLazyReusePool(int(minAddr), int(maxAddr))
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	ueIPPool := &UeIPPool{
		ueSubNet: ipNet,
		pool:     newPool,
	}
	return ueIPPool
}

func calcAddrRange(ipNet *net.IPNet) (minAddr, maxAddr uint32, err error) {
	maskVal := binary.BigEndian.Uint32(ipNet.Mask)
	baseIPVal := binary.BigEndian.Uint32(ipNet.IP)
	if maskVal == math.MaxUint32 {
		return baseIPVal, baseIPVal, nil
	}
	minAddr = (baseIPVal & maskVal) + 1  // 0 is network address
	maxAddr = (baseIPVal | ^maskVal) - 1 // all 1 is broadcast address
	if minAddr > maxAddr {
		return minAddr, maxAddr, errors.New("Mask is invalid.")
	}
	return minAddr, maxAddr, nil
}

func isOverlap(pools []*UeIPPool) bool {
	if len(pools) < 2 {
		// no need to check
		return false
	}
	for i := 0; i < len(pools)-1; i++ {
		for j := i + 1; j < len(pools); j++ {
			if pools[i].pool.IsJoint(pools[j].pool) {
				return true
			}
		}
	}
	return false
}

func (upi *UserPlaneInformation) GetUPFNameByIp(ip string) string {
	return upi.UPFIPToName[ip]
}

func (upi *UserPlaneInformation) GetUPFNodeIDByName(name string) pfcpType.NodeID {
	return upi.UPFs[name].NodeID
}

func (upi *UserPlaneInformation) GetUPFNodeByIP(ip string) *UPNode {
	upfName := upi.GetUPFNameByIp(ip)
	return upi.UPFs[upfName]
}

func (upi *UserPlaneInformation) GetUPFIDByIP(ip string) string {
	return upi.UPFsIPtoID[ip]
}

func (upi *UserPlaneInformation) GetDefaultUserPlanePathByDNN(selection *UPFSelectionParams) (path UPPath) {
	path, pathExist := upi.DefaultUserPlanePath[selection.String()]
	logger.CtxLog.Traceln("In GetDefaultUserPlanePathByDNN")
	logger.CtxLog.Traceln("selection: ", selection.String())
	if pathExist {
		return
	} else {
		pathExist = upi.GenerateDefaultPath(selection)
		if pathExist {
			return upi.DefaultUserPlanePath[selection.String()]
		}
	}
	return nil
}

func (upi *UserPlaneInformation) GetDefaultUserPlanePathByDNNAndUPF(selection *UPFSelectionParams, upf *UPNode) (path UPPath) {
	nodeID := upf.NodeID.ResolveNodeIdToIp().String()
	var pathExist bool

	if upi.DefaultUserPlanePathToUPF[selection.String()] != nil {
		path, pathExist := upi.DefaultUserPlanePathToUPF[selection.String()][nodeID]
		logger.CtxLog.Traceln("In GetDefaultUserPlanePathByDNN")
		logger.CtxLog.Traceln("selection: ", selection.String())
		if pathExist {
			return path
		}
	}
	pathExist = upi.GenerateDefaultPathToUPF(selection, upf)
	if pathExist {
		return upi.DefaultUserPlanePathToUPF[selection.String()][nodeID]
	}
	return nil
}

func (upi *UserPlaneInformation) ExistDefaultPath(dnn string) bool {
	_, exist := upi.DefaultUserPlanePath[dnn]
	return exist
}

func GenerateDataPath(upPath UPPath, smContext *SMContext) *DataPath {
	if len(upPath) < 1 {
		logger.CtxLog.Errorf("Invalid data path")
		return nil
	}
	lowerBound := 0
	upperBound := len(upPath) - 1
	var root *DataPathNode
	var curDataPathNode *DataPathNode
	var prevDataPathNode *DataPathNode

	for idx, upNode := range upPath {
		curDataPathNode = NewDataPathNode()
		curDataPathNode.UPF = upNode.UPF

		if idx == lowerBound {
			root = curDataPathNode
			root.AddPrev(nil)
		}
		if idx == upperBound {
			curDataPathNode.AddNext(nil)
		}
		if prevDataPathNode != nil {
			prevDataPathNode.AddNext(curDataPathNode)
			curDataPathNode.AddPrev(prevDataPathNode)
		}
		prevDataPathNode = curDataPathNode
	}

	dataPath := &DataPath{
		Destination: Destination{
			DestinationIP:   "",
			DestinationPort: "",
			Url:             "",
		},
		FirstDPNode: root,
	}
	return dataPath
}

func (upi *UserPlaneInformation) GenerateDefaultPath(selection *UPFSelectionParams) bool {
	var source *UPNode
	var destinations []*UPNode

	for _, node := range upi.AccessNetwork {
		if node.Type == UPNODE_AN {
			source = node
			break
		}
	}

	if source == nil {
		logger.CtxLog.Errorf("There is no AN Node in config file!")
		return false
	}

	destinations = upi.selectMatchUPF(selection)

	if len(destinations) == 0 {
		logger.CtxLog.Errorf("Can't find UPF with DNN[%s] S-NSSAI[sst: %d sd: %s] DNAI[%s]\n", selection.Dnn,
			selection.SNssai.Sst, selection.SNssai.Sd, selection.Dnai)
		return false
	} else {
		logger.CtxLog.Tracef("Find UPF with DNN[%s] S-NSSAI[sst: %d sd: %s] DNAI[%s]\n", selection.Dnn,
			selection.SNssai.Sst, selection.SNssai.Sd, selection.Dnai)
	}

	// Run DFS
	visited := make(map[*UPNode]bool)

	for _, upNode := range upi.UPNodes {
		visited[upNode] = false
	}

	path, pathExist := getPathBetween(source, destinations[0], visited, selection)

	if pathExist {
		if path[0].Type == UPNODE_AN {
			path = path[1:]
		}
		upi.DefaultUserPlanePath[selection.String()] = path
	}

	return pathExist
}

func (upi *UserPlaneInformation) GenerateDefaultPathToUPF(selection *UPFSelectionParams, destination *UPNode) bool {
	var source *UPNode

	for _, node := range upi.AccessNetwork {
		if node.Type == UPNODE_AN {
			source = node
			break
		}
	}

	if source == nil {
		logger.CtxLog.Errorf("There is no AN Node in config file!")
		return false
	}

	// Run DFS
	visited := make(map[*UPNode]bool)

	for _, upNode := range upi.UPNodes {
		visited[upNode] = false
	}

	path, pathExist := getPathBetween(source, destination, visited, selection)

	if pathExist {
		if path[0].Type == UPNODE_AN {
			path = path[1:]
		}
		if upi.DefaultUserPlanePathToUPF[selection.String()] == nil {
			upi.DefaultUserPlanePathToUPF[selection.String()] = make(map[string][]*UPNode)
		}
		upi.DefaultUserPlanePathToUPF[selection.String()][destination.NodeID.ResolveNodeIdToIp().String()] = path
	}

	return pathExist
}

func (upi *UserPlaneInformation) selectMatchUPF(selection *UPFSelectionParams) []*UPNode {
	upList := make([]*UPNode, 0)

	for _, upNode := range upi.UPFs {
		for _, snssaiInfo := range upNode.UPF.SNssaiInfos {
			currentSnssai := &snssaiInfo.SNssai
			targetSnssai := selection.SNssai

			if currentSnssai.Equal(targetSnssai) {
				for _, dnnInfo := range snssaiInfo.DnnList {
					if dnnInfo.Dnn == selection.Dnn && dnnInfo.ContainsDNAI(selection.Dnai) {
						upList = append(upList, upNode)
						break
					}
				}
			}
		}
	}
	return upList
}

func getPathBetween(cur *UPNode, dest *UPNode, visited map[*UPNode]bool,
	selection *UPFSelectionParams) (path []*UPNode, pathExist bool) {
	visited[cur] = true

	if reflect.DeepEqual(*cur, *dest) {
		path = make([]*UPNode, 0)
		path = append(path, cur)
		pathExist = true
		return
	}

	selectedSNssai := selection.SNssai

	for _, nodes := range cur.Links {
		if !visited[nodes] {
			if !nodes.UPF.isSupportSnssai(selectedSNssai) {
				visited[nodes] = true
				continue
			}

			path_tail, path_exist := getPathBetween(nodes, dest, visited, selection)

			if path_exist {
				path = make([]*UPNode, 0)
				path = append(path, cur)

				path = append(path, path_tail...)
				pathExist = true

				return
			}
		}
	}

	return nil, false
}

func (upi *UserPlaneInformation) selectAnchorUPF(source *UPNode, selection *UPFSelectionParams) []*UPNode {
	upList := make([]*UPNode, 0)
	visited := make(map[*UPNode]bool)
	queue := make([]*UPNode, 0)
	targetSnssai := selection.SNssai

	queue = append(queue, source)
	for {
		node := queue[0]
		queue = queue[1:]
		findNewNode := false
		visited[node] = true
		for _, link := range node.Links {
			if !visited[link] {
				for _, snssaiInfo := range link.UPF.SNssaiInfos {
					currentSnssai := &snssaiInfo.SNssai
					if currentSnssai.Equal(targetSnssai) {
						for _, dnnInfo := range snssaiInfo.DnnList {
							if dnnInfo.Dnn == selection.Dnn && dnnInfo.ContainsDNAI(selection.Dnai) {
								queue = append(queue, link)
								findNewNode = true
								break
							}
						}
					}
				}
			}
		}
		if !findNewNode {
			upList = append(upList, node)
		}
		if len(queue) == 0 {
			break
		}
	}
	return upList
}

func (upi *UserPlaneInformation) sortUPFListByName(upfList []*UPNode) []*UPNode {
	keys := make([]string, 0, len(upi.UPFs))
	for k := range upi.UPFs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedUpList := make([]*UPNode, 0)
	for _, name := range keys {
		for _, node := range upfList {
			if name == upi.GetUPFNameByIp(node.NodeID.ResolveNodeIdToIp().String()) {
				sortedUpList = append(sortedUpList, node)
			}
		}
	}
	return sortedUpList
}

func (upi *UserPlaneInformation) selectUPPathSource() (*UPNode, error) {
	// if multiple gNBs exist, select one according to some criterion
	for _, node := range upi.AccessNetwork {
		if node.Type == UPNODE_AN {
			return node, nil
		}
	}
	return nil, errors.New("AN Node not found")
}

func (upi *UserPlaneInformation) SelectUPFAndAllocUEIP(selection *UPFSelectionParams) (*UPNode, net.IP) {
	source, err := upi.selectUPPathSource()
	if err != nil {
		return nil, nil
	}
	UPFList := upi.selectAnchorUPF(source, selection)
	listLength := len(UPFList)
	if listLength == 0 {
		logger.CtxLog.Warnf("Can't find UPF with DNN[%s] S-NSSAI[sst: %d sd: %s] DNAI[%s]\n", selection.Dnn,
			selection.SNssai.Sst, selection.SNssai.Sd, selection.Dnai)
		return nil, nil
	}
	UPFList = upi.sortUPFListByName(UPFList)
	sortedUPFList := createUPFListForSelection(UPFList)
	for _, upf := range sortedUPFList {
		logger.CtxLog.Debugf("check start UPF: %s",
			upi.GetUPFNameByIp(upf.NodeID.ResolveNodeIdToIp().String()))
		pools := getUEIPPool(upf, selection)
		if pools == nil || len(pools) == 0 {
			continue
		}
		sortedPoolList := createPoolListForSelection(pools)
		for _, pool := range sortedPoolList {
			logger.CtxLog.Debugf("check start UEIPPool(%+v)", pool.ueSubNet)
			addr := pool.allocate()
			if addr != nil {
				logger.CtxLog.Infof("Selected UPF: %s",
					upi.GetUPFNameByIp(upf.NodeID.ResolveNodeIdToIp().String()))
				return upf, addr
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
	return nil, nil
}

func createUPFListForSelection(inputList []*UPNode) (outputList []*UPNode) {
	offset := rand.Intn(len(inputList))
	return append(inputList[offset:], inputList[:offset]...)
}

func createPoolListForSelection(inputList []*UeIPPool) (outputList []*UeIPPool) {
	offset := rand.Intn(len(inputList))
	return append(inputList[offset:], inputList[:offset]...)
}

func getUEIPPool(upNode *UPNode, selection *UPFSelectionParams) []*UeIPPool {
	for _, snssaiInfo := range upNode.UPF.SNssaiInfos {
		currentSnssai := &snssaiInfo.SNssai
		targetSnssai := selection.SNssai

		if currentSnssai.Equal(targetSnssai) {
			for _, dnnInfo := range snssaiInfo.DnnList {
				if dnnInfo.Dnn == selection.Dnn && dnnInfo.ContainsDNAI(selection.Dnai) {
					return dnnInfo.UeIPPools
				}
			}
		}
	}
	return nil
}

func (ueIPPool *UeIPPool) allocate() net.IP {
	allocVal, res := ueIPPool.pool.Allocate()
	if !res {
		logger.CtxLog.Warnf("Pool is empty: %+v", ueIPPool.ueSubNet)
		return nil
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(allocVal))
	logger.CtxLog.Infof("Allocated UE IP address: %v", net.IPv4(buf[0], buf[1], buf[2], buf[3]))
	return buf
}

func (upi *UserPlaneInformation) ReleaseUEIP(upf *UPNode, addr net.IP) {
	pool := findPoolByAddr(upf, addr)
	if pool == nil {
		//nothing to do
		logger.CtxLog.Warnf("Fail to release UE IP address: %v to UPF: %s",
			upi.GetUPFNameByIp(upf.NodeID.ResolveNodeIdToIp().String()), addr)
		return
	}
	pool.release(addr)
}

func findPoolByAddr(upf *UPNode, addr net.IP) *UeIPPool {
	for _, snssaiInfo := range upf.UPF.SNssaiInfos {
		for _, dnnInfo := range snssaiInfo.DnnList {
			for _, pool := range dnnInfo.UeIPPools {
				if pool.ueSubNet.Contains(addr) {
					return pool
				}
			}
		}
	}
	return nil
}

func (ueIPPool *UeIPPool) release(addr net.IP) {
	addrVal := binary.BigEndian.Uint32(addr)
	res := ueIPPool.pool.Free(int(addrVal))
	if !res {
		logger.CtxLog.Warnf("failed to release UE Address: %s", addr)
	}
	logger.CtxLog.Debug(ueIPPool.dump())
}

func (ueIPPool *UeIPPool) dump() string {
	str := "["
	elements := ueIPPool.pool.Dump()
	for index, element := range elements {
		var firstAddr net.IP
		var lastAddr net.IP
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[0]))
		firstAddr = buf
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[1]))
		lastAddr = buf
		if index > 0 {
			str += ("->")
		}
		str += fmt.Sprintf("{%s - %s}", firstAddr.String(), lastAddr.String())
	}
	str += ("]")
	return str
}
