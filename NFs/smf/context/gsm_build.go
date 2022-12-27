package context

import (
	"encoding/hex"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/smf/logger"
)

func BuildGSMPDUSessionEstablishmentAccept(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionEstablishmentAccept = nasMessage.NewPDUSessionEstablishmentAccept(0x0)
	pDUSessionEstablishmentAccept := m.PDUSessionEstablishmentAccept

	sessRule := smContext.SelectedSessionRule()
	authDefQos := sessRule.AuthDefQos

	pDUSessionEstablishmentAccept.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionEstablishmentAccept.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	pDUSessionEstablishmentAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionEstablishmentAccept.SetPTI(smContext.Pti)

	if v := smContext.EstAcceptCause5gSMValue; v != 0 {
		pDUSessionEstablishmentAccept.Cause5GSM = nasType.NewCause5GSM(nasMessage.PDUSessionEstablishmentAcceptCause5GSMType)
		pDUSessionEstablishmentAccept.Cause5GSM.SetCauseValue(v)
	}
	pDUSessionEstablishmentAccept.SetPDUSessionType(smContext.SelectedPDUSessionType)

	pDUSessionEstablishmentAccept.SetSSCMode(1)
	pDUSessionEstablishmentAccept.SessionAMBR = nasConvert.ModelsToSessionAMBR(sessRule.AuthSessAmbr)
	pDUSessionEstablishmentAccept.SessionAMBR.SetLen(uint8(len(pDUSessionEstablishmentAccept.SessionAMBR.Octet)))

	qoSRules := QoSRules{
		QoSRule{
			Identifier:    0x01,
			DQR:           0x01,
			OperationCode: OperationCodeCreateNewQoSRule,
			Precedence:    0xff,
			QFI:           uint8(authDefQos.Var5qi),
			PacketFilterList: []PacketFilter{
				{
					Identifier:    0x01,
					Direction:     PacketFilterDirectionBidirectional,
					ComponentType: PacketFilterComponentTypeMatchAll,
				},
			},
		},
	}

	qosRulesBytes, err := qoSRules.MarshalBinary()
	if err != nil {
		return nil, err
	}

	pDUSessionEstablishmentAccept.AuthorizedQosRules.SetLen(uint16(len(qosRulesBytes)))
	pDUSessionEstablishmentAccept.AuthorizedQosRules.SetQosRule(qosRulesBytes)

	if smContext.PDUAddress != nil {
		addr, addrLen := smContext.PDUAddressToNAS()
		pDUSessionEstablishmentAccept.PDUAddress =
			nasType.NewPDUAddress(nasMessage.PDUSessionEstablishmentAcceptPDUAddressType)
		pDUSessionEstablishmentAccept.PDUAddress.SetLen(addrLen)
		pDUSessionEstablishmentAccept.PDUAddress.SetPDUSessionTypeValue(smContext.SelectedPDUSessionType)
		pDUSessionEstablishmentAccept.PDUAddress.SetPDUAddressInformation(addr)
	}

	pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions =
		nasType.NewAuthorizedQosFlowDescriptions(nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType)
	pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions.SetLen(6)
	pDUSessionEstablishmentAccept.SetQoSFlowDescriptions([]uint8{uint8(authDefQos.Var5qi), 0x20, 0x41, 0x01, 0x01, 0x09})

	var sd [3]uint8

	if byteArray, err := hex.DecodeString(smContext.Snssai.Sd); err != nil {
		return nil, err
	} else {
		copy(sd[:], byteArray)
	}

	pDUSessionEstablishmentAccept.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
	pDUSessionEstablishmentAccept.SNSSAI.SetLen(4)
	pDUSessionEstablishmentAccept.SNSSAI.SetSST(uint8(smContext.Snssai.Sst))
	pDUSessionEstablishmentAccept.SNSSAI.SetSD(sd)

	dnn := []byte(smContext.Dnn)
	pDUSessionEstablishmentAccept.DNN = nasType.NewDNN(nasMessage.ULNASTransportDNNType)
	pDUSessionEstablishmentAccept.DNN.SetLen(uint8(len(dnn)))
	pDUSessionEstablishmentAccept.DNN.SetDNN(dnn)

	if smContext.ProtocolConfigurationOptions.DNSIPv4Request ||
		smContext.ProtocolConfigurationOptions.DNSIPv6Request ||
		smContext.ProtocolConfigurationOptions.PCSCFIPv4Request ||
		smContext.ProtocolConfigurationOptions.IPv4LinkMTURequest {
		pDUSessionEstablishmentAccept.ExtendedProtocolConfigurationOptions =
			nasType.NewExtendedProtocolConfigurationOptions(
				nasMessage.PDUSessionEstablishmentAcceptExtendedProtocolConfigurationOptionsType,
			)
		protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()

		// IPv4 DNS
		if smContext.ProtocolConfigurationOptions.DNSIPv4Request {
			err := protocolConfigurationOptions.AddDNSServerIPv4Address(smContext.DNNInfo.DNS.IPv4Addr)
			if err != nil {
				logger.GsmLog.Warnln("Error while adding DNS IPv4 Addr: ", err)
			}
		}

		// IPv6 DNS
		if smContext.ProtocolConfigurationOptions.DNSIPv6Request {
			err := protocolConfigurationOptions.AddDNSServerIPv6Address(smContext.DNNInfo.DNS.IPv6Addr)
			if err != nil {
				logger.GsmLog.Warnln("Error while adding DNS IPv6 Addr: ", err)
			}
		}

		// IPv4 PCSCF (need for ims DNN)
		if smContext.ProtocolConfigurationOptions.PCSCFIPv4Request {
			err := protocolConfigurationOptions.AddPCSCFIPv4Address(smContext.DNNInfo.PCSCF.IPv4Addr)
			if err != nil {
				logger.GsmLog.Warnln("Error while adding PCSCF IPv4 Addr: ", err)
			}
		}

		// MTU
		if smContext.ProtocolConfigurationOptions.IPv4LinkMTURequest {
			err := protocolConfigurationOptions.AddIPv4LinkMTU(1400)
			if err != nil {
				logger.GsmLog.Warnln("Error while adding MTU: ", err)
			}
		}

		pcoContents := protocolConfigurationOptions.Marshal()
		pcoContentsLength := len(pcoContents)
		pDUSessionEstablishmentAccept.
			ExtendedProtocolConfigurationOptions.
			SetLen(uint16(pcoContentsLength))
		pDUSessionEstablishmentAccept.
			ExtendedProtocolConfigurationOptions.
			SetExtendedProtocolConfigurationOptionsContents(pcoContents)
	}
	return m.PlainNasEncode()
}

