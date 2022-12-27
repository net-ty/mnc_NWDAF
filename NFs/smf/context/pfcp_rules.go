package context

import (
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/util_3gpp"
)

const (
	RULE_INITIAL RuleState = 0
	RULE_CREATE  RuleState = 1
	RULE_UPDATE  RuleState = 2
	RULE_REMOVE  RuleState = 3
)

type RuleState uint8

// Packet Detection Rule. Table 7.5.2.2-1
type PDR struct {
	PDRID uint16

	Precedence         uint32
	PDI                PDI
	OuterHeaderRemoval *pfcpType.OuterHeaderRemoval

	FAR *FAR
	URR *URR
	QER []*QER

	State RuleState
}

// Packet Detection. 7.5.2.2-2
type PDI struct {
	SourceInterface pfcpType.SourceInterface
	LocalFTeid      *pfcpType.FTEID
	NetworkInstance util_3gpp.Dnn
	UEIPAddress     *pfcpType.UEIPAddress
	SDFFilter       *pfcpType.SDFFilter
	ApplicationID   string
}

// Forwarding Action Rule. 7.5.2.3-1
type FAR struct {
	FARID uint32

	ApplyAction          pfcpType.ApplyAction
	ForwardingParameters *ForwardingParameters

	BAR   *BAR
	State RuleState
}

// Forwarding Parameters. 7.5.2.3-2
type ForwardingParameters struct {
	DestinationInterface pfcpType.DestinationInterface
	NetworkInstance      util_3gpp.Dnn
	OuterHeaderCreation  *pfcpType.OuterHeaderCreation
	ForwardingPolicyID   string
}

// Buffering Action Rule 7.5.2.6-1
type BAR struct {
	BARID uint8

	DownlinkDataNotificationDelay  pfcpType.DownlinkDataNotificationDelay
	SuggestedBufferingPacketsCount pfcpType.SuggestedBufferingPacketsCount

	State RuleState
}

// QoS Enhancement Rule
type QER struct {
	QERID uint32

	QFI pfcpType.QFI

	GateStatus *pfcpType.GateStatus
	MBR        *pfcpType.MBR
	GBR        *pfcpType.GBR

	State RuleState
}

// Usage Report Rule
type URR struct {
}
