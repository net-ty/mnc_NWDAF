/*
 * NSSF Testing Utility
 */

package test

import (
	"flag"

	. "github.com/free5gc/nssf/plugin"
	"github.com/free5gc/path_util"
)

var (
	ConfigFileFromArgs string
	DefaultConfigFile  string = path_util.Free5gcPath("github.com/free5gc/nssf/test/conf/test_nssf_config.yaml")
)

type TestingUtil struct {
	ConfigFile string
}

type TestingNsselection struct {
	ConfigFile string

	GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

	GenerateRoamingQueryParameter func() NsselectionQueryParameter
}

type TestingNssaiavailability struct {
	ConfigFile string

	NfId string

	SubscriptionId string

	NfNssaiAvailabilityUri string
}

func init() {
	flag.StringVar(&ConfigFileFromArgs, "config-file", DefaultConfigFile, "Configuration file")
	flag.Parse()
}
