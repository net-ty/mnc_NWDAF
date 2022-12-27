package context

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/free5gc/idgenerator"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/factory"
	"github.com/free5gc/pcf/logger"
)

var pcfCtx *PCFContext

func init() {
	pcfCtx = new(PCFContext)
	pcfCtx.Name = "pcf"
	pcfCtx.UriScheme = models.UriScheme_HTTPS
	pcfCtx.TimeFormat = "2006-01-02 15:04:05"
	pcfCtx.DefaultBdtRefId = "BdtPolicyId-"
	pcfCtx.NfService = make(map[models.ServiceName]models.NfService)
	pcfCtx.PcfServiceUris = make(map[models.ServiceName]string)
	pcfCtx.PcfSuppFeats = make(map[models.ServiceName]openapi.SupportedFeature)
	pcfCtx.BdtPolicyIDGenerator = idgenerator.NewGenerator(1, math.MaxInt64)
}

type PCFContext struct {
	NfId            string
	Name            string
	UriScheme       models.UriScheme
	BindingIPv4     string
	RegisterIPv4    string
	SBIPort         int
	TimeFormat      string
	DefaultBdtRefId string
	NfService       map[models.ServiceName]models.NfService
	PcfServiceUris  map[models.ServiceName]string
	PcfSuppFeats    map[models.ServiceName]openapi.SupportedFeature
	NrfUri          string
	DefaultUdrURI   string
	// UePool          map[string]*UeContext
	UePool sync.Map
	// Bdt Policy related
	BdtPolicyPool        sync.Map
	BdtPolicyIDGenerator *idgenerator.IDGenerator
	// App Session related
	AppSessionPool sync.Map
	// AMF Status Change Subscription related
	AMFStatusSubsData sync.Map // map[string]AMFStatusSubscriptionData; subscriptionID as key

	// lock
	DefaultUdrURILock sync.RWMutex
}

type AMFStatusSubscriptionData struct {
	AmfUri       string
	AmfStatusUri string
	GuamiList    []models.Guami
}

type AppSessionData struct {
	AppSessionId      string
	AppSessionContext *models.AppSessionContext
	// (compN/compN-subCompN/appId-%s) map to PccRule
	RelatedPccRuleIds    map[string]string
	PccRuleIdMapToCompId map[string]string
	// EventSubscription
	Events   map[models.AfEvent]models.AfNotifMethod
	EventUri string
	// related Session
	SmPolicyData *UeSmPolicyData
}

// Create new PCF context
func PCF_Self() *PCFContext {
	return pcfCtx
}

func GetTimeformat() string {
	return pcfCtx.TimeFormat
}

func GetUri(name models.ServiceName) string {
	return pcfCtx.PcfServiceUris[name]
}

var (
	PolicyAuthorizationUri = "/npcf-policyauthorization/v1/app-sessions/"
	SmUri                  = "/npcf-smpolicycontrol/v1"
	IPv4Address            = "192.168."
	IPv6Address            = "ffab::"
	CheckNotifiUri         = "/npcf-callback/v1/nudr-notify/"
	Ipv4_pool              = make(map[string]string)
	Ipv6_pool              = make(map[string]string)
)

// BdtPolicy default value
const DefaultBdtRefId = "BdtPolicyId-"

func (c *PCFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", c.UriScheme, c.RegisterIPv4, c.SBIPort)
}

// Init NfService with supported service list ,and version of services
func (c *PCFContext) InitNFService(serviceList []factory.Service, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, service := range serviceList {
		name := models.ServiceName(service.ServiceName)
		c.NfService[name] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          c.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       c.GetIPv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: c.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(c.SBIPort),
				},
			},
			SupportedFeatures: service.SuppFeat,
		}
	}
}

// Allocate PCF Ue with supi and add to pcf Context and returns allocated ue
func (c *PCFContext) NewPCFUe(Supi string) (*UeContext, error) {
	if strings.HasPrefix(Supi, "imsi-") {
		newUeContext := &UeContext{}
		newUeContext.SmPolicyData = make(map[string]*UeSmPolicyData)
		newUeContext.AMPolicyData = make(map[string]*UeAMPolicyData)
		newUeContext.PolAssociationIDGenerator = 1
		newUeContext.AppSessionIDGenerator = idgenerator.NewGenerator(1, math.MaxInt64)
		newUeContext.Supi = Supi
		c.UePool.Store(Supi, newUeContext)
		return newUeContext, nil
	} else {
		return nil, fmt.Errorf(" add Ue context fail ")
	}
}

// Return Bdt Policy Id with format "BdtPolicyId-%d" which be allocated
func (c *PCFContext) AllocBdtPolicyID() (bdtPolicyID string, err error) {
	var allocID int64
	if allocID, err = c.BdtPolicyIDGenerator.Allocate(); err != nil {
		logger.CtxLog.Warnf("Allocate pathID error: %+v", err)
		return "", err
	}

	bdtPolicyID = fmt.Sprintf("BdtPolicyId-%d", allocID)
	return bdtPolicyID, nil
}

// Find PcfUe which the policyId belongs to
func (c *PCFContext) PCFUeFindByPolicyId(PolicyId string) *UeContext {
	index := strings.LastIndex(PolicyId, "-")
	if index == -1 {
		return nil
	}
	supi := PolicyId[:index]
	if supi != "" {
		if value, ok := c.UePool.Load(supi); ok {
			ueContext := value.(*UeContext)
			return ueContext
		}
	}
	return nil
}

// Find PcfUe which the AppSessionId belongs to
func (c *PCFContext) PCFUeFindByAppSessionId(appSessionId string) *UeContext {
	index := strings.LastIndex(appSessionId, "-")
	if index == -1 {
		return nil
	}
	supi := appSessionId[:index]
	if supi != "" {
		if value, ok := c.UePool.Load(supi); ok {
			ueContext := value.(*UeContext)
			return ueContext
		}
	}
	return nil
}

