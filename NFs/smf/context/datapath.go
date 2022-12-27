package context

import (
	"fmt"
	"strconv"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/logger"
	"github.com/free5gc/smf/util"
	"github.com/free5gc/util_3gpp"
)

// GTPTunnel represents the GTP tunnel information
type GTPTunnel struct {
	SrcEndPoint  *DataPathNode
	DestEndPoint *DataPathNode

	TEID uint32
	PDR  *PDR
}

type DataPathNode struct {
	UPF *UPF
	// DataPathToAN *DataPathDownLink
	// DataPathToDN map[string]*DataPathUpLink //uuid to DataPathLink

	UpLinkTunnel   *GTPTunnel
	DownLinkTunnel *GTPTunnel
	// for UE Routing Topology
	// for special case:
	// branching & leafnode

	// InUse                bool
	IsBranchingPoint bool
	// DLDataPathLinkForPSA *DataPathUpLink
	// BPUpLinkPDRs         map[string]*DataPathDownLink // uuid to UpLink
}

type DataPath struct {
	// meta data
	Activated         bool
	IsDefaultPath     bool
	Destination       Destination
	HasBranchingPoint bool
	// Data Path Double Link List
	FirstDPNode *DataPathNode
}

type DataPathPool map[int64]*DataPath

type Destination struct {
	DestinationIP   string
	DestinationPort string
	Url             string
}

func NewDataPathNode() *DataPathNode {
	node := &DataPathNode{
		UpLinkTunnel:   &GTPTunnel{},
		DownLinkTunnel: &GTPTunnel{},
	}
	return node
}

func NewDataPath() *DataPath {
	dataPath := &DataPath{
		Destination: Destination{
			DestinationIP:   "",
			DestinationPort: "",
			Url:             "",
		},
	}

	return dataPath
}

func NewDataPathPool() DataPathPool {
	pool := make(map[int64]*DataPath)
	return pool
}

func (node *DataPathNode) AddNext(next *DataPathNode) {
	node.DownLinkTunnel.SrcEndPoint = next
}

func (node *DataPathNode) AddPrev(prev *DataPathNode) {
	node.UpLinkTunnel.SrcEndPoint = prev
}

func (node *DataPathNode) Next() *DataPathNode {
	if node.DownLinkTunnel == nil {
		return nil
	}
	next := node.DownLinkTunnel.SrcEndPoint
	return next
}

func (node *DataPathNode) Prev() *DataPathNode {
	if node.UpLinkTunnel == nil {
		return nil
	}
	prev := node.UpLinkTunnel.SrcEndPoint
	return prev
}

func (node *DataPathNode) ActivateUpLinkTunnel(smContext *SMContext) error {
	var err error
	logger.CtxLog.Traceln("In ActivateUpLinkTunnel")
	node.UpLinkTunnel.SrcEndPoint = node.Prev()
	node.UpLinkTunnel.DestEndPoint = node

	destUPF := node.UPF
	if node.UpLinkTunnel.PDR, err = destUPF.AddPDR(); err != nil {
		logger.CtxLog.Errorln("In ActivateUpLinkTunnel UPF IP: ", node.UPF.NodeID.ResolveNodeIdToIp().String())
		logger.CtxLog.Errorln("Allocate PDR Error: ", err)
		return fmt.Errorf("Add PDR failed: %s", err)
	}

	if err = smContext.PutPDRtoPFCPSession(destUPF.NodeID, node.UpLinkTunnel.PDR); err != nil {
		logger.CtxLog.Errorln("Put PDR Error: ", err)
		return err
	}

	if teid, err := destUPF.GenerateTEID(); err != nil {
		logger.CtxLog.Errorf("Generate uplink TEID fail: %s", err)
		return err
	} else {
		node.UpLinkTunnel.TEID = teid
	}

	return nil
}

