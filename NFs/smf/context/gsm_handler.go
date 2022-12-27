package context

import (
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/smf/logger"
)

func (smContext *SMContext) HandlePDUSessionEstablishmentRequest(req *nasMessage.PDUSessionEstablishmentRequest) {
	// Retrieve PDUSessionID
	smContext.PDUSessionID = int32(req.PDUSessionID.GetPDUSessionID())
	logger.GsmLog.Infoln("In HandlePDUSessionEstablishmentRequest")

	// Retrieve PTI (Procedure transaction identity)
	smContext.Pti = req.GetPTI()

	// Handle PDUSessionType
	if req.PDUSessionType != nil {
		requestedPDUSessionType := req.PDUSessionType.GetPDUSessionTypeValue()
		if err := smContext.isAllowedPDUSessionType(requestedPDUSessionType); err != nil {
			logger.CtxLog.Errorf("%s", err)
			return
		}
	} else {
		// Set to default supported PDU Session Type
		switch SMF_Self().SupportedPDUSessionType {
		case "IPv4":
			smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4
		case "IPv6":
			smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv6
		case "IPv4v6":
			smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4IPv6
		case "Ethernet":
			smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeEthernet
		default:
			smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4
		}
	}

	if req.ExtendedProtocolConfigurationOptions != nil {
		EPCOContents := req.ExtendedProtocolConfigurationOptions.GetExtendedProtocolConfigurationOptionsContents()
		protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
		unmarshalErr := protocolConfigurationOptions.UnMarshal(EPCOContents)
		if unmarshalErr != nil {
			logger.GsmLog.Errorf("Parsing PCO failed: %s", unmarshalErr)
		}
		logger.GsmLog.Infoln("Protocol Configuration Options")
		logger.GsmLog.Infoln(protocolConfigurationOptions)

		for _, container := range protocolConfigurationOptions.ProtocolOrContainerList {
			logger.GsmLog.Traceln("Container ID: ", container.ProtocolOrContainerID)
			logger.GsmLog.Traceln("Container Length: ", container.LengthOfContents)
			switch container.ProtocolOrContainerID {
			case nasMessage.PCSCFIPv6AddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type PCSCFIPv6AddressRequestUL")
			case nasMessage.IMCNSubsystemSignalingFlagUL:
				logger.GsmLog.Infoln("Didn't Implement container type IMCNSubsystemSignalingFlagUL")
			case nasMessage.DNSServerIPv6AddressRequestUL:
				smContext.ProtocolConfigurationOptions.DNSIPv6Request = true
			case nasMessage.NotSupportedUL:
				logger.GsmLog.Infoln("Didn't Implement container type NotSupportedUL")
			case nasMessage.MSSupportOfNetworkRequestedBearerControlIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSSupportOfNetworkRequestedBearerControlIndicatorUL")
			case nasMessage.DSMIPv6HomeAgentAddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6HomeAgentAddressRequestUL")
			case nasMessage.DSMIPv6HomeNetworkPrefixRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6HomeNetworkPrefixRequestUL")
			case nasMessage.DSMIPv6IPv4HomeAgentAddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6IPv4HomeAgentAddressRequestUL")
			case nasMessage.IPAddressAllocationViaNASSignallingUL:
				logger.GsmLog.Infoln("Didn't Implement container type IPAddressAllocationViaNASSignallingUL")
			case nasMessage.IPv4AddressAllocationViaDHCPv4UL:
				logger.GsmLog.Infoln("Didn't Implement container type IPv4AddressAllocationViaDHCPv4UL")
			case nasMessage.PCSCFIPv4AddressRequestUL:
				smContext.ProtocolConfigurationOptions.PCSCFIPv4Request = true
			case nasMessage.DNSServerIPv4AddressRequestUL:
				smContext.ProtocolConfigurationOptions.DNSIPv4Request = true
			case nasMessage.MSISDNRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSISDNRequestUL")
			case nasMessage.IFOMSupportRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type IFOMSupportRequestUL")
			case nasMessage.IPv4LinkMTURequestUL:
				smContext.ProtocolConfigurationOptions.IPv4LinkMTURequest = true
			case nasMessage.MSSupportOfLocalAddressInTFTIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSSupportOfLocalAddressInTFTIndicatorUL")
			case nasMessage.PCSCFReSelectionSupportUL:
				logger.GsmLog.Infoln("Didn't Implement container type PCSCFReSelectionSupportUL")
			case nasMessage.NBIFOMRequestIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type NBIFOMRequestIndicatorUL")
			case nasMessage.NBIFOMModeUL:
				logger.GsmLog.Infoln("Didn't Implement container type NBIFOMModeUL")
			case nasMessage.NonIPLinkMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type NonIPLinkMTURequestUL")
			case nasMessage.APNRateControlSupportIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type APNRateControlSupportIndicatorUL")
			case nasMessage.UEStatus3GPPPSDataOffUL:
				logger.GsmLog.Infoln("Didn't Implement container type UEStatus3GPPPSDataOffUL")
			case nasMessage.ReliableDataServiceRequestIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type ReliableDataServiceRequestIndicatorUL")
			case nasMessage.AdditionalAPNRateControlForExceptionDataSupportIndicatorUL:
				logger.GsmLog.Infoln(
					"Didn't Implement container type AdditionalAPNRateControlForExceptionDataSupportIndicatorUL",
				)
			case nasMessage.PDUSessionIDUL:
				logger.GsmLog.Infoln("Didn't Implement container type PDUSessionIDUL")
			case nasMessage.EthernetFramePayloadMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type EthernetFramePayloadMTURequestUL")
			case nasMessage.UnstructuredLinkMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type UnstructuredLinkMTURequestUL")
			case nasMessage.I5GSMCauseValueUL:
				logger.GsmLog.Infoln("Didn't Implement container type 5GSMCauseValueUL")
			case nasMessage.QoSRulesWithTheLengthOfTwoOctetsSupportIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type QoSRulesWithTheLengthOfTwoOctetsSupportIndicatorUL")
			case nasMessage.QoSFlowDescriptionsWithTheLengthOfTwoOctetsSupportIndicatorUL:
				logger.GsmLog.Infoln(
					"Didn't Implement container type QoSFlowDescriptionsWithTheLengthOfTwoOctetsSupportIndicatorUL",
				)
			case nasMessage.LinkControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type LinkControlProtocolUL")
			case nasMessage.PushAccessControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type PushAccessControlProtocolUL")
			case nasMessage.ChallengeHandshakeAuthenticationProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type ChallengeHandshakeAuthenticationProtocolUL")
			case nasMessage.InternetProtocolControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type InternetProtocolControlProtocolUL")
			default:
				logger.GsmLog.Infof("Unknown Container ID [%d]", container.ProtocolOrContainerID)
			}
		}
	}
}

func (smContext *SMContext) HandlePDUSessionReleaseRequest(req *nasMessage.PDUSessionReleaseRequest) {
	logger.GsmLog.Infof("Handle Pdu Session Release Request")

	// Retrieve PTI (Procedure transaction identity)
	smContext.Pti = req.GetPTI()
}
