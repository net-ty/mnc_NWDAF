/*
 * UDM Configuration Factory
 */

package factory

import (
	"github.com/free5gc/logger_util"
)

const (
	UDM_EXPECTED_CONFIG_VERSION = "1.0.0"
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
	UDM_DEFAULT_IPV4     = "127.0.0.3"
	UDM_DEFAULT_PORT     = "8000"
	UDM_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	UdmName         string   `yaml:"udmName,omitempty"`
	Sbi             *Sbi     `yaml:"sbi,omitempty"`
	ServiceNameList []string `yaml:"serviceNameList,omitempty"`
	NrfUri          string   `yaml:"nrfUri,omitempty"`
	Keys            *Keys    `yaml:"keys,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
	Tls         *Tls   `yaml:"tls,omitempty"`
}

type Tls struct {
	Log string `yaml:"log,omitempty"`
	Pem string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

type Keys struct {
	UdmProfileAHNPrivateKey string `yaml:"udmProfileAHNPrivateKey,omitempty"`
	UdmProfileAHNPublicKey  string `yaml:"udmProfileAHNPublicKey,omitempty"`
	UdmProfileBHNPrivateKey string `yaml:"udmProfileBHNPrivateKey,omitempty"`
	UdmProfileBHNPublicKey  string `yaml:"udmProfileBHNPublicKey,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
