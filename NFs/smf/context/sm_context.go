package context

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Namf_Communication"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Npcf_SMPolicyControl"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/logger"
)

var (
	smContextPool    sync.Map
	canonicalRef     sync.Map
	seidSMContextMap sync.Map
)

var smContextCount uint64

type SMContextState int

const (
	InActive SMContextState = iota
	ActivePending
	Active
	InActivePending
	ModificationPending
	PFCPModification
)

func init() {
}

func GetSMContextCount() uint64 {
	atomic.AddUint64(&smContextCount, 1)
	return smContextCount
}

type SMContext struct {
	Ref string

	LocalSEID  uint64
	RemoteSEID uint64

	UnauthenticatedSupi bool
	// SUPI or PEI
	Supi           string
	Pei            string
	Identifier     string
	Gpsi           string
	PDUSessionID   int32
	Dnn            string
	Snssai         *models.Snssai
	HplmnSnssai    *models.Snssai
	ServingNetwork *models.PlmnId
	ServingNfId    string

	UpCnxState models.UpCnxState

	AnType          models.AccessType
	RatType         models.RatType
	PresenceInLadn  models.PresenceState
	UeLocation      *models.UserLocation
	UeTimeZone      string
	AddUeLocation   *models.UserLocation
	OldPduSessionId int32
	HoState         models.HoState

	PDUAddress             net.IP
	SelectedPDUSessionType uint8

	DnnConfiguration models.DnnConfiguration

	// Client
	SMPolicyClient      *Npcf_SMPolicyControl.APIClient
	CommunicationClient *Namf_Communication.APIClient

	AMFProfile         models.NfProfile
	SelectedPCFProfile models.NfProfile
	SmStatusNotifyUri  string

	SMContextState SMContextState

	Tunnel      *UPTunnel
	SelectedUPF *UPNode
	BPManager   *BPManager
	// NodeID(string form) to PFCP Session Context
	PFCPContext                         map[string]*PFCPSessionContext
	SBIPFCPCommunicationChan            chan PFCPSessionResponseStatus
	PendingUPF                          PendingUPF
	PDUSessionRelease_DUE_TO_DUP_PDU_ID bool

	DNNInfo *SnssaiSmfDnnInfo

	// SM Policy related
	PCCRules           map[string]*PCCRule
	SessionRules       map[string]*SessionRule
	TrafficControlPool map[string]*TrafficControlData

	// NAS
	Pti                     uint8
	EstAcceptCause5gSMValue uint8

	// PCO Related
	ProtocolConfigurationOptions *ProtocolConfigurationOptions

	// lock
	SMLock sync.Mutex
}

func canonicalName(identifier string, pduSessID int32) (canonical string) {
	return fmt.Sprintf("%s-%d", identifier, pduSessID)
}

func ResolveRef(identifier string, pduSessID int32) (ref string, err error) {
	if value, ok := canonicalRef.Load(canonicalName(identifier, pduSessID)); ok {
		ref = value.(string)
		err = nil
	} else {
		ref = ""
		err = fmt.Errorf(
			"UE '%s' - PDUSessionID '%d' not found in SMContext", identifier, pduSessID)
	}
	return
}

func NewSMContext(identifier string, pduSessID int32) (smContext *SMContext) {
	smContext = new(SMContext)
	// Create Ref and identifier
	smContext.Ref = uuid.New().URN()
	smContextPool.Store(smContext.Ref, smContext)
	canonicalRef.Store(canonicalName(identifier, pduSessID), smContext.Ref)

	smContext.SMContextState = InActive
	smContext.Identifier = identifier
	smContext.PDUSessionID = pduSessID
	smContext.PFCPContext = make(map[string]*PFCPSessionContext)
	smContext.LocalSEID = GetSMContextCount()

	// initialize SM Policy Data
	smContext.PCCRules = make(map[string]*PCCRule)
	smContext.SessionRules = make(map[string]*SessionRule)
	smContext.TrafficControlPool = make(map[string]*TrafficControlData)
	smContext.SBIPFCPCommunicationChan = make(chan PFCPSessionResponseStatus, 1)

	smContext.ProtocolConfigurationOptions = &ProtocolConfigurationOptions{}

	return smContext
}

//*** add unit test ***//
func GetSMContext(ref string) (smContext *SMContext) {
	if value, ok := smContextPool.Load(ref); ok {
		smContext = value.(*SMContext)
	}

	return
}

//*** add unit test ***//
func RemoveSMContext(ref string) {
	var smContext *SMContext
	if value, ok := smContextPool.Load(ref); ok {
		smContext = value.(*SMContext)
	}

	if smContext.SelectedUPF != nil {
		logger.PduSessLog.Infof("UE[%s] PDUSessionID[%d] Release IP[%s]",
			smContext.Supi, smContext.PDUSessionID, smContext.PDUAddress.String())
		GetUserPlaneInformation().ReleaseUEIP(smContext.SelectedUPF, smContext.PDUAddress)
	}

	for _, pfcpSessionContext := range smContext.PFCPContext {
		seidSMContextMap.Delete(pfcpSessionContext.LocalSEID)
	}

	smContextPool.Delete(ref)
}

