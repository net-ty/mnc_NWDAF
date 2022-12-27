/*
 * NSSF Service
 */

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

	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/nssf/consumer"
	"github.com/free5gc/nssf/context"
	"github.com/free5gc/nssf/factory"
	"github.com/free5gc/nssf/logger"
	"github.com/free5gc/nssf/nssaiavailability"
	"github.com/free5gc/nssf/nsselection"
	"github.com/free5gc/nssf/util"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
)

type NSSF struct{}

type (
	// Config information.
	Config struct {
		nssfcfg string
	}
)

var config Config

var nssfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nssfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NSSF) GetCliCmd() (flags []cli.Flag) {
	return nssfCLi
}

func (nssf *NSSF) Initialize(c *cli.Context) error {
	config = Config{
		nssfcfg: c.String("nssfcfg"),
	}

	if config.nssfcfg != "" {
		if err := factory.InitConfigFactory(config.nssfcfg); err != nil {
			return err
		}
	} else {
		DefaultNssfConfigPath := path_util.Free5gcPath("free5gc/config/nssfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultNssfConfigPath); err != nil {
			return err
		}
	}

	context.InitNssfContext()

	nssf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nssf *NSSF) setLogLevel() {
	if factory.NssfConfig.Logger == nil {
		initLog.Warnln("NSSF config without log level setting!!!")
		return
	}

	if factory.NssfConfig.Logger.NSSF != nil {
		if factory.NssfConfig.Logger.NSSF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NssfConfig.Logger.NSSF.DebugLevel); err != nil {
				initLog.Warnf("NSSF Log level [%s] is invalid, set to [info] level",
					factory.NssfConfig.Logger.NSSF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("NSSF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("NSSF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NssfConfig.Logger.NSSF.ReportCaller)
	}

	if factory.NssfConfig.Logger.PathUtil != nil {
		if factory.NssfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NssfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.NssfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.NssfConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.NssfConfig.Logger.OpenApi != nil {
		if factory.NssfConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NssfConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.NssfConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.NssfConfig.Logger.OpenApi.ReportCaller)
	}
}

func (nssf *NSSF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nssf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nssf *NSSF) Start() {
	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	nssaiavailability.AddService(router)
	nsselection.AddService(router)

	self := context.NSSF_Self()
	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	// Register to NRF
	profile, err := consumer.BuildNFProfile(self)
	if err != nil {
		initLog.Error("Failed to build NSSF profile")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("Failed to register NSSF to NRF: %s", err.Error())
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		nssf.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, util.NSSF_LOG_PATH, router)

	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.NssfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.NSSF_PEM_PATH, util.NSSF_KEY_PATH)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (nssf *NSSF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("nssfcfg"))
	args := nssf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nssf", args...)

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
		if err = command.Start(); err != nil {
			fmt.Printf("NSSF Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (nssf *NSSF) Terminate() {
	logger.InitLog.Infof("Terminating NSSF...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}

	logger.InitLog.Infof("NSSF terminated")
}
