package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// OTelConfig represents the structure of the OTel Collector configuration file
type OTelConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Extensions map[string]interface{} `yaml:"extensions"`
	Service    struct {
		Extensions []string `yaml:"extensions"`
		Pipelines  map[string]struct {
			Receivers  []string `yaml:"receivers"`
			Processors []string `yaml:"processors"`
			Exporters  []string `yaml:"exporters"`
		} `yaml:"pipelines"`
	} `yaml:"service"`
}

// Parse reads one or more OTel configuration files and merges them into a single OTelConfig.
func Parse(files []string) (*OTelConfig, error) {
	mergedConfig := &OTelConfig{
		Receivers:  make(map[string]interface{}),
		Processors: make(map[string]interface{}),
		Exporters:  make(map[string]interface{}),
		Extensions: make(map[string]interface{}),
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var config OTelConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		// Merge maps
		mergeMap(mergedConfig.Receivers, config.Receivers)
		mergeMap(mergedConfig.Processors, config.Processors)
		mergeMap(mergedConfig.Exporters, config.Exporters)
		mergeMap(mergedConfig.Extensions, config.Extensions)
	}

	return mergedConfig, nil
}

func mergeMap(dest, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