//*** add unit test ***//
func GetSMContextBySEID(SEID uint64) (smContext *SMContext) {
	if value, ok := seidSMContextMap.Load(SEID); ok {
		smContext = value.(*SMContext)
	}
	return
}

//*** add unit test ***//
func (smContext *SMContext) SetCreateData(createData *models.SmContextCreateData) {
	smContext.Gpsi = createData.Gpsi
	smContext.Supi = createData.Supi
	smContext.Dnn = createData.Dnn
	smContext.Snssai = createData.SNssai
	smContext.HplmnSnssai = createData.HplmnSnssai
	smContext.ServingNetwork = createData.ServingNetwork
	smContext.AnType = createData.AnType
	smContext.RatType = createData.RatType
	smContext.PresenceInLadn = createData.PresenceInLadn
	smContext.UeLocation = createData.UeLocation
	smContext.UeTimeZone = createData.UeTimeZone
	smContext.AddUeLocation = createData.AddUeLocation
	smContext.OldPduSessionId = createData.OldPduSessionId
	smContext.ServingNfId = createData.ServingNfId
}

func (smContext *SMContext) BuildCreatedData() (createdData *models.SmContextCreatedData) {
	createdData = new(models.SmContextCreatedData)
	createdData.SNssai = smContext.Snssai
	return
}

func (smContext *SMContext) PDUAddressToNAS() (addr [12]byte, addrLen uint8) {
	copy(addr[:], smContext.PDUAddress)
	switch smContext.SelectedPDUSessionType {
	case nasMessage.PDUSessionTypeIPv4:
		addrLen = 4 + 1
	case nasMessage.PDUSessionTypeIPv6:
	case nasMessage.PDUSessionTypeIPv4IPv6:
		addrLen = 12 + 1
	}
	return
}

// PCFSelection will select PCF for this SM Context
func (smContext *SMContext) PCFSelection() error {
	// Send NFDiscovery for find PCF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	rep, res, err := SMF_Self().
		NFDiscoveryClient.
		NFInstancesStoreApi.
		SearchNFInstances(context.TODO(), models.NfType_PCF, models.NfType_SMF, &localVarOptionals)
	if err != nil {
		return err
	}
	defer func() {
		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
			logger.PduSessLog.Errorf("SmfEventExposureNotification response body cannot close: %+v", rspCloseErr)
		}
	}()

	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			apiError := err.(openapi.GenericOpenAPIError)
			problemDetails := apiError.Model().(models.ProblemDetails)

			logger.CtxLog.Warningf("NFDiscovery PCF return status: %d\n", status)
			logger.CtxLog.Warningf("Detail: %v\n", problemDetails.Title)
		}
	}

	// Select PCF from available PCF

	smContext.SelectedPCFProfile = rep.NfInstances[0]

	// Create SMPolicyControl Client for this SM Context
	for _, service := range *smContext.SelectedPCFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NPCF_SMPOLICYCONTROL {
			SmPolicyControlConf := Npcf_SMPolicyControl.NewConfiguration()
			SmPolicyControlConf.SetBasePath(service.ApiPrefix)
			smContext.SMPolicyClient = Npcf_SMPolicyControl.NewAPIClient(SmPolicyControlConf)
		}
	}

	return nil
}

func (smContext *SMContext) GetNodeIDByLocalSEID(seid uint64) (nodeID pfcpType.NodeID) {
	for _, pfcpCtx := range smContext.PFCPContext {
		if pfcpCtx.LocalSEID == seid {
			nodeID = pfcpCtx.NodeID
		}
	}

	return
}

func (smContext *SMContext) AllocateLocalSEIDForUPPath(path UPPath) {
	for _, upNode := range path {
		NodeIDtoIP := upNode.NodeID.ResolveNodeIdToIp().String()
		if _, exist := smContext.PFCPContext[NodeIDtoIP]; !exist {
			allocatedSEID := AllocateLocalSEID()

			smContext.PFCPContext[NodeIDtoIP] = &PFCPSessionContext{
				PDRs:      make(map[uint16]*PDR),
				NodeID:    upNode.NodeID,
				LocalSEID: allocatedSEID,
			}

			seidSMContextMap.Store(allocatedSEID, smContext)
		}
	}
}

func (smContext *SMContext) AllocateLocalSEIDForDataPath(dataPath *DataPath) {
	logger.PduSessLog.Traceln("In AllocateLocalSEIDForDataPath")
	for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
		NodeIDtoIP := curDataPathNode.UPF.NodeID.ResolveNodeIdToIp().String()
		logger.PduSessLog.Traceln("NodeIDtoIP: ", NodeIDtoIP)
		if _, exist := smContext.PFCPContext[NodeIDtoIP]; !exist {
			allocatedSEID := AllocateLocalSEID()
			smContext.PFCPContext[NodeIDtoIP] = &PFCPSessionContext{
				PDRs:      make(map[uint16]*PDR),
				NodeID:    curDataPathNode.UPF.NodeID,
				LocalSEID: allocatedSEID,
			}

			seidSMContextMap.Store(allocatedSEID, smContext)
		}
	}
}

