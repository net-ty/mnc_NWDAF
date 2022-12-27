package context

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/UeauCommon"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/fsm"
	"github.com/free5gc/idgenerator"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

type OnGoingProcedure string

const (
	OnGoingProcedureNothing      OnGoingProcedure = "Nothing"
	OnGoingProcedurePaging       OnGoingProcedure = "Paging"
	OnGoingProcedureN2Handover   OnGoingProcedure = "N2Handover"
	OnGoingProcedureRegistration OnGoingProcedure = "Registration"
)

const (
	NgRanCgiPresentNRCGI    int32 = 0
	NgRanCgiPresentEUTRACGI int32 = 1
)

const (
	RecommendRanNodePresentRanNode int32 = 0
	RecommendRanNodePresentTAI     int32 = 1
)

// GMM state for UE
const (
	Deregistered            fsm.StateType = "Deregistered"
	DeregistrationInitiated fsm.StateType = "DeregistrationInitiated"
	Authentication          fsm.StateType = "Authentication"
	SecurityMode            fsm.StateType = "SecurityMode"
	ContextSetup            fsm.StateType = "ContextSetup"
	Registered              fsm.StateType = "Registered"
)

type AmfUe struct {
	/* the AMF which serving this AmfUe now */
	servingAMF *AMFContext // never nil

	/* Gmm State */
	State map[models.AccessType]*fsm.State
	/* Registration procedure related context */
	RegistrationType5GS                uint8
	IdentityTypeUsedForRegistration    uint8
	RegistrationRequest                *nasMessage.RegistrationRequest
	ServingAmfChanged                  bool
	DeregistrationTargetAccessType     uint8 // only used when deregistration procedure is initialized by the network
	RegistrationAcceptForNon3GPPAccess []byte
	RetransmissionOfInitialNASMsg      bool
	/* Used for AMF relocation */
	TargetAmfProfile *models.NfProfile
	TargetAmfUri     string
	/* Ue Identity*/
	PlmnId              models.PlmnId
	Suci                string
	Supi                string
	UnauthenticatedSupi bool
	Gpsi                string
	Pei                 string
	Tmsi                int32 // 5G-Tmsi
	Guti                string
	GroupID             string
	EBI                 int32
	/* Ue Identity*/
	EventSubscriptionsInfo map[string]*AmfUeEventSubscription
	/* User Location*/
	RatType                  models.RatType
	Location                 models.UserLocation
	Tai                      models.Tai
	LocationChanged          bool
	LastVisitedRegisteredTai models.Tai
	TimeZone                 string
	/* context about udm */
	UdmId                             string
	NudmUECMUri                       string
	NudmSDMUri                        string
	ContextValid                      bool
	Reachability                      models.UeReachability
	SubscribedData                    models.SubscribedData
	SmfSelectionData                  *models.SmfSelectionSubscriptionData
	UeContextInSmfData                *models.UeContextInSmfData
	TraceData                         *models.TraceData
	UdmGroupId                        string
	SubscribedNssai                   []models.SubscribedSnssai
	AccessAndMobilitySubscriptionData *models.AccessAndMobilitySubscriptionData
	/* contex abut ausf */
	AusfGroupId                       string
	AusfId                            string
	AusfUri                           string
	RoutingIndicator                  string
	AuthenticationCtx                 *models.UeAuthenticationCtx
	AuthFailureCauseSynchFailureTimes int
	ABBA                              []uint8
	Kseaf                             string
	Kamf                              string
	/* context about PCF */
	PcfId                        string
	PcfUri                       string
	PolicyAssociationId          string
	AmPolicyUri                  string
	AmPolicyAssociation          *models.PolicyAssociation
	RequestTriggerLocationChange bool // true if AmPolicyAssociation.Trigger contains RequestTrigger_LOC_CH
	ConfigurationUpdateMessage   []byte
	/* UeContextForHandover*/
	HandoverNotifyUri string
	/* N1N2Message */
	N1N2MessageIDGenerator          *idgenerator.IDGenerator
	N1N2Message                     *N1N2Message
	N1N2MessageSubscribeIDGenerator *idgenerator.IDGenerator
	// map[int64]models.UeN1N2InfoSubscriptionCreateData; use n1n2MessageSubscriptionID as key
	N1N2MessageSubscription sync.Map
	/* Pdu Sesseion context */
	SmContextList sync.Map // map[int32]*SmContext, pdu session id as key
	/* Related Context*/
	RanUe map[models.AccessType]*RanUe
	/* other */
	onGoing                       map[models.AccessType]*OnGoing
	UeRadioCapability             string // OCTET string
	Capability5GMM                nasType.Capability5GMM
	ConfigurationUpdateIndication nasType.ConfigurationUpdateIndication
	/* context related to Paging */
	UeRadioCapabilityForPaging                 *UERadioCapabilityForPaging
	InfoOnRecommendedCellsAndRanNodesForPaging *InfoOnRecommendedCellsAndRanNodesForPaging
	UESpecificDRX                              uint8
	/* Security Context */
	SecurityContextAvailable bool
	UESecurityCapability     nasType.UESecurityCapability // for security command
	NgKsi                    models.NgKsi
	MacFailed                bool      // set to true if the integrity check of current NAS message is failed
	KnasInt                  [16]uint8 // 16 byte
	KnasEnc                  [16]uint8 // 16 byte
	Kgnb                     []uint8   // 32 byte
	Kn3iwf                   []uint8   // 32 byte
	NH                       []uint8   // 32 byte
	NCC                      uint8     // 0..7
	ULCount                  security.Count
	DLCount                  security.Count
	CipheringAlg             uint8
	IntegrityAlg             uint8
	/* Registration Area */
	RegistrationArea map[models.AccessType][]models.Tai
	LadnInfo         []LADN
	/* Network Slicing related context and Nssf */
	NssfId                            string
	NssfUri                           string
	NetworkSliceInfo                  *models.AuthorizedNetworkSliceInfo
	AllowedNssai                      map[models.AccessType][]models.AllowedSnssai
	ConfiguredNssai                   []models.ConfiguredSnssai
	NetworkSlicingSubscriptionChanged bool
	SdmSubscriptionId                 string
	UeCmRegistered                    bool
	/* T3513(Paging) */
	T3513 *Timer // for paging
	/* T3565(Notification) */
	T3565 *Timer // for NAS Notification
	/* T3560 (for authentication request/security mode command retransmission) */
	T3560 *Timer
	/* T3550 (for registration accept retransmission) */
	T3550 *Timer
	/* T3522 (for deregistration request) */
	T3522 *Timer
	/* Ue Context Release Cause */
	ReleaseCause map[models.AccessType]*CauseAll
	/* T3502 (Assigned by AMF, and used by UE to initialize registration procedure) */
	T3502Value                      int // Second
	T3512Value                      int // default 54 min
	Non3gppDeregistrationTimerValue int // default 54 min
	Lock                            sync.Mutex

	// logger
	NASLog      *logrus.Entry
	GmmLog      *logrus.Entry
	ProducerLog *logrus.Entry
}

