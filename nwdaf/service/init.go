package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	"github.com/free5gc/path_util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"nwdaf.com/anlf"
	"nwdaf.com/consumer"
	nwdaf_context "nwdaf.com/context"
	"nwdaf.com/factory"
	"nwdaf.com/logger"
	"nwdaf.com/mtlf"
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

	initLog.Infoln("Server started")

	if !util.InitNWDAFContext() {
		initLog.Error("Initicating context failed")
		return
	}

	wg := sync.WaitGroup{}

	self := nwdaf_context.NWDAF_Self()
	util.InitNwdafContext(self)

	addr := fmt.Sprintf("127.0.0.1:24242")
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	mtlf.AddService(router)
	anlf.AddService(router)

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

	server, err := http2_util.NewServer(addr, "nwdafsslkey.log", router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}
	if err != nil {
		initLog.Warnln("Initialize HTTP server:", err)
	}
	serverScheme := factory.NwdafConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServe() //TODO: changing to HTTPS (TLS)
	}

	if err != nil {
		initLog.Fatalln("HTTP server setup failed:", err)
	}
	initLog.Info("NWDAF running...")

	wg.Wait()
}

func (nwdaf *NWDAF) Exec(c *cli.Context) error {

	initLog.Traceln("args:", c.String("nwdafcfg"))
	args := nwdaf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nwdaf", args...)

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
			initLog.Errorf("NWDAF start error: %v", errCom)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
