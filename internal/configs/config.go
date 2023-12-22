package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/tranHieuDev23/cato/configs"
)

type ConfigFilePath string

type Config struct {
	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
	HTTP     HTTP     `yaml:"http"`
}

func NewConfig(filePath ConfigFilePath) (config Config, err error) {
	configBytes := configs.DefaultConfigBytes

	if filePath != "" {
		configBytes, err = os.ReadFile(string(filePath))
		if err != nil {
			return Config{}, fmt.Errorf("failed to read YAML file: %v", err)
		}
	}

	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return config, nil
}
