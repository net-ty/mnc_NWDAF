/*
 * PCF Configuration Factory
 */

package factory

import (
	"github.com/free5gc/logger_util"
)

const (
	PCF_EXPECTED_CONFIG_VERSION = "1.0.0"
)

type Config struct {
	Info          *Info               `yaml:"info"`
	Configuration *Configuration      `yaml:"configuration"`
	Logger        *logger_util.Logger `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

const (
	PCF_DEFAULT_IPV4     = "127.0.0.7"
	PCF_DEFAULT_PORT     = "8000"
	PCF_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	PcfName         string    `yaml:"pcfName,omitempty"`
	Sbi             *Sbi      `yaml:"sbi,omitempty"`
	TimeFormat      string    `yaml:"timeFormat,omitempty"`
	DefaultBdtRefId string    `yaml:"defaultBdtRefId,omitempty"`
	NrfUri          string    `yaml:"nrfUri,omitempty"`
	ServiceList     []Service `yaml:"serviceList,omitempty"`
	Mongodb         *Mongodb  `yaml:"mongodb"`
}

type Service struct {
	ServiceName string `yaml:"serviceName"`
	SuppFeat    string `yaml:"suppFeat,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	// IPv6Addr  string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
}

type Mongodb struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