type AmfUeEventSubscription struct {
	Timestamp         time.Time
	AnyUe             bool
	RemainReports     *int32
	EventSubscription *models.AmfEventSubscription
}

type N1N2Message struct {
	Request     models.N1N2MessageTransferRequest
	Status      models.N1N2MessageTransferCause
	ResourceUri string
}

type OnGoing struct {
	Procedure OnGoingProcedure
	Ppi       int32 // Paging priority
}

type UERadioCapabilityForPaging struct {
	NR    string // OCTET string
	EUTRA string // OCTET string
}

// TS 38.413 9.3.1.100
type InfoOnRecommendedCellsAndRanNodesForPaging struct {
	RecommendedCells    []RecommendedCell  // RecommendedCellsForPaging
	RecommendedRanNodes []RecommendRanNode // RecommendedRanNodesForPaging
}

// TS 38.413 9.3.1.71
type RecommendedCell struct {
	NgRanCGI         NGRANCGI
	TimeStayedInCell *int64
}

// TS 38.413 9.3.1.101
type RecommendRanNode struct {
	Present         int32
	GlobalRanNodeId *models.GlobalRanNodeId
	Tai             *models.Tai
}

type NGRANCGI struct {
	Present  int32
	NRCGI    *models.Ncgi
	EUTRACGI *models.Ecgi
}