func BuildGSMPDUSessionEstablishmentReject(smContext *SMContext, cause uint8) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentReject)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionEstablishmentReject = nasMessage.NewPDUSessionEstablishmentReject(0x0)
	pDUSessionEstablishmentReject := m.PDUSessionEstablishmentReject

	pDUSessionEstablishmentReject.SetMessageType(nas.MsgTypePDUSessionEstablishmentReject)
	pDUSessionEstablishmentReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionEstablishmentReject.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionEstablishmentReject.SetCauseValue(cause)

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionReleaseCommand(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionReleaseCommand = nasMessage.NewPDUSessionReleaseCommand(0x0)
	pDUSessionReleaseCommand := m.PDUSessionReleaseCommand

	pDUSessionReleaseCommand.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	pDUSessionReleaseCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionReleaseCommand.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionReleaseCommand.SetPTI(smContext.Pti)
	pDUSessionReleaseCommand.SetCauseValue(0x0)

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionModificationCommand(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationCommand)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionModificationCommand = nasMessage.NewPDUSessionModificationCommand(0x0)
	pDUSessionModificationCommand := m.PDUSessionModificationCommand

	pDUSessionModificationCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionModificationCommand.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionModificationCommand.SetPTI(smContext.Pti)
	pDUSessionModificationCommand.SetMessageType(nas.MsgTypePDUSessionModificationCommand)
	// pDUSessionModificationCommand.SetQosRule()
	// pDUSessionModificationCommand.AuthorizedQosRules.SetLen()
	// pDUSessionModificationCommand.SessionAMBR.SetSessionAMBRForDownlink([2]uint8{0x11, 0x11})
	// pDUSessionModificationCommand.SessionAMBR.SetSessionAMBRForUplink([2]uint8{0x11, 0x11})
	// pDUSessionModificationCommand.SessionAMBR.SetUnitForSessionAMBRForDownlink(10)
	// pDUSessionModificationCommand.SessionAMBR.SetUnitForSessionAMBRForUplink(10)
	// pDUSessionModificationCommand.SessionAMBR.SetLen(uint8(len(pDUSessionModificationCommand.SessionAMBR.Octet)))

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionReleaseReject(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseReject)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionReleaseReject = nasMessage.NewPDUSessionReleaseReject(0x0)
	pDUSessionReleaseReject := m.PDUSessionReleaseReject

	pDUSessionReleaseReject.SetMessageType(nas.MsgTypePDUSessionReleaseReject)
	pDUSessionReleaseReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)

	pDUSessionReleaseReject.SetPDUSessionID(uint8(smContext.PDUSessionID))

	pDUSessionReleaseReject.SetPTI(smContext.Pti)
	// TODO: fix to real value
	pDUSessionReleaseReject.SetCauseValue(nasMessage.Cause5GSMRequestRejectedUnspecified)

	return m.PlainNasEncode()
}
