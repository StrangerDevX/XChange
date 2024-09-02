package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Token string `yaml:"token"`
}

var configFileName = "config.yaml"

func LoadConfig() (*Config, error) {
	config := &Config{}
	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file %s not found", configFileName)
		}
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func CreateConfigFile() error {
	_, err := os.ReadFile(configFileName)
	if err == nil {
		fmt.Println("Config file already exists")
		return nil
	}
	config := &Config{}
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFileName, out, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Created config file " + configFileName)
	return nil
}

func SetToken(token string) error {
	config := &Config{}
	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file %s not found", configFileName)
		}
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return err
	}

	config.Token = token

	newConfigFile, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFileName, newConfigFile, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Token updated")
	return nil

}