func (node *DataPathNode) ActivateDownLinkTunnel(smContext *SMContext) error {
	var err error
	node.DownLinkTunnel.SrcEndPoint = node.Next()
	node.DownLinkTunnel.DestEndPoint = node

	destUPF := node.UPF
	if node.DownLinkTunnel.PDR, err = destUPF.AddPDR(); err != nil {
		logger.CtxLog.Errorln("In ActivateDownLinkTunnel UPF IP: ", node.UPF.NodeID.ResolveNodeIdToIp().String())
		logger.CtxLog.Errorln("Allocate PDR Error: ", err)
		return fmt.Errorf("Add PDR failed: %s", err)
	}

	if err = smContext.PutPDRtoPFCPSession(destUPF.NodeID, node.DownLinkTunnel.PDR); err != nil {
		logger.CtxLog.Errorln("Put PDR Error: ", err)
		return err
	}

	if teid, err := destUPF.GenerateTEID(); err != nil {
		logger.CtxLog.Errorf("Generate downlink TEID fail: %s", err)
		return err
	} else {
		node.DownLinkTunnel.TEID = teid
	}

	return nil
}

func (node *DataPathNode) DeactivateUpLinkTunnel(smContext *SMContext) {
	if pdr := node.UpLinkTunnel.PDR; pdr != nil {
		smContext.RemovePDRfromPFCPSession(node.UPF.NodeID, pdr)
		err := node.UPF.RemovePDR(pdr)
		if err != nil {
			logger.CtxLog.Warnln("Deactivaed UpLinkTunnel", err)
		}

		if far := pdr.FAR; far != nil {
			err = node.UPF.RemoveFAR(far)
			if err != nil {
				logger.CtxLog.Warnln("Deactivaed UpLinkTunnel", err)
			}

			bar := far.BAR
			if bar != nil {
				err = node.UPF.RemoveBAR(bar)
				if err != nil {
					logger.CtxLog.Warnln("Deactivaed UpLinkTunnel", err)
				}
			}
		}
		if qerList := pdr.QER; qerList != nil {
			for _, qer := range qerList {
				if qer != nil {
					err = node.UPF.RemoveQER(qer)
					if err != nil {
						logger.CtxLog.Warnln("Deactivaed UpLinkTunnel", err)
					}
				}
			}
		}
	}

	teid := node.DownLinkTunnel.TEID
	node.UPF.teidGenerator.FreeID(int64(teid))
	node.DownLinkTunnel = &GTPTunnel{}
}

func (node *DataPathNode) DeactivateDownLinkTunnel(smContext *SMContext) {
	if pdr := node.DownLinkTunnel.PDR; pdr != nil {
		smContext.RemovePDRfromPFCPSession(node.UPF.NodeID, pdr)
		err := node.UPF.RemovePDR(pdr)
		if err != nil {
			logger.CtxLog.Warnln("Deactivaed DownLinkTunnel", err)
		}

		if far := pdr.FAR; far != nil {
			err = node.UPF.RemoveFAR(far)
			if err != nil {
				logger.CtxLog.Warnln("Deactivaed DownLinkTunnel", err)
			}

			bar := far.BAR
			if bar != nil {
				err = node.UPF.RemoveBAR(bar)
				if err != nil {
					logger.CtxLog.Warnln("Deactivaed DownLinkTunnel", err)
				}
			}
		}
		if qerList := pdr.QER; qerList != nil {
			for _, qer := range qerList {
				if qer != nil {
					err = node.UPF.RemoveQER(qer)
					if err != nil {
						logger.CtxLog.Warnln("Deactivaed UpLinkTunnel", err)
					}
				}
			}
		}
	}

	teid := node.DownLinkTunnel.TEID
	node.UPF.teidGenerator.FreeID(int64(teid))
	node.DownLinkTunnel = &GTPTunnel{}
}

func (node *DataPathNode) GetUPFID() (id string, err error) {
	node_ip := node.GetNodeIP()
	var exist bool

	if id, exist = smfContext.UserPlaneInformation.UPFsIPtoID[node_ip]; !exist {
		err = fmt.Errorf("UPNode IP %s doesn't exist in smfcfg.yaml", node_ip)
		return "", err
	}

	return id, nil
}

func (node *DataPathNode) GetNodeIP() (ip string) {
	ip = node.UPF.NodeID.ResolveNodeIdToIp().String()
	return
}

func (node *DataPathNode) IsANUPF() bool {
	if node.Prev() == nil {
		return true
	} else {
		return false
	}
}

func (node *DataPathNode) IsAnchorUPF() bool {
	if node.Next() == nil {
		return true
	} else {
		return false
	}
}

