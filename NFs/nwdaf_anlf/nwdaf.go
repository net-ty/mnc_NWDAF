package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

//	"github.com/free5gc/version"
	"nwdaf.com/logger"
	"nwdaf.com/service"
)

// func main() {
// 	message := service.Hello("servicehello")
// 	fmt.Println(message)
// 	message2 := logger.Hello("loggerhello")
// 	fmt.Println(message2)
// 	message3 := factory.Hello("factoryhello")
// 	fmt.Println(message3)

// }

var NWDAF = &service.NWDAF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "nwdaf"
	appLog.Infoln(app.Name)
//	appLog.Infoln("NWDAF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -nwdafcfg nwdaf configuration file"
	app.Action = action
	app.Flags = NWDAF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("NWDAF Run error: %v", err)
		return
	}
}

func action(c *cli.Context) error {
	if err := NWDAF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	NWDAF.Start()

	return nil
}
