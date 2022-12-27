package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/n3iwf/logger"
	"github.com/free5gc/n3iwf/service"
	"github.com/free5gc/version"
)

var N3IWF = &service.N3IWF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "n3iwf"
	appLog.Infoln(app.Name)
	appLog.Infoln("N3IWF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -n3iwfcfg n3iwf configuration file"
	app.Action = action
	app.Flags = N3IWF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("N3IWF Run Error: %v", err)
	}
}

func action(c *cli.Context) error {
	if err := N3IWF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	N3IWF.Start()

	return nil
}
