/*
 * NSSF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/nssf/logger"
)

var (
	NssfConfig Config
	Configured bool
	ConfigLock sync.RWMutex
)

func init() {
	Configured = false
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		NssfConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NssfConfig); yamlErr != nil {
			return yamlErr
		}

		Configured = true
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := NssfConfig.GetVersion()

	if currentVersion != NSSF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, NSSF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
