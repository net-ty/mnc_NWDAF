package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	aperLogger "github.com/free5gc/aper/logger"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	nasLogger "github.com/free5gc/nas/logger"
	ngapLogger "github.com/free5gc/ngap/logger"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
	pfcpLogger "github.com/free5gc/pfcp/logger"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/callback"
	"github.com/free5gc/smf/consumer"
	"github.com/free5gc/smf/context"
	"github.com/free5gc/smf/eventexposure"
	"github.com/free5gc/smf/factory"
	"github.com/free5gc/smf/logger"
	"github.com/free5gc/smf/oam"
	"github.com/free5gc/smf/pdusession"
	"github.com/free5gc/smf/pfcp"
	"github.com/free5gc/smf/pfcp/message"
	"github.com/free5gc/smf/pfcp/udp"
	"github.com/free5gc/smf/util"
)

type SMF struct{}

type (
	// Config information.
	Config struct {
		smfcfg    string
		uerouting string
	}
)

var config Config

var smfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "smfcfg",
		Usage: "config file",
	},
	cli.StringFlag{
		Name:  "uerouting",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*SMF) GetCliCmd() (flags []cli.Flag) {
	return smfCLi
}

func (smf *SMF) Initialize(c *cli.Context) error {
	config = Config{
		smfcfg:    c.String("smfcfg"),
		uerouting: c.String("uerouting"),
	}

	if config.smfcfg != "" {
		if err := factory.InitConfigFactory(config.smfcfg); err != nil {
			return err
		}
	} else {
		DefaultSmfConfigPath := path_util.Free5gcPath("free5gc/config/smfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultSmfConfigPath); err != nil {
			return err
		}
	}

	if config.uerouting != "" {
		if err := factory.InitRoutingConfigFactory(config.uerouting); err != nil {
			return err
		}
	} else {
		DefaultUERoutingPath := path_util.Free5gcPath("free5gc/config/uerouting.yaml")
		if err := factory.InitRoutingConfigFactory(DefaultUERoutingPath); err != nil {
			return err
		}
	}

	smf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (smf *SMF) setLogLevel() {
	if factory.SmfConfig.Logger == nil {
		initLog.Warnln("SMF config without log level setting!!!")
		return
	}

	if factory.SmfConfig.Logger.SMF != nil {
		if factory.SmfConfig.Logger.SMF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.SMF.DebugLevel); err != nil {
				initLog.Warnf("SMF Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.SMF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("SMF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("SMF Log level is default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.SmfConfig.Logger.SMF.ReportCaller)
	}

	if factory.SmfConfig.Logger.NAS != nil {
		if factory.SmfConfig.Logger.NAS.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.NAS.DebugLevel); err != nil {
				nasLogger.NasLog.Warnf("NAS Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.NAS.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				nasLogger.SetLogLevel(level)
			}
		} else {
			nasLogger.NasLog.Warnln("NAS Log level not set. Default set to [info] level")
			nasLogger.SetLogLevel(logrus.InfoLevel)
		}
		nasLogger.SetReportCaller(factory.SmfConfig.Logger.NAS.ReportCaller)
	}

	if factory.SmfConfig.Logger.NGAP != nil {
		if factory.SmfConfig.Logger.NGAP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.NGAP.DebugLevel); err != nil {
				ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.NGAP.DebugLevel)
				ngapLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				ngapLogger.SetLogLevel(level)
			}
		} else {
			ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
			ngapLogger.SetLogLevel(logrus.InfoLevel)
		}
		ngapLogger.SetReportCaller(factory.SmfConfig.Logger.NGAP.ReportCaller)
	}

	if factory.SmfConfig.Logger.Aper != nil {
		if factory.SmfConfig.Logger.Aper.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.Aper.DebugLevel); err != nil {
				aperLogger.AperLog.Warnf("Aper Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.Aper.DebugLevel)
				aperLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				aperLogger.SetLogLevel(level)
			}
		} else {
			aperLogger.AperLog.Warnln("Aper Log level not set. Default set to [info] level")
			aperLogger.SetLogLevel(logrus.InfoLevel)
		}
		aperLogger.SetReportCaller(factory.SmfConfig.Logger.Aper.ReportCaller)
	}

	if factory.SmfConfig.Logger.PathUtil != nil {
		if factory.SmfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.SmfConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.SmfConfig.Logger.OpenApi != nil {
		if factory.SmfConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.SmfConfig.Logger.OpenApi.ReportCaller)
	}

	if factory.SmfConfig.Logger.PFCP != nil {
		if factory.SmfConfig.Logger.PFCP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.SmfConfig.Logger.PFCP.DebugLevel); err != nil {
				pfcpLogger.PFCPLog.Warnf("PFCP Log level [%s] is invalid, set to [info] level",
					factory.SmfConfig.Logger.PFCP.DebugLevel)
				pfcpLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pfcpLogger.SetLogLevel(level)
			}
		} else {
			pfcpLogger.PFCPLog.Warnln("PFCP Log level not set. Default set to [info] level")
			pfcpLogger.SetLogLevel(logrus.InfoLevel)
		}
		pfcpLogger.SetReportCaller(factory.SmfConfig.Logger.PFCP.ReportCaller)
	}
}

func (smf *SMF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range smf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (smf *SMF) Start() {
	context.InitSmfContext(&factory.SmfConfig)
	// allocate id for each upf
	context.AllocateUPFID()
	context.InitSMFUERouting(&factory.UERoutingConfig)

	initLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	err := consumer.SendNFRegistration()
	if err != nil {
		retry_err := consumer.RetrySendNFRegistration(10)
		if retry_err != nil {
			logger.InitLog.Errorln(retry_err)
			return
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		smf.Terminate()
		os.Exit(0)
	}()

	oam.AddService(router)
	callback.AddService(router)
	for _, serviceName := range factory.SmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NSMF_PDUSESSION:
			pdusession.AddService(router)
		case models.ServiceName_NSMF_EVENT_EXPOSURE:
			eventexposure.AddService(router)
		}
	}
	udp.Run(pfcp.Dispatch)

	for _, upf := range context.SMF_Self().UserPlaneInformation.UPFs {
		if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
			logger.AppLog.Infof("Send PFCP Association Request to UPF[%s](%s)\n", upf.NodeID.NodeIdValue,
				upf.NodeID.ResolveNodeIdToIp().String())
		} else {
			logger.AppLog.Infof("Send PFCP Association Request to UPF[%s]\n", upf.NodeID.ResolveNodeIdToIp().String())
		}
		message.SendPfcpAssociationSetupRequest(upf.NodeID)
	}

	time.Sleep(1000 * time.Millisecond)

	HTTPAddr := fmt.Sprintf("%s:%d", context.SMF_Self().BindingIPv4, context.SMF_Self().SBIPort)
	server, err := http2_util.NewServer(HTTPAddr, util.SmfLogPath, router)

	if server == nil {
		initLog.Error("Initialize HTTP server failed:", err)
		return
	}

	if err != nil {
		initLog.Warnln("Initialize HTTP server:", err)
	}

	serverScheme := factory.SmfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.SmfPemPath, util.SmfKeyPath)
	}

	if err != nil {
		initLog.Fatalln("HTTP server setup failed:", err)
	}
}

func (smf *SMF) Terminate() {
	logger.InitLog.Infof("Terminating SMF...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
}

func (smf *SMF) Exec(c *cli.Context) error {
	return nil
}
