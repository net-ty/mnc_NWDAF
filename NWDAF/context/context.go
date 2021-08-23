package context

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/Nudm_SubscriberDataManagement"
	"github.com/free5gc/openapi/models"
	"nwdaf.com/factory"
	"nwdaf.com/logger"
)

func init() {
	nwdafContext.NfInstanceID = uuid.New().String()
}

var nwdafContext NWDAFContext

type NWDAFContext struct {
	Name            string
	URIScheme       models.UriScheme
	UriScheme       models.UriScheme
	BindingIPv4     string
	RegisterIPv4    string
	SBIPort         int
	HttpIPv6Address string
	NfInstanceID    string
	NfId            string
	// Key    string
	// PEM    string
	// KeyLog string

	NrfUri string

	// Now only "IPv4" supported
	// TODO: support "IPv6", "IPv4v6", "Ethernet"
	SupportedPDUSessionType string

	//*** For ULCL ** //
	// ULCLSupport    bool
	// ULCLGroups     map[string][]string
	// LocalSEIDCount uint64

	NFManagementClient             *Nnrf_NFManagement.APIClient
	NFDiscoveryClient              *Nnrf_NFDiscovery.APIClient
	SubscriberDataManagementClient *Nudm_SubscriberDataManagement.APIClient
}

// RetrieveDnnInformation gets the corresponding dnn info from S-NSSAI and DNN

// func AllocateLocalSEID() uint64 {
// 	atomic.AddUint64(&nwdafContext.LocalSEIDCount, 1)
// 	return nwdafContext.LocalSEIDCount //if error delete this
// }

func InitNwdafContext(config *factory.Config) {
	if config == nil {
		logger.CtxLog.Error("Config is nil")
		return
	}

	logger.CtxLog.Infof("nwdafconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	if configuration.NwdafName != "" {
		nwdafContext.Name = configuration.NwdafName
	}

	sbi := configuration.Sbi
	if sbi == nil {
		logger.CtxLog.Errorln("Configuration needs \"sbi\" value")
		return
	} else {
		nwdafContext.URIScheme = models.UriScheme(sbi.Scheme)
		nwdafContext.RegisterIPv4 = factory.NWDAF_DEFAULT_IPV4 // default localhost
		nwdafContext.SBIPort = factory.NWDAF_DEFAULT_PORT_INT  // default port
		if sbi.RegisterIPv4 != "" {
			nwdafContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			nwdafContext.SBIPort = sbi.Port
		}

		nwdafContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if nwdafContext.BindingIPv4 != "" {
			logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			nwdafContext.BindingIPv4 = sbi.BindingIPv4
			if nwdafContext.BindingIPv4 == "" {
				logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				nwdafContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}

	if configuration.NrfUri != "" {
		nwdafContext.NrfUri = configuration.NrfUri
	} else {
		logger.CtxLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		nwdafContext.NrfUri = fmt.Sprintf("%s://%s:%d", nwdafContext.URIScheme, "127.0.0.1", 29510)
	}

	//SetupNFProfile(config)
}

func NWDAF_Self() *NWDAFContext {
	return &nwdafContext
}