func (node *DataPathNode) GetUpLinkPDR() (pdr *PDR) {
	return node.UpLinkTunnel.PDR
}

func (node *DataPathNode) GetUpLinkFAR() (far *FAR) {
	return node.UpLinkTunnel.PDR.FAR
}

func (dataPathPool DataPathPool) GetDefaultPath() (dataPath *DataPath) {
	for _, path := range dataPathPool {
		if path.IsDefaultPath {
			dataPath = path
			return
		}
	}
	return
}

func (dataPath *DataPath) String() string {
	firstDPNode := dataPath.FirstDPNode

	var str string

	str += "DataPath Meta Information\n"
	str += "Activated: " + strconv.FormatBool(dataPath.Activated) + "\n"
	str += "IsDefault Path: " + strconv.FormatBool(dataPath.IsDefaultPath) + "\n"
	str += "Has Braching Point: " + strconv.FormatBool(dataPath.HasBranchingPoint) + "\n"
	str += "Destination IP: " + dataPath.Destination.DestinationIP + "\n"
	str += "Destination Port: " + dataPath.Destination.DestinationPort + "\n"

	str += "DataPath Routing Information\n"
	index := 1
	for curDPNode := firstDPNode; curDPNode != nil; curDPNode = curDPNode.Next() {
		str += strconv.Itoa(index) + "th Node in the Path\n"
		str += "Current UPF IP: " + curDPNode.GetNodeIP() + "\n"
		if curDPNode.Prev() != nil {
			str += "Previous UPF IP: " + curDPNode.Prev().GetNodeIP() + "\n"
		} else {
			str += "Previous UPF IP: None\n"
		}

		if curDPNode.Next() != nil {
			str += "Next UPF IP: " + curDPNode.Next().GetNodeIP() + "\n"
		} else {
			str += "Next UPF IP: None\n"
		}

		index++
	}

	return str
}