func (ue *AmfUe) init() {
	ue.servingAMF = AMF_Self()
	ue.State = make(map[models.AccessType]*fsm.State)
	ue.State[models.AccessType__3_GPP_ACCESS] = fsm.NewState(Deregistered)
	ue.State[models.AccessType_NON_3_GPP_ACCESS] = fsm.NewState(Deregistered)
	ue.UnauthenticatedSupi = true
	ue.EventSubscriptionsInfo = make(map[string]*AmfUeEventSubscription)
	ue.RanUe = make(map[models.AccessType]*RanUe)
	ue.RegistrationArea = make(map[models.AccessType][]models.Tai)
	ue.AllowedNssai = make(map[models.AccessType][]models.AllowedSnssai)
	ue.N1N2MessageIDGenerator = idgenerator.NewGenerator(1, 2147483647)
	ue.N1N2MessageSubscribeIDGenerator = idgenerator.NewGenerator(1, 2147483647)
	ue.onGoing = make(map[models.AccessType]*OnGoing)
	ue.onGoing[models.AccessType_NON_3_GPP_ACCESS] = new(OnGoing)
	ue.onGoing[models.AccessType_NON_3_GPP_ACCESS].Procedure = OnGoingProcedureNothing
	ue.onGoing[models.AccessType__3_GPP_ACCESS] = new(OnGoing)
	ue.onGoing[models.AccessType__3_GPP_ACCESS].Procedure = OnGoingProcedureNothing
	ue.ReleaseCause = make(map[models.AccessType]*CauseAll)
}

func (ue *AmfUe) ServingAMF() *AMFContext {
	return ue.servingAMF
}

func (ue *AmfUe) CmConnect(anType models.AccessType) bool {
	if _, ok := ue.RanUe[anType]; !ok {
		return false
	}
	return true
}

func (ue *AmfUe) CmIdle(anType models.AccessType) bool {
	return !ue.CmConnect(anType)
}

func (ue *AmfUe) Remove() {
	for _, ranUe := range ue.RanUe {
		if err := ranUe.Remove(); err != nil {
			logger.ContextLog.Errorf("Remove RanUe error: %v", err)
		}
	}
	tmsiGenerator.FreeID(int64(ue.Tmsi))
	if len(ue.Supi) > 0 {
		AMF_Self().UePool.Delete(ue.Supi)
	}
}

func (ue *AmfUe) DetachRanUe(anType models.AccessType) {
	delete(ue.RanUe, anType)
}

func (ue *AmfUe) AttachRanUe(ranUe *RanUe) {
	ue.RanUe[ranUe.Ran.AnType] = ranUe
	ranUe.AmfUe = ue
}

func (ue *AmfUe) GetAnType() models.AccessType {
	if ue.CmConnect(models.AccessType__3_GPP_ACCESS) {
		return models.AccessType__3_GPP_ACCESS
	} else if ue.CmConnect(models.AccessType_NON_3_GPP_ACCESS) {
		return models.AccessType_NON_3_GPP_ACCESS
	}
	return ""
}

