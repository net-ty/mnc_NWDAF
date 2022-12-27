/*
 * AUSF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/ausf/logger"
)

var AusfConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		AusfConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &AusfConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := AusfConfig.GetVersion()

	if currentVersion != AUSF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, AUSF_EXPECTED_CONFIG_VERSION)
	}
	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
