package context

import (
	"fmt"
	"time"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/factory"
)

var NFServices *[]models.NfService

var NfServiceVersion *[]models.NfServiceVersion

var SmfInfo *models.SmfInfo

func SetupNFProfile(config *factory.Config) {
	// Set time
	nfSetupTime := time.Now()

	// set NfServiceVersion
	NfServiceVersion = &[]models.NfServiceVersion{
		{
			ApiVersionInUri: "v1",
			ApiFullVersion:  fmt.Sprintf("https://%s:%d/nsmf-pdusession/v1", SMF_Self().RegisterIPv4, SMF_Self().SBIPort),
			Expiry:          &nfSetupTime,
		},
	}

	// set NFServices
	NFServices = new([]models.NfService)
	for _, serviceName := range config.Configuration.ServiceNameList {
		*NFServices = append(*NFServices, models.NfService{
			ServiceInstanceId: SMF_Self().NfInstanceID + serviceName,
			ServiceName:       models.ServiceName(serviceName),
			Versions:          NfServiceVersion,
			Scheme:            models.UriScheme_HTTPS,
			NfServiceStatus:   models.NfServiceStatus_REGISTERED,
			ApiPrefix:         fmt.Sprintf("%s://%s:%d", SMF_Self().URIScheme, SMF_Self().RegisterIPv4, SMF_Self().SBIPort),
		})
	}

	// set smfInfo
	SmfInfo = &models.SmfInfo{
		SNssaiSmfInfoList: SNssaiSmfInfo(),
	}
}

func SNssaiSmfInfo() *[]models.SnssaiSmfInfoItem {
	snssaiInfo := make([]models.SnssaiSmfInfoItem, 0)
	for _, snssai := range smfContext.SnssaiInfos {
		var snssaiInfoModel models.SnssaiSmfInfoItem
		snssaiInfoModel.SNssai = &models.Snssai{
			Sst: snssai.Snssai.Sst,
			Sd:  snssai.Snssai.Sd,
		}
		dnnModelList := make([]models.DnnSmfInfoItem, 0)

		for dnn := range snssai.DnnInfos {
			dnnModelList = append(dnnModelList, models.DnnSmfInfoItem{
				Dnn: dnn,
			})
		}

		snssaiInfoModel.DnnSmfInfoList = &dnnModelList

		snssaiInfo = append(snssaiInfo, snssaiInfoModel)
	}

	return &snssaiInfo
}