// Find PcfUe which Ipv4 belongs to
func (c *PCFContext) PcfUeFindByIPv4(v4 string) *UeContext {
	var ue *UeContext
	c.UePool.Range(func(key, value interface{}) bool {
		ue = value.(*UeContext)
		if ue.SMPolicyFindByIpv4(v4) != nil {
			return false
		} else {
			return true
		}
	})

	return ue
}

// Find PcfUe which Ipv6 belongs to
func (c *PCFContext) PcfUeFindByIPv6(v6 string) *UeContext {
	var ue *UeContext
	c.UePool.Range(func(key, value interface{}) bool {
		ue = value.(*UeContext)
		if ue.SMPolicyFindByIpv6(v6) != nil {
			return false
		} else {
			return true
		}
	})

	return ue
}

// Find SMPolicy with AppSessionContext
func ueSMPolicyFindByAppSessionContext(ue *UeContext, req *models.AppSessionContextReqData) (*UeSmPolicyData, error) {
	var policy *UeSmPolicyData
	var err error

	if req.UeIpv4 != "" {
		policy = ue.SMPolicyFindByIdentifiersIpv4(req.UeIpv4, req.SliceInfo, req.Dnn, req.IpDomain)
		if policy == nil {
			err = fmt.Errorf("Can't find Ue with Ipv4[%s]", req.UeIpv4)
		}
	} else if req.UeIpv6 != "" {
		policy = ue.SMPolicyFindByIdentifiersIpv6(req.UeIpv6, req.SliceInfo, req.Dnn)
		if policy == nil {
			err = fmt.Errorf("Can't find Ue with Ipv6 prefix[%s]", req.UeIpv6)
		}
	} else {
		// TODO: find by MAC address
		err = fmt.Errorf("Ue finding by MAC address does not support")
	}
	return policy, err
}

// SessionBinding from application request to get corresponding Sm policy
func (c *PCFContext) SessionBinding(req *models.AppSessionContextReqData) (*UeSmPolicyData, error) {
	var selectedUE *UeContext
	var policy *UeSmPolicyData
	var err error

	if req.Supi != "" {
		if val, exist := c.UePool.Load(req.Supi); exist {
			selectedUE = val.(*UeContext)
		}
	}

	if req.Gpsi != "" && selectedUE == nil {
		c.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*UeContext)
			if ue.Gpsi == req.Gpsi {
				selectedUE = ue
				return false
			} else {
				return true
			}
		})
	}

	if selectedUE != nil {
		policy, err = ueSMPolicyFindByAppSessionContext(selectedUE, req)
	} else {
		c.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*UeContext)
			policy, err = ueSMPolicyFindByAppSessionContext(ue, req)
			return true
		})
	}
	if policy == nil && err == nil {
		err = fmt.Errorf("No SM policy found")
	}
	return policy, err
}

// SetDefaultUdrURI ... function to set DefaultUdrURI
func (c *PCFContext) SetDefaultUdrURI(uri string) {
	c.DefaultUdrURILock.Lock()
	defer c.DefaultUdrURILock.Unlock()
	c.DefaultUdrURI = uri
}

func Ipv4Pool(ipindex int32) string {
	ipv4address := IPv4Address + fmt.Sprint((int(ipindex)/255)+1) + "." + fmt.Sprint(int(ipindex)%255)
	return ipv4address
}

func Ipv4Index() int32 {
	if len(Ipv4_pool) == 0 {
		Ipv4_pool["1"] = Ipv4Pool(1)
	} else {
		for i := 1; i <= len(Ipv4_pool); i++ {
			if Ipv4_pool[fmt.Sprint(i)] == "" {
				Ipv4_pool[fmt.Sprint(i)] = Ipv4Pool(int32(i))
				return int32(i)
			}
		}

		Ipv4_pool[fmt.Sprint(int32(len(Ipv4_pool)+1))] = Ipv4Pool(int32(len(Ipv4_pool) + 1))
		return int32(len(Ipv4_pool))
	}
	return 1
}

func GetIpv4Address(ipindex int32) string {
	return Ipv4_pool[fmt.Sprint(ipindex)]
}

func DeleteIpv4index(Ipv4index int32) {
	delete(Ipv4_pool, fmt.Sprint(Ipv4index))
}

func Ipv6Pool(ipindex int32) string {
	ipv6address := IPv6Address + fmt.Sprintf("%x\n", ipindex)
	return ipv6address
}

func Ipv6Index() int32 {
	if len(Ipv6_pool) == 0 {
		Ipv6_pool["1"] = Ipv6Pool(1)
	} else {
		for i := 1; i <= len(Ipv6_pool); i++ {
			if Ipv6_pool[fmt.Sprint(i)] == "" {
				Ipv6_pool[fmt.Sprint(i)] = Ipv6Pool(int32(i))
				return int32(i)
			}
		}

		Ipv6_pool[fmt.Sprint(int32(len(Ipv6_pool)+1))] = Ipv6Pool(int32(len(Ipv6_pool) + 1))
		return int32(len(Ipv6_pool))
	}
	return 1
}

func GetIpv6Address(ipindex int32) string {
	return Ipv6_pool[fmt.Sprint(ipindex)]
}

func DeleteIpv6index(Ipv6index int32) {
	delete(Ipv6_pool, fmt.Sprint(Ipv6index))
}

func (c *PCFContext) NewAmfStatusSubscription(subscriptionID string, subscriptionData AMFStatusSubscriptionData) {
	c.AMFStatusSubsData.Store(subscriptionID, subscriptionData)
}
