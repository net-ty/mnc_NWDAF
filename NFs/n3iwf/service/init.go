package service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	aperLogger "github.com/free5gc/aper/logger"
	"github.com/free5gc/n3iwf/factory"
	ike_service "github.com/free5gc/n3iwf/ike/service"
	"github.com/free5gc/n3iwf/logger"
	ngap_service "github.com/free5gc/n3iwf/ngap/service"
	nwucp_service "github.com/free5gc/n3iwf/nwucp/service"
	nwuup_service "github.com/free5gc/n3iwf/nwuup/service"
	"github.com/free5gc/n3iwf/util"
	ngapLogger "github.com/free5gc/ngap/logger"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
)

type N3IWF struct{}

type (
	// Config information.
	Config struct {
		n3iwfcfg string
	}
)

var config Config

var n3iwfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "n3iwfcfg",
		Usage: "n3iwf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*N3IWF) GetCliCmd() (flags []cli.Flag) {
	return n3iwfCLi
}

func (n3iwf *N3IWF) Initialize(c *cli.Context) error {
	config = Config{
		n3iwfcfg: c.String("n3iwfcfg"),
	}

	if config.n3iwfcfg != "" {
		if err := factory.InitConfigFactory(config.n3iwfcfg); err != nil {
			return err
		}
	} else {
		DefaultN3iwfConfigPath := path_util.Free5gcPath("free5gc/config/n3iwfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultN3iwfConfigPath); err != nil {
			return err
		}
	}

	n3iwf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (n3iwf *N3IWF) setLogLevel() {
	if factory.N3iwfConfig.Logger == nil {
		initLog.Warnln("N3IWF config without log level setting!!!")
		return
	}

	if factory.N3iwfConfig.Logger.N3IWF != nil {
		if factory.N3iwfConfig.Logger.N3IWF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.N3iwfConfig.Logger.N3IWF.DebugLevel); err != nil {
				initLog.Warnf("N3IWF Log level [%s] is invalid, set to [info] level",
					factory.N3iwfConfig.Logger.N3IWF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("N3IWF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("N3IWF Log level is default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.N3iwfConfig.Logger.N3IWF.ReportCaller)
	}

	if factory.N3iwfConfig.Logger.NGAP != nil {
		if factory.N3iwfConfig.Logger.NGAP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.N3iwfConfig.Logger.NGAP.DebugLevel); err != nil {
				ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
					factory.N3iwfConfig.Logger.NGAP.DebugLevel)
				ngapLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				ngapLogger.SetLogLevel(level)
			}
		} else {
			ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
			ngapLogger.SetLogLevel(logrus.InfoLevel)
		}
		ngapLogger.SetReportCaller(factory.N3iwfConfig.Logger.NGAP.ReportCaller)
	}

	if factory.N3iwfConfig.Logger.Aper != nil {
		if factory.N3iwfConfig.Logger.Aper.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.N3iwfConfig.Logger.Aper.DebugLevel); err != nil {
				aperLogger.AperLog.Warnf("Aper Log level [%s] is invalid, set to [info] level",
					factory.N3iwfConfig.Logger.Aper.DebugLevel)
				aperLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				aperLogger.SetLogLevel(level)
			}
		} else {
			aperLogger.AperLog.Warnln("Aper Log level not set. Default set to [info] level")
			aperLogger.SetLogLevel(logrus.InfoLevel)
		}
		aperLogger.SetReportCaller(factory.N3iwfConfig.Logger.Aper.ReportCaller)
	}

	if factory.N3iwfConfig.Logger.PathUtil != nil {
		if factory.N3iwfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.N3iwfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.N3iwfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.N3iwfConfig.Logger.PathUtil.ReportCaller)
	}
}

func (n3iwf *N3IWF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range n3iwf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (n3iwf *N3IWF) Start() {
	initLog.Infoln("Server started")

	if !util.InitN3IWFContext() {
		initLog.Error("Initicating context failed")
		return
	}

	wg := sync.WaitGroup{}

	// NGAP
	if err := ngap_service.Run(); err != nil {
		initLog.Errorf("Start NGAP service failed: %+v", err)
		return
	}
	initLog.Info("NGAP service running.")
	wg.Add(1)

	// Relay listeners
	// Control plane
	if err := nwucp_service.Run(); err != nil {
		initLog.Errorf("Listen NWu control plane traffic failed: %+v", err)
		return
	}
	initLog.Info("NAS TCP server successfully started.")
	wg.Add(1)

	// User plane
	if err := nwuup_service.Run(); err != nil {
		initLog.Errorf("Listen NWu user plane traffic failed: %+v", err)
		return
	}
	initLog.Info("Listening NWu user plane traffic")
	wg.Add(1)

	// IKE
	if err := ike_service.Run(); err != nil {
		initLog.Errorf("Start IKE service failed: %+v", err)
		return
	}
	initLog.Info("IKE service running.")
	wg.Add(1)

	initLog.Info("N3IWF running...")

	wg.Wait()
}

func (n3iwf *N3IWF) Exec(c *cli.Context) error {
	// N3IWF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("n3iwfcfg"))
	args := n3iwf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./n3iwf", args...)

	wg := sync.WaitGroup{}
	wg.Add(3)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
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
		if errCom := command.Start(); errCom != nil {
			initLog.Errorf("N3IWF start error: %v", errCom)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
