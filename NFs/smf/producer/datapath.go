package producer

import (
	"github.com/free5gc/pfcp/pfcpType"
	smf_context "github.com/free5gc/smf/context"
	"github.com/free5gc/smf/logger"
	pfcp_message "github.com/free5gc/smf/pfcp/message"
)

type PFCPState struct {
	nodeID  pfcpType.NodeID
	pdrList []*smf_context.PDR
	farList []*smf_context.FAR
	qerList []*smf_context.QER
}

// SendPFCPRule send one datapath to UPF
func SendPFCPRule(smContext *smf_context.SMContext, dataPath *smf_context.DataPath) {
	logger.PduSessLog.Infoln("Send PFCP Rule")
	logger.PduSessLog.Infoln("DataPath: ", dataPath)
	for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		pdrList := make([]*smf_context.PDR, 0, 2)
		farList := make([]*smf_context.FAR, 0, 2)
		qerList := make([]*smf_context.QER, 0, 2)

		if curDataPathNode.UpLinkTunnel != nil && curDataPathNode.UpLinkTunnel.PDR != nil {
			pdrList = append(pdrList, curDataPathNode.UpLinkTunnel.PDR)
			farList = append(farList, curDataPathNode.UpLinkTunnel.PDR.FAR)
			if curDataPathNode.DownLinkTunnel.PDR.QER != nil {
				qerList = append(qerList, curDataPathNode.DownLinkTunnel.PDR.QER...)
			}
		}
		if curDataPathNode.DownLinkTunnel != nil && curDataPathNode.DownLinkTunnel.PDR != nil {
			pdrList = append(pdrList, curDataPathNode.DownLinkTunnel.PDR)
			farList = append(farList, curDataPathNode.DownLinkTunnel.PDR.FAR)
		}

		sessionContext, exist := smContext.PFCPContext[curDataPathNode.GetNodeIP()]
		if !exist || sessionContext.RemoteSEID == 0 {
			pfcp_message.SendPfcpSessionEstablishmentRequest(
				curDataPathNode.UPF.NodeID, smContext, pdrList, farList, nil, qerList)
		} else {
			pfcp_message.SendPfcpSessionModificationRequest(
				curDataPathNode.UPF.NodeID, smContext, pdrList, farList, nil, qerList)
		}
	}
}

// SendPFCPRules send all datapaths to UPFs
func SendPFCPRules(smContext *smf_context.SMContext) {
	pfcpPool := make(map[string]*PFCPState)

	for _, dataPath := range smContext.Tunnel.DataPathPool {
		if dataPath.Activated {
			for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
				pdrList := make([]*smf_context.PDR, 0, 2)
				farList := make([]*smf_context.FAR, 0, 2)
				qerList := make([]*smf_context.QER, 0, 2)

				if curDataPathNode.UpLinkTunnel != nil && curDataPathNode.UpLinkTunnel.PDR != nil {
					pdrList = append(pdrList, curDataPathNode.UpLinkTunnel.PDR)
					farList = append(farList, curDataPathNode.UpLinkTunnel.PDR.FAR)
					if curDataPathNode.UpLinkTunnel.PDR.QER != nil {
						qerList = append(qerList, curDataPathNode.UpLinkTunnel.PDR.QER...)
					}
				}
				if curDataPathNode.DownLinkTunnel != nil && curDataPathNode.DownLinkTunnel.PDR != nil {
					pdrList = append(pdrList, curDataPathNode.DownLinkTunnel.PDR)
					farList = append(farList, curDataPathNode.DownLinkTunnel.PDR.FAR)
					// skip send QER because uplink and downlink shared one QER
				}

				pfcpState := pfcpPool[curDataPathNode.GetNodeIP()]
				if pfcpState == nil {
					pfcpPool[curDataPathNode.GetNodeIP()] = &PFCPState{
						nodeID:  curDataPathNode.UPF.NodeID,
						pdrList: pdrList,
						farList: farList,
						qerList: qerList,
					}
				} else {
					pfcpState.pdrList = append(pfcpState.pdrList, pdrList...)
					pfcpState.farList = append(pfcpState.farList, farList...)
					pfcpState.qerList = append(pfcpState.qerList, qerList...)
				}
			}
		}
	}
	for ip, pfcp := range pfcpPool {
		sessionContext, exist := smContext.PFCPContext[ip]
		if !exist || sessionContext.RemoteSEID == 0 {
			pfcp_message.SendPfcpSessionEstablishmentRequest(
				pfcp.nodeID, smContext, pfcp.pdrList, pfcp.farList, nil, pfcp.qerList)
		} else {
			pfcp_message.SendPfcpSessionModificationRequest(
				pfcp.nodeID, smContext, pfcp.pdrList, pfcp.farList, nil, pfcp.qerList)
		}
	}
}

func removeDataPath(smContext *smf_context.SMContext, datapath *smf_context.DataPath) {
	for curDPNode := datapath.FirstDPNode; curDPNode != nil; curDPNode = curDPNode.Next() {
		if curDPNode.DownLinkTunnel != nil && curDPNode.DownLinkTunnel.PDR != nil {
			curDPNode.DownLinkTunnel.PDR.State = smf_context.RULE_REMOVE
			curDPNode.DownLinkTunnel.PDR.FAR.State = smf_context.RULE_REMOVE
		}
		if curDPNode.UpLinkTunnel != nil && curDPNode.UpLinkTunnel.PDR != nil {
			curDPNode.UpLinkTunnel.PDR.State = smf_context.RULE_REMOVE
			curDPNode.UpLinkTunnel.PDR.FAR.State = smf_context.RULE_REMOVE
		}
	}
}

// UpdateDataPathToUPF update the datapath of the UPF
func UpdateDataPathToUPF(smContext *smf_context.SMContext, oldDataPath, updateDataPath *smf_context.DataPath) {
	if oldDataPath == nil {
		SendPFCPRule(smContext, updateDataPath)
		return
	} else {
		removeDataPath(smContext, oldDataPath)
		SendPFCPRule(smContext, updateDataPath)
	}
}