func (ue *AmfUe) GetCmInfo() (cmInfos []models.CmInfo) {
	var cmInfo models.CmInfo
	cmInfo.AccessType = models.AccessType__3_GPP_ACCESS
	if ue.CmConnect(cmInfo.AccessType) {
		cmInfo.CmState = models.CmState_CONNECTED
	} else {
		cmInfo.CmState = models.CmState_IDLE
	}
	cmInfos = append(cmInfos, cmInfo)
	cmInfo.AccessType = models.AccessType_NON_3_GPP_ACCESS
	if ue.CmConnect(cmInfo.AccessType) {
		cmInfo.CmState = models.CmState_CONNECTED
	} else {
		cmInfo.CmState = models.CmState_IDLE
	}
	cmInfos = append(cmInfos, cmInfo)
	return
}

func (ue *AmfUe) InAllowedNssai(targetSNssai models.Snssai, anType models.AccessType) bool {
	for _, allowedSnssai := range ue.AllowedNssai[anType] {
		if reflect.DeepEqual(*allowedSnssai.AllowedSnssai, targetSNssai) {
			return true
		}
	}
	return false
}

func (ue *AmfUe) InSubscribedNssai(targetSNssai models.Snssai) bool {
	for _, sNssai := range ue.SubscribedNssai {
		if reflect.DeepEqual(*sNssai.SubscribedSnssai, targetSNssai) {
			return true
		}
	}
	return false
}

func (ue *AmfUe) GetNsiInformationFromSnssai(anType models.AccessType, snssai models.Snssai) *models.NsiInformation {
	for _, allowedSnssai := range ue.AllowedNssai[anType] {
		if reflect.DeepEqual(*allowedSnssai.AllowedSnssai, snssai) {
			// TODO: select NsiInformation based on operator policy
			if len(allowedSnssai.NsiInformationList) != 0 {
				return &allowedSnssai.NsiInformationList[0]
			}
		}
	}
	return nil
}

func (ue *AmfUe) TaiListInRegistrationArea(taiList []models.Tai, accessType models.AccessType) bool {
	for _, tai := range taiList {
		if !InTaiList(tai, ue.RegistrationArea[accessType]) {
			return false
		}
	}
	return true
}

func (ue *AmfUe) HasWildCardSubscribedDNN() bool {
	for _, snssaiInfo := range ue.SmfSelectionData.SubscribedSnssaiInfos {
		for _, dnnInfo := range snssaiInfo.DnnInfos {
			if dnnInfo.Dnn == "*" {
				return true
			}
		}
	}
	return false
}

func (ue *AmfUe) SecurityContextIsValid() bool {
	return ue.SecurityContextAvailable && ue.NgKsi.Ksi != nasMessage.NasKeySetIdentifierNoKeyIsAvailable && !ue.MacFailed
}

