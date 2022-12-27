//+build debug

package util

import (
	"github.com/free5gc/path_util"
)

var (
	UdrLogPath           = path_util.Free5gcPath("free5gc/udrsslkey.log")
	UdrPemPath           = path_util.Free5gcPath("free5gc/support/TLS/_debug.pem")
	UdrKeyPath           = path_util.Free5gcPath("free5gc/support/TLS/_debug.key")
	DefaultUdrConfigPath = path_util.Free5gcPath("free5gc/config/udrcfg.yaml")
)
