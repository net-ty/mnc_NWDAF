package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/free5gc/path_util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"nwdaf.com/consumer"
	nwdaf_context "nwdaf.com/context"
	"nwdaf.com/factory"
	"nwdaf.com/logger"
	"nwdaf.com/util"
)

type NWDAF struct{}

type (
	// config information
	Config struct {
		nwdafcfg string
	}
)

var config Config

var nwdafCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nwdafcfg",
		Usage: "nwdaf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NWDAF) GetCliCmd() (flags []cli.Flag) {
	return nwdafCLi
}

func (nwdaf *NWDAF) Initialize(c *cli.Context) error {
	config = Config{
		nwdafcfg: c.String("nwdafcfg"),
	}

	if config.nwdafcfg != "" {
		if err := factory.InitConfigFactory(config.nwdafcfg); err != nil {
			return err
		}
	} else {
		DefaultNWDAFConfigPath := path_util.Free5gcPath("free5gc/config/nwdafcfg.yaml")
		if err := factory.InitConfigFactory(DefaultNWDAFConfigPath); err != nil {
			return err
		}
	}

	nwdaf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nwdaf *NWDAF) setLogLevel() {
	if factory.NwdafConfig.Logger == nil {
		initLog.Warnln("NWDAF config without log level setting!!!")
		return
	}

	logger.SetLogLevel(logrus.InfoLevel)

}

func (nwdaf *NWDAF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nwdaf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nwdaf *NWDAF) Start() {

	initLog.Infoln("Registered to NRF")

	if !util.InitNWDAFContext() {
		initLog.Error("Initicating context failed")
		return
	}


	self := nwdaf_context.NWDAF_Self()
	util.InitNwdafContext(self)


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
		os.Exit(0)
	}()
}
