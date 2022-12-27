package util

import (
	"os"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/context"
	"github.com/free5gc/udm/factory"
	"github.com/free5gc/udm/logger"
)

func InitUDMContext(udmContext *context.UDMContext) {
	config := factory.UdmConfig
	logger.UtilLog.Info("udmconfig Info: Version[", config.Info.Version, "] Description[", config.Info.Description, "]")
	configuration := config.Configuration
	udmContext.NfId = uuid.New().String()
	if configuration.UdmName != "" {
		udmContext.Name = configuration.UdmName
	}
	sbi := configuration.Sbi
	udmContext.UriScheme = ""
	udmContext.SBIPort = factory.UDM_DEFAULT_PORT_INT
	udmContext.RegisterIPv4 = factory.UDM_DEFAULT_IPV4
	if sbi != nil {
		if sbi.Scheme != "" {
			udmContext.UriScheme = models.UriScheme(sbi.Scheme)
		}
		if sbi.RegisterIPv4 != "" {
			udmContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			udmContext.SBIPort = sbi.Port
		}

		udmContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if udmContext.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			udmContext.BindingIPv4 = sbi.BindingIPv4
			if udmContext.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				udmContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	udmContext.NrfUri = configuration.NrfUri
	servingNameList := configuration.ServiceNameList

	udmContext.Keys = configuration.Keys

	udmContext.InitNFService(servingNameList, config.Info.Version)
}