func (dataPath *DataPath) ActivateTunnelAndPDR(smContext *SMContext, precedence uint32) {
	smContext.AllocateLocalSEIDForDataPath(dataPath)

	firstDPNode := dataPath.FirstDPNode
	logger.PduSessLog.Traceln("In ActivateTunnelAndPDR")
	logger.PduSessLog.Traceln(dataPath.String())
	// Activate Tunnels
	for curDataPathNode := firstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		logger.PduSessLog.Traceln("Current DP Node IP: ", curDataPathNode.UPF.NodeID.ResolveNodeIdToIp().String())
		if err := curDataPathNode.ActivateUpLinkTunnel(smContext); err != nil {
			logger.CtxLog.Warnln(err)
			return
		}
		if err := curDataPathNode.ActivateDownLinkTunnel(smContext); err != nil {
			logger.CtxLog.Warnln(err)
			return
		}
	}

	sessionRule := smContext.SelectedSessionRule()
	AuthDefQos := sessionRule.AuthDefQos

	// Activate PDR
	for curDataPathNode := firstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		var flowQER *QER

		// if has the sess QoS (QER == 1), this value **SHOULD BE** uint32
		if oldQER, ok := curDataPathNode.UPF.qerPool.Load(uint32(1)); ok {
			flowQER = oldQER.(*QER)
		} else {
			if newQER, err := curDataPathNode.UPF.AddQER(); err != nil {
				logger.PduSessLog.Errorln("new QER failed")
				return
			} else {
				newQER.QFI.QFI = uint8(AuthDefQos.Var5qi)
				newQER.GateStatus = &pfcpType.GateStatus{
					ULGate: pfcpType.GateOpen,
					DLGate: pfcpType.GateOpen,
				}
				newQER.MBR = &pfcpType.MBR{
					ULMBR: util.BitRateTokbps(sessionRule.AuthSessAmbr.Uplink),
					DLMBR: util.BitRateTokbps(sessionRule.AuthSessAmbr.Downlink),
				}

				flowQER = newQER
			}
		}

		logger.CtxLog.Traceln("Calculate ", curDataPathNode.UPF.PFCPAddr().String())
		curULTunnel := curDataPathNode.UpLinkTunnel
		curDLTunnel := curDataPathNode.DownLinkTunnel

		// Setup UpLink PDR
		if curULTunnel != nil {
			ULPDR := curULTunnel.PDR
			ULDestUPF := curULTunnel.DestEndPoint.UPF
			ULPDR.QER = append(ULPDR.QER, flowQER)

			ULPDR.Precedence = precedence

			var iface *UPFInterfaceInfo
			if curDataPathNode.IsANUPF() {
				iface = ULDestUPF.GetInterface(models.UpInterfaceType_N3, smContext.Dnn)
			} else {
				iface = ULDestUPF.GetInterface(models.UpInterfaceType_N9, smContext.Dnn)
			}

			if upIP, err := iface.IP(smContext.SelectedPDUSessionType); err != nil {
				logger.CtxLog.Errorln("ActivateTunnelAndPDR failed", err)
				return
			} else {
				ULPDR.PDI = PDI{
					SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceAccess},
					LocalFTeid: &pfcpType.FTEID{
						V4:          true,
						Ipv4Address: upIP,
						Teid:        curULTunnel.TEID,
					},
					UEIPAddress: &pfcpType.UEIPAddress{
						V4:          true,
						Ipv4Address: smContext.PDUAddress.To4(),
					},
				}
			}

			ULPDR.OuterHeaderRemoval = &pfcpType.OuterHeaderRemoval{
				OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4,
			}

			ULFAR := ULPDR.FAR
			ULFAR.ApplyAction = pfcpType.ApplyAction{
				Buff: false,
				Drop: false,
				Dupl: false,
				Forw: true,
				Nocp: false,
			}
			ULFAR.ForwardingParameters = &ForwardingParameters{
				DestinationInterface: pfcpType.DestinationInterface{
					InterfaceValue: pfcpType.DestinationInterfaceCore,
				},
				NetworkInstance: []byte(smContext.Dnn),
			}

			if curDataPathNode.IsAnchorUPF() {
				ULFAR.ForwardingParameters.
					DestinationInterface.InterfaceValue = pfcpType.DestinationInterfaceSgiLanN6Lan
			}

			if nextULDest := curDataPathNode.Next(); nextULDest != nil {
				nextULTunnel := nextULDest.UpLinkTunnel
				iface = nextULTunnel.DestEndPoint.UPF.GetInterface(models.UpInterfaceType_N9, smContext.Dnn)

				if upIP, err := iface.IP(smContext.SelectedPDUSessionType); err != nil {
					logger.CtxLog.Errorln("ActivateTunnelAndPDR failed", err)
					return
				} else {
					ULFAR.ForwardingParameters.OuterHeaderCreation = &pfcpType.OuterHeaderCreation{
						OuterHeaderCreationDescription: pfcpType.OuterHeaderCreationGtpUUdpIpv4,
						Ipv4Address:                    upIP,
						Teid:                           nextULTunnel.TEID,
					}
				}
			}
		}

		// Setup DownLink
		if curDLTunnel != nil {
			var iface *UPFInterfaceInfo
			DLPDR := curDLTunnel.PDR
			DLDestUPF := curDLTunnel.DestEndPoint.UPF
			DLPDR.QER = append(DLPDR.QER, flowQER)

			DLPDR.Precedence = precedence

			// TODO: Should delete this after FR5GC-1029 is solved
			if curDataPathNode.IsAnchorUPF() {
				DLPDR.PDI.UEIPAddress = &pfcpType.UEIPAddress{
					V4:          true,
					Ipv4Address: smContext.PDUAddress.To4(),
				}
			} else {
				DLPDR.OuterHeaderRemoval = &pfcpType.OuterHeaderRemoval{
					OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4,
				}

				iface = DLDestUPF.GetInterface(models.UpInterfaceType_N9, smContext.Dnn)
				if upIP, err := iface.IP(smContext.SelectedPDUSessionType); err != nil {
					logger.CtxLog.Errorln("ActivateTunnelAndPDR failed", err)
					return
				} else {
					DLPDR.PDI = PDI{
						SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceCore},
						LocalFTeid: &pfcpType.FTEID{
							V4:          true,
							Ipv4Address: upIP,
							Teid:        curDLTunnel.TEID,
						},

						// TODO: Should Uncomment this after FR5GC-1029 is solved
						// UEIPAddress: &pfcpType.UEIPAddress{
						// 	V4:          true,
						// 	Ipv4Address: smContext.PDUAddress.To4(),
						// },
					}
				}
			}

			DLFAR := DLPDR.FAR

			logger.PduSessLog.Traceln("Current DP Node IP: ", curDataPathNode.UPF.NodeID.ResolveNodeIdToIp().String())
			logger.PduSessLog.Traceln("Before DLPDR OuterHeaderCreation")
			if nextDLDest := curDataPathNode.Prev(); nextDLDest != nil {
				logger.PduSessLog.Traceln("In DLPDR OuterHeaderCreation")
				nextDLTunnel := nextDLDest.DownLinkTunnel

				DLFAR.ApplyAction = pfcpType.ApplyAction{
					Buff: false,
					Drop: false,
					Dupl: false,
					Forw: true,
					Nocp: false,
				}

				iface = nextDLDest.UPF.GetInterface(models.UpInterfaceType_N9, smContext.Dnn)

				if upIP, err := iface.IP(smContext.SelectedPDUSessionType); err != nil {
					logger.CtxLog.Errorln("ActivateTunnelAndPDR failed", err)
					return
				} else {
					DLFAR.ForwardingParameters = &ForwardingParameters{
						DestinationInterface: pfcpType.DestinationInterface{InterfaceValue: pfcpType.DestinationInterfaceAccess},
						OuterHeaderCreation: &pfcpType.OuterHeaderCreation{
							OuterHeaderCreationDescription: pfcpType.OuterHeaderCreationGtpUUdpIpv4,
							Ipv4Address:                    upIP,
							Teid:                           nextDLTunnel.TEID,
						},
					}
				}
			} else {
				if anIP := smContext.Tunnel.ANInformation.IPAddress; anIP != nil {
					ANUPF := dataPath.FirstDPNode
					DLPDR := ANUPF.DownLinkTunnel.PDR
					DLFAR := DLPDR.FAR
					DLFAR.ForwardingParameters = new(ForwardingParameters)
					DLFAR.ForwardingParameters.DestinationInterface.InterfaceValue = pfcpType.DestinationInterfaceAccess
					DLFAR.ForwardingParameters.NetworkInstance = []byte(smContext.Dnn)
					DLFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)

					dlOuterHeaderCreation := DLFAR.ForwardingParameters.OuterHeaderCreation
					dlOuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
					dlOuterHeaderCreation.Teid = smContext.Tunnel.ANInformation.TEID
					dlOuterHeaderCreation.Ipv4Address = smContext.Tunnel.ANInformation.IPAddress.To4()
				}
			}
		}
		if curDataPathNode.DownLinkTunnel != nil {
			if curDataPathNode.DownLinkTunnel.SrcEndPoint == nil {
				DNDLPDR := curDataPathNode.DownLinkTunnel.PDR
				DNDLPDR.PDI = PDI{
					SourceInterface: pfcpType.SourceInterface{InterfaceValue: pfcpType.SourceInterfaceCore},
					NetworkInstance: util_3gpp.Dnn(smContext.Dnn),
					UEIPAddress: &pfcpType.UEIPAddress{
						V4:          true,
						Ipv4Address: smContext.PDUAddress.To4(),
					},
				}
			}
		}
	}

	dataPath.Activated = true
}

func (dataPath *DataPath) DeactivateTunnelAndPDR(smContext *SMContext) {
	firstDPNode := dataPath.FirstDPNode

	// Deactivate Tunnels
	for curDataPathNode := firstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		curDataPathNode.DeactivateUpLinkTunnel(smContext)
		curDataPathNode.DeactivateDownLinkTunnel(smContext)
	}

	dataPath.Activated = false
}

func (dataPath *DataPath) CopyFirstDPNode() *DataPathNode {
	if dataPath.FirstDPNode == nil {
		return nil
	}
	var firstNode *DataPathNode = nil
	var parentNode *DataPathNode = nil
	for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		newNode := NewDataPathNode()
		if firstNode == nil {
			firstNode = newNode
		}
		newNode.UPF = curDataPathNode.UPF
		if parentNode != nil {
			newNode.AddPrev(parentNode)
			parentNode.AddNext(newNode)
		}
		parentNode = newNode
	}
	return firstNode
}