func (smContext *SMContext) PutPDRtoPFCPSession(nodeID pfcpType.NodeID, pdr *PDR) error {
	NodeIDtoIP := nodeID.ResolveNodeIdToIp().String()
	if pfcpSessCtx, exist := smContext.PFCPContext[NodeIDtoIP]; exist {
		pfcpSessCtx.PDRs[pdr.PDRID] = pdr
	} else {
		return fmt.Errorf("Can't find PFCPContext[%s] to put PDR(%d)", NodeIDtoIP, pdr.PDRID)
	}
	return nil
}

func (smContext *SMContext) RemovePDRfromPFCPSession(nodeID pfcpType.NodeID, pdr *PDR) {
	NodeIDtoIP := nodeID.ResolveNodeIdToIp().String()
	pfcpSessCtx := smContext.PFCPContext[NodeIDtoIP]
	delete(pfcpSessCtx.PDRs, pdr.PDRID)
}

func (smContext *SMContext) isAllowedPDUSessionType(requestedPDUSessionType uint8) error {
	dnnPDUSessionType := smContext.DnnConfiguration.PduSessionTypes
	if dnnPDUSessionType == nil {
		return fmt.Errorf("this SMContext[%s] has no subscription pdu session type info", smContext.Ref)
	}

	allowIPv4 := false
	allowIPv6 := false
	allowEthernet := false

	for _, allowedPDUSessionType := range smContext.DnnConfiguration.PduSessionTypes.AllowedSessionTypes {
		switch allowedPDUSessionType {
		case models.PduSessionType_IPV4:
			allowIPv4 = true
		case models.PduSessionType_IPV6:
			allowIPv6 = true
		case models.PduSessionType_IPV4_V6:
			allowIPv4 = true
			allowIPv6 = true
		case models.PduSessionType_ETHERNET:
			allowEthernet = true
		}
	}

	supportedPDUSessionType := SMF_Self().SupportedPDUSessionType
	switch supportedPDUSessionType {
	case "IPv4":
		if !allowIPv4 {
			return fmt.Errorf("No SupportedPDUSessionType[%q] in DNN[%s] configuration", supportedPDUSessionType, smContext.Dnn)
		}
	case "IPv6":
		if !allowIPv6 {
			return fmt.Errorf("No SupportedPDUSessionType[%q] in DNN[%s] configuration", supportedPDUSessionType, smContext.Dnn)
		}
	case "IPv4v6":
		if !allowIPv4 && !allowIPv6 {
			return fmt.Errorf("No SupportedPDUSessionType[%q] in DNN[%s] configuration", supportedPDUSessionType, smContext.Dnn)
		}
	case "Ethernet":
		if !allowEthernet {
			return fmt.Errorf("No SupportedPDUSessionType[%q] in DNN[%s] configuration", supportedPDUSessionType, smContext.Dnn)
		}
	}

	smContext.EstAcceptCause5gSMValue = 0
	switch nasConvert.PDUSessionTypeToModels(requestedPDUSessionType) {
	case models.PduSessionType_IPV4:
		if allowIPv4 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4)
		} else {
			return fmt.Errorf("PduSessionType_IPV4 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_IPV6:
		if allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV6)
		} else {
			return fmt.Errorf("PduSessionType_IPV6 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_IPV4_V6:
		if allowIPv4 && allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4_V6)
		} else if allowIPv4 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4)
			smContext.EstAcceptCause5gSMValue = nasMessage.Cause5GSMPDUSessionTypeIPv4OnlyAllowed
		} else if allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV6)
			smContext.EstAcceptCause5gSMValue = nasMessage.Cause5GSMPDUSessionTypeIPv6OnlyAllowed
		} else {
			return fmt.Errorf("PduSessionType_IPV4_V6 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_ETHERNET:
		if allowEthernet {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_ETHERNET)
		} else {
			return fmt.Errorf("PduSessionType_ETHERNET is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	default:
		return fmt.Errorf("Requested PDU Sesstion type[%d] is not supported", requestedPDUSessionType)
	}
	return nil
}

// SM Policy related operation

// SelectedSessionRule - return the SMF selected session rule for this SM Context
func (smContext *SMContext) SelectedSessionRule() *SessionRule {
	for _, sessionRule := range smContext.SessionRules {
		if sessionRule.isActivate {
			return sessionRule
		}
	}

	return nil
}

func (smContextState SMContextState) String() string {
	switch smContextState {
	case InActive:
		return "InActive"
	case ActivePending:
		return "ActivePending"
	case Active:
		return "Active"
	case InActivePending:
		return "InActivePending"
	case ModificationPending:
		return "ModificationPending"
	case PFCPModification:
		return "PFCPModification"
	default:
		return "Unknown State"
	}
}
