package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/antihax/optional"
	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
	"github.com/free5gc/pcf/ampolicy"
	"github.com/free5gc/pcf/bdtpolicy"
	"github.com/free5gc/pcf/consumer"
	"github.com/free5gc/pcf/context"
	"github.com/free5gc/pcf/factory"
	"github.com/free5gc/pcf/httpcallback"
	"github.com/free5gc/pcf/internal/notifyevent"
	"github.com/free5gc/pcf/logger"
	"github.com/free5gc/pcf/oam"
	"github.com/free5gc/pcf/policyauthorization"
	"github.com/free5gc/pcf/smpolicy"
	"github.com/free5gc/pcf/uepolicy"
	"github.com/free5gc/pcf/util"
)

type PCF struct{}

type (
	// Config information.
	Config struct {
		pcfcfg string
	}
)

var config Config

var pcfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "pcfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*PCF) GetCliCmd() (flags []cli.Flag) {
	return pcfCLi
}

func (pcf *PCF) Initialize(c *cli.Context) error {
	config = Config{
		pcfcfg: c.String("pcfcfg"),
	}
	if config.pcfcfg != "" {
		if err := factory.InitConfigFactory(config.pcfcfg); err != nil {
			return err
		}
	} else {
		DefaultPcfConfigPath := path_util.Free5gcPath("free5gc/config/pcfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultPcfConfigPath); err != nil {
			return err
		}
	}

	pcf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (pcf *PCF) setLogLevel() {
	if factory.PcfConfig.Logger == nil {
		initLog.Warnln("PCF config without log level setting!!!")
		return
	}

	if factory.PcfConfig.Logger.PCF != nil {
		if factory.PcfConfig.Logger.PCF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.PcfConfig.Logger.PCF.DebugLevel); err != nil {
				initLog.Warnf("PCF Log level [%s] is invalid, set to [info] level",
					factory.PcfConfig.Logger.PCF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("PCF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("PCF Log level is default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.PcfConfig.Logger.PCF.ReportCaller)
	}

	if factory.PcfConfig.Logger.PathUtil != nil {
		if factory.PcfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.PcfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.PcfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.PcfConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.PcfConfig.Logger.OpenApi != nil {
		if factory.PcfConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.PcfConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.PcfConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.PcfConfig.Logger.OpenApi.ReportCaller)
	}
}

func (pcf *PCF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range pcf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (pcf *PCF) Start() {
	initLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	bdtpolicy.AddService(router)
	smpolicy.AddService(router)
	ampolicy.AddService(router)
	uepolicy.AddService(router)
	policyauthorization.AddService(router)
	httpcallback.AddService(router)
	oam.AddService(router)

	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	if err := notifyevent.RegisterNotifyDispatcher(); err != nil {
		initLog.Error("Register NotifyDispatcher Error")
	}

	self := context.PCF_Self()
	util.InitpcfContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	profile, err := consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build PCF Profile Error")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("PCF register to NRF Error[%s]", err.Error())
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NUDR_DR}),
	}
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_UDR, models.NfType_PCF, param)
	for _, nfProfile := range resp.NfInstances {
		udruri := util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED)
		if udruri != "" {
			self.SetDefaultUdrURI(udruri)
			break
		}
	}
	if err != nil {
		initLog.Errorln(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		pcf.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, util.PCF_LOG_PATH, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.PcfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.PCF_PEM_PATH, util.PCF_KEY_PATH)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (pcf *PCF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("pcfcfg"))
	args := pcf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./pcf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	go func() {
		in := bufio.NewScanner(stderr)
		fmt.Println("PCF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("PCF start")
		if err = command.Start(); err != nil {
			fmt.Printf("command.Start() error: %v", err)
		}
		fmt.Println("PCF end")
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (pcf *PCF) Terminate() {
	logger.InitLog.Infof("Terminating PCF...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
	logger.InitLog.Infof("PCF terminated")
}