// Kamf Derivation function defined in TS 33.501 Annex A.7
func (ue *AmfUe) DerivateKamf() {
	supiRegexp, err := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	groups := supiRegexp.FindStringSubmatch(ue.Supi)
	if groups == nil {
		logger.NasLog.Errorln("supi is not correct")
		return
	}

	P0 := []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 := ue.ABBA
	L1 := UeauCommon.KDFLen(P1)

	KseafDecode, err := hex.DecodeString(ue.Kseaf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	KamfBytes := UeauCommon.GetKDFValue(KseafDecode, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
	ue.Kamf = hex.EncodeToString(KamfBytes)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *AmfUe) DerivateAlgKey() {
	// Security Key
	P0 := []byte{security.NNASEncAlg}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.CipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, err := hex.DecodeString(ue.Kamf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	kenc := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasEnc[:], kenc[16:32])

	// Integrity Key
	P0 = []byte{security.NNASIntAlg}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntegrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasInt[:], kint[16:32])
}

// Access Network key Derivation function defined in TS 33.501 Annex A.9
func (ue *AmfUe) DerivateAnKey(anType models.AccessType) {
	accessType := security.AccessType3GPP // Defalut 3gpp
	P0 := make([]byte, 4)
	binary.BigEndian.PutUint32(P0, ue.ULCount.Get())
	L0 := UeauCommon.KDFLen(P0)
	if anType == models.AccessType_NON_3_GPP_ACCESS {
		accessType = security.AccessTypeNon3GPP
	}
	P1 := []byte{accessType}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, err := hex.DecodeString(ue.Kamf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	key := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_KGNB_KN3IWF_DERIVATION, P0, L0, P1, L1)
	switch accessType {
	case security.AccessType3GPP:
		ue.Kgnb = key
	case security.AccessTypeNon3GPP:
		ue.Kn3iwf = key
	}
}

// NH Derivation function defined in TS 33.501 Annex A.10
func (ue *AmfUe) DerivateNH(syncInput []byte) {
	P0 := syncInput
	L0 := UeauCommon.KDFLen(P0)

	KamfBytes, err := hex.DecodeString(ue.Kamf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	ue.NH = UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_NH_DERIVATION, P0, L0)
}

func (ue *AmfUe) UpdateSecurityContext(anType models.AccessType) {
	ue.DerivateAnKey(anType)
	switch anType {
	case models.AccessType__3_GPP_ACCESS:
		ue.DerivateNH(ue.Kgnb)
	case models.AccessType_NON_3_GPP_ACCESS:
		ue.DerivateNH(ue.Kn3iwf)
	}
	ue.NCC = 1
}

func (ue *AmfUe) UpdateNH() {
	ue.NCC++
	ue.DerivateNH(ue.NH)
}

func (ue *AmfUe) SelectSecurityAlg(intOrder, encOrder []uint8) {
	ue.CipheringAlg = security.AlgCiphering128NEA0
	ue.IntegrityAlg = security.AlgIntegrity128NIA0

	ueSupported := uint8(0)
	for _, intAlg := range intOrder {
		switch intAlg {
		case security.AlgIntegrity128NIA0:
			ueSupported = ue.UESecurityCapability.GetIA0_5G()
		case security.AlgIntegrity128NIA1:
			ueSupported = ue.UESecurityCapability.GetIA1_128_5G()
		case security.AlgIntegrity128NIA2:
			ueSupported = ue.UESecurityCapability.GetIA2_128_5G()
		case security.AlgIntegrity128NIA3:
			ueSupported = ue.UESecurityCapability.GetIA3_128_5G()
		}
		if ueSupported == 1 {
			ue.IntegrityAlg = intAlg
			break
		}
	}

	ueSupported = uint8(0)
	for _, encAlg := range encOrder {
		switch encAlg {
		case security.AlgCiphering128NEA0:
			ueSupported = ue.UESecurityCapability.GetEA0_5G()
		case security.AlgCiphering128NEA1:
			ueSupported = ue.UESecurityCapability.GetEA1_128_5G()
		case security.AlgCiphering128NEA2:
			ueSupported = ue.UESecurityCapability.GetEA2_128_5G()
		case security.AlgCiphering128NEA3:
			ueSupported = ue.UESecurityCapability.GetEA3_128_5G()
		}
		if ueSupported == 1 {
			ue.CipheringAlg = encAlg
			break
		}
	}
}

func (ue *AmfUe) ClearRegistrationRequestData(accessType models.AccessType) {
	ue.RegistrationRequest = nil
	ue.RegistrationType5GS = 0
	ue.IdentityTypeUsedForRegistration = 0
	ue.AuthFailureCauseSynchFailureTimes = 0
	ue.ServingAmfChanged = false
	ue.RegistrationAcceptForNon3GPPAccess = nil
	ue.RanUe[accessType].UeContextRequest = false
	ue.RetransmissionOfInitialNASMsg = false
	ue.onGoing[accessType].Procedure = OnGoingProcedureNothing
}

func (ue *AmfUe) SetOnGoing(anType models.AccessType, onGoing *OnGoing) {
	prevOnGoing := ue.onGoing[anType]
	ue.onGoing[anType] = onGoing
	ue.GmmLog.Debugf("OnGoing[%s]->[%s] PPI[%d]->[%d]", prevOnGoing.Procedure, onGoing.Procedure,
		prevOnGoing.Ppi, onGoing.Ppi)
}

func (ue *AmfUe) OnGoing(anType models.AccessType) OnGoing {
	return *ue.onGoing[anType]
}

func (ue *AmfUe) RemoveAmPolicyAssociation() {
	ue.AmPolicyAssociation = nil
	ue.PolicyAssociationId = ""
}

func (ue *AmfUe) CopyDataFromUeContextModel(ueContext models.UeContext) {
	if ueContext.Supi != "" {
		ue.Supi = ueContext.Supi
		ue.UnauthenticatedSupi = ueContext.SupiUnauthInd
	}

	if ueContext.Pei != "" {
		ue.Pei = ueContext.Pei
	}

	if ueContext.UdmGroupId != "" {
		ue.UdmGroupId = ueContext.UdmGroupId
	}

	if ueContext.AusfGroupId != "" {
		ue.AusfGroupId = ueContext.AusfGroupId
	}

	if ueContext.RoutingIndicator != "" {
		ue.RoutingIndicator = ueContext.RoutingIndicator
	}

	if ueContext.SubUeAmbr != nil {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		if ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr == nil {
			ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr = new(models.AmbrRm)
		}

		subAmbr := ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr
		subAmbr.Uplink = ueContext.SubUeAmbr.Uplink
		subAmbr.Downlink = ueContext.SubUeAmbr.Downlink
	}

	if ueContext.SubRfsp != 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.RfspIndex = ueContext.SubRfsp
	}

	if len(ueContext.RestrictedRatList) > 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.RatRestrictions = ueContext.RestrictedRatList
	}

	if len(ueContext.ForbiddenAreaList) > 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.ForbiddenAreas = ueContext.ForbiddenAreaList
	}

	if ueContext.ServiceAreaRestriction != nil {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.ServiceAreaRestriction = ueContext.ServiceAreaRestriction
	}

	if ueContext.SeafData != nil {
		seafData := ueContext.SeafData

		ue.NgKsi = *seafData.NgKsi
		if seafData.KeyAmf != nil {
			if seafData.KeyAmf.KeyType == models.KeyAmfType_KAMF {
				ue.Kamf = seafData.KeyAmf.KeyVal
			}
		}
		if nh, err := hex.DecodeString(seafData.Nh); err != nil {
			logger.ContextLog.Error(err)
			return
		} else {
			ue.NH = nh
		}
		ue.NCC = uint8(seafData.Ncc)
	}

	if ueContext.PcfId != "" {
		ue.PcfId = ueContext.PcfId
	}

	if ueContext.PcfAmPolicyUri != "" {
		ue.AmPolicyUri = ueContext.PcfAmPolicyUri
	}

	if len(ueContext.AmPolicyReqTriggerList) > 0 {
		if ue.AmPolicyAssociation == nil {
			ue.AmPolicyAssociation = new(models.PolicyAssociation)
		}
		for _, trigger := range ueContext.AmPolicyReqTriggerList {
			switch trigger {
			case models.AmPolicyReqTrigger_LOCATION_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_LOC_CH)
			case models.AmPolicyReqTrigger_PRA_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_PRA_CH)
			case models.AmPolicyReqTrigger_SARI_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_SERV_AREA_CH)
			case models.AmPolicyReqTrigger_RFSP_INDEX_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_RFSP_CH)
			}
		}
	}

	if len(ueContext.SessionContextList) > 0 {
		for _, pduSessionContext := range ueContext.SessionContextList {
			smContext := SmContext{
				pduSessionID: pduSessionContext.PduSessionId,
				smContextRef: pduSessionContext.SmContextRef,
				snssai:       *pduSessionContext.SNssai,
				dnn:          pduSessionContext.Dnn,
				accessType:   pduSessionContext.AccessType,
				hSmfID:       pduSessionContext.HsmfId,
				vSmfID:       pduSessionContext.VsmfId,
				nsInstance:   pduSessionContext.NsInstance,
			}
			ue.StoreSmContext(pduSessionContext.PduSessionId, &smContext)
		}
	}

	if len(ueContext.MmContextList) > 0 {
		for _, mmContext := range ueContext.MmContextList {
			if mmContext.AccessType == models.AccessType__3_GPP_ACCESS {
				if nasSecurityMode := mmContext.NasSecurityMode; nasSecurityMode != nil {
					switch nasSecurityMode.IntegrityAlgorithm {
					case models.IntegrityAlgorithm_NIA0:
						ue.IntegrityAlg = security.AlgIntegrity128NIA0
					case models.IntegrityAlgorithm_NIA1:
						ue.IntegrityAlg = security.AlgIntegrity128NIA1
					case models.IntegrityAlgorithm_NIA2:
						ue.IntegrityAlg = security.AlgIntegrity128NIA2
					case models.IntegrityAlgorithm_NIA3:
						ue.IntegrityAlg = security.AlgIntegrity128NIA3
					}

					switch nasSecurityMode.CipheringAlgorithm {
					case models.CipheringAlgorithm_NEA0:
						ue.CipheringAlg = security.AlgCiphering128NEA0
					case models.CipheringAlgorithm_NEA1:
						ue.CipheringAlg = security.AlgCiphering128NEA1
					case models.CipheringAlgorithm_NEA2:
						ue.CipheringAlg = security.AlgCiphering128NEA2
					case models.CipheringAlgorithm_NEA3:
						ue.CipheringAlg = security.AlgCiphering128NEA3
					}

					if mmContext.NasDownlinkCount != 0 {
						overflow := uint16((uint32(mmContext.NasDownlinkCount) & 0x00ffff00) >> 8)
						sqn := uint8(uint32(mmContext.NasDownlinkCount & 0x000000ff))
						ue.DLCount.Set(overflow, sqn)
					}

					if mmContext.NasUplinkCount != 0 {
						overflow := uint16((uint32(mmContext.NasUplinkCount) & 0x00ffff00) >> 8)
						sqn := uint8(uint32(mmContext.NasUplinkCount & 0x000000ff))
						ue.ULCount.Set(overflow, sqn)
					}

					// TS 29.518 Table 6.1.6.3.2.1
					if mmContext.UeSecurityCapability != "" {
						// ue.SecurityCapabilities
						buf, err := base64.StdEncoding.DecodeString(mmContext.UeSecurityCapability)
						if err != nil {
							logger.ContextLog.Error(err)
							return
						}
						ue.UESecurityCapability.Buffer = buf
						ue.UESecurityCapability.SetLen(uint8(len(buf)))
					}
				}
			}

			if mmContext.AllowedNssai != nil {
				for _, snssai := range mmContext.AllowedNssai {
					allowedSnssai := models.AllowedSnssai{
						AllowedSnssai: &snssai,
					}
					ue.AllowedNssai[mmContext.AccessType] = append(ue.AllowedNssai[mmContext.AccessType], allowedSnssai)
				}
			}
		}
	}
	if ueContext.TraceData != nil {
		ue.TraceData = ueContext.TraceData
	}
}

// SM Context realted function

func (ue *AmfUe) StoreSmContext(pduSessionID int32, smContext *SmContext) {
	ue.SmContextList.Store(pduSessionID, smContext)
}

func (ue *AmfUe) SmContextFindByPDUSessionID(pduSessionID int32) (*SmContext, bool) {
	if value, ok := ue.SmContextList.Load(pduSessionID); ok {
		return value.(*SmContext), true
	} else {
		return nil, false
	}
}
