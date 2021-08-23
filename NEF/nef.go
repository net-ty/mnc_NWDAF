package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/version"
	"nef.com/logger"
	"nef.com/service"
)

// func main() {
// 	message := service.Hello("servicehello")
// 	fmt.Println(message)
// 	message2 := logger.Hello("loggerhello")
// 	fmt.Println(message2)
// 	message3 := factory.Hello("factoryhello")
// 	fmt.Println(message3)

// }

var NEF = &service.NEF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "nef"
	appLog.Infoln(app.Name)
	appLog.Infoln("NEF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -nefcfg nef configuration file"
	app.Action = action
	app.Flags = NEF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("NEF Run error: %v", err)
		return
	}
}

func action(c *cli.Context) error {
	if err := NEF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	NEF.Start()

	return nil
}
