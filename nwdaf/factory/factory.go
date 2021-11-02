/*
 * NWDAF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"nwdaf.com/logger"
)

var NwdafConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		NwdafConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NwdafConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := NwdafConfig.GetVersion()

	if currentVersion != NWDAF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, NWDAF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
