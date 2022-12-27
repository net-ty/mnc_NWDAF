package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/MongoDBLibrary"
	mongoDBLibLogger "github.com/free5gc/MongoDBLibrary/logger"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
	"github.com/free5gc/udr/consumer"
	udr_context "github.com/free5gc/udr/context"
	"github.com/free5gc/udr/datarepository"
	"github.com/free5gc/udr/factory"
	"github.com/free5gc/udr/logger"
	"github.com/free5gc/udr/util"
)

type UDR struct{}

type (
	// Config information.
	Config struct {
		udrcfg string
	}
)

var config Config

var udrCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "udrcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*UDR) GetCliCmd() (flags []cli.Flag) {
	return udrCLi
}

func (udr *UDR) Initialize(c *cli.Context) error {
	config = Config{
		udrcfg: c.String("udrcfg"),
	}

	if config.udrcfg != "" {
		if err := factory.InitConfigFactory(config.udrcfg); err != nil {
			return err
		}
	} else {
		DefaultUdrConfigPath := path_util.Free5gcPath("free5gc/config/udrcfg.yaml")
		if err := factory.InitConfigFactory(DefaultUdrConfigPath); err != nil {
			return err
		}
	}

	udr.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (udr *UDR) setLogLevel() {
	if factory.UdrConfig.Logger == nil {
		initLog.Warnln("UDR config without log level setting!!!")
		return
	}

	if factory.UdrConfig.Logger.UDR != nil {
		if factory.UdrConfig.Logger.UDR.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.UdrConfig.Logger.UDR.DebugLevel); err != nil {
				initLog.Warnf("UDR Log level [%s] is invalid, set to [info] level",
					factory.UdrConfig.Logger.UDR.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("UDR Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("UDR Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.UdrConfig.Logger.UDR.ReportCaller)
	}

	if factory.UdrConfig.Logger.PathUtil != nil {
		if factory.UdrConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.UdrConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.UdrConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.UdrConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.UdrConfig.Logger.OpenApi != nil {
		if factory.UdrConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.UdrConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.UdrConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.UdrConfig.Logger.OpenApi.ReportCaller)
	}

	if factory.UdrConfig.Logger.MongoDBLibrary != nil {
		if factory.UdrConfig.Logger.MongoDBLibrary.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.UdrConfig.Logger.MongoDBLibrary.DebugLevel); err != nil {
				mongoDBLibLogger.MongoDBLog.Warnf("MongoDBLibrary Log level [%s] is invalid, set to [info] level",
					factory.UdrConfig.Logger.MongoDBLibrary.DebugLevel)
				mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				mongoDBLibLogger.SetLogLevel(level)
			}
		} else {
			mongoDBLibLogger.MongoDBLog.Warnln("MongoDBLibrary Log level not set. Default set to [info] level")
			mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
		}
		mongoDBLibLogger.SetReportCaller(factory.UdrConfig.Logger.MongoDBLibrary.ReportCaller)
	}
}

func (udr *UDR) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range udr.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (udr *UDR) Start() {
	// get config file info
	config := factory.UdrConfig
	mongodb := config.Configuration.Mongodb

	initLog.Infof("UDR Config Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	// Connect to MongoDB
	MongoDBLibrary.SetMongoDB(mongodb.Name, mongodb.Url)

	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	datarepository.AddService(router)

	udrLogPath := util.UdrLogPath
	udrPemPath := util.UdrPemPath
	udrKeyPath := util.UdrKeyPath

	self := udr_context.UDR_Self()
	util.InitUdrContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
	profile := consumer.BuildNFInstance(self)
	var newNrfUri string
	var err error
	newNrfUri, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, profile.NfInstanceId, profile)
	if err == nil {
		self.NrfUri = newNrfUri
	} else {
		initLog.Errorf("Send Register NFInstance Error[%s]", err.Error())
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		udr.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, udrLogPath, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := factory.UdrConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(udrPemPath, udrKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (udr *UDR) Exec(c *cli.Context) error {
	// UDR.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("udrcfg"))
	args := udr.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./udr", args...)

	if err := udr.Initialize(c); err != nil {
		return err
	}

	var stdout io.ReadCloser
	if readCloser, err := command.StdoutPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stdout = readCloser
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

	var stderr io.ReadCloser
	if readCloser, err := command.StderrPipe(); err != nil {
		initLog.Fatalln(err)
	} else {
		stderr = readCloser
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	var err error
	go func() {
		if errormessage := command.Start(); err != nil {
			fmt.Println("command.Start Fails!")
			err = errormessage
		}
		wg.Done()
	}()

	wg.Wait()
	return err
}

func (udr *UDR) Terminate() {
	logger.InitLog.Infof("Terminating UDR...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
	logger.InitLog.Infof("UDR terminated")
}
