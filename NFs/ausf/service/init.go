package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/ausf/consumer"
	ausf_context "github.com/free5gc/ausf/context"
	"github.com/free5gc/ausf/factory"
	"github.com/free5gc/ausf/logger"
	"github.com/free5gc/ausf/ueauthentication"
	"github.com/free5gc/ausf/util"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
)

type AUSF struct{}

type (
	// Config information.
	Config struct {
		ausfcfg string
	}
)

var config Config

var ausfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "ausfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*AUSF) GetCliCmd() (flags []cli.Flag) {
	return ausfCLi
}

func (ausf *AUSF) Initialize(c *cli.Context) error {
	config = Config{
		ausfcfg: c.String("ausfcfg"),
	}

	if config.ausfcfg != "" {
		if err := factory.InitConfigFactory(config.ausfcfg); err != nil {
			return err
		}
	} else {
		DefaultAusfConfigPath := path_util.Free5gcPath("free5gc/config/ausfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultAusfConfigPath); err != nil {
			return err
		}
	}

	ausf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (ausf *AUSF) setLogLevel() {
	if factory.AusfConfig.Logger == nil {
		initLog.Warnln("AUSF config without log level setting!!!")
		return
	}

	if factory.AusfConfig.Logger.AUSF != nil {
		if factory.AusfConfig.Logger.AUSF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.AusfConfig.Logger.AUSF.DebugLevel); err != nil {
				initLog.Warnf("AUSF Log level [%s] is invalid, set to [info] level",
					factory.AusfConfig.Logger.AUSF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("AUSF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Warnln("AUSF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.AusfConfig.Logger.AUSF.ReportCaller)
	}

	if factory.AusfConfig.Logger.PathUtil != nil {
		if factory.AusfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.AusfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.AusfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.AusfConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.AusfConfig.Logger.OpenApi != nil {
		if factory.AusfConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.AusfConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.AusfConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.AusfConfig.Logger.OpenApi.ReportCaller)
	}
}

func (ausf *AUSF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ausf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ausf *AUSF) Start() {
	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)
	ueauthentication.AddService(router)

	ausf_context.Init()
	self := ausf_context.GetSelf()
	// Register to NRF
	profile, err := consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build AUSF Profile Error")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("AUSF register to NRF Error[%s]", err.Error())
	}

	ausfLogPath := util.AusfLogPath

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		ausf.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, ausfLogPath, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.AusfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.AusfPemPath, util.AusfKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (ausf *AUSF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("ausfcfg"))
	args := ausf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ausf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
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
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		startErr := command.Start()
		if startErr != nil {
			initLog.Fatalln(startErr)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (ausf *AUSF) Terminate() {
	logger.InitLog.Infof("Terminating AUSF...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}

	logger.InitLog.Infof("AUSF terminated")
}
