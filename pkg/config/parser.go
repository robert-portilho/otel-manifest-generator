package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// rawOTelConfig is used for initial unmarshaling to handle flexible types
type rawOTelConfig struct {
	Receivers  interface{} `yaml:"receivers"`
	Processors interface{} `yaml:"processors"`
	Exporters  interface{} `yaml:"exporters"`
	Extensions interface{} `yaml:"extensions"`
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

		var rawConfig rawOTelConfig
		if err := yaml.Unmarshal(data, &rawConfig); err != nil {
			return nil, err
		}

		// Helper to resolve potentially external config
		resolve := func(raw interface{}) (map[string]interface{}, error) {
			if raw == nil {
				return make(map[string]interface{}), nil
			}

			// case 1: it's a map
			if m, ok := raw.(map[string]interface{}); ok {
				return m, nil
			}

			// case 2: it's a string (external file)
			if s, ok := raw.(string); ok {
				return loadExternalFile(s, file)
			}

			return nil, fmt.Errorf("unexpected type for component config: %T", raw)
		}

		receivers, err := resolve(rawConfig.Receivers)
		if err != nil {
			return nil, fmt.Errorf("error resolving receivers: %w", err)
		}
		processors, err := resolve(rawConfig.Processors)
		if err != nil {
			return nil, fmt.Errorf("error resolving processors: %w", err)
		}
		exporters, err := resolve(rawConfig.Exporters)
		if err != nil {
			return nil, fmt.Errorf("error resolving exporters: %w", err)
		}
		extensions, err := resolve(rawConfig.Extensions)
		if err != nil {
			return nil, fmt.Errorf("error resolving extensions: %w", err)
		}

		// Merge maps
		mergeMap(mergedConfig.Receivers, receivers)
		mergeMap(mergedConfig.Processors, processors)
		mergeMap(mergedConfig.Exporters, exporters)
		mergeMap(mergedConfig.Extensions, extensions)

		// Copy service config (simple overwrite for now as in original,
		// could technically merge pipelines too but kept simple)
		mergedConfig.Service = rawConfig.Service
	}

	return mergedConfig, nil
}

func loadExternalFile(ref string, parentFile string) (map[string]interface{}, error) {
	const prefix = "${file:"
	const suffix = "}"

	if !strings.HasPrefix(ref, prefix) || !strings.HasSuffix(ref, suffix) {
		return nil, fmt.Errorf("invalid file reference format: %s", ref)
	}

	pathStr := ref[len(prefix) : len(ref)-len(suffix)]

	// resolve path relative to parent file
	parentDir := filepath.Dir(parentFile)
	absPath := filepath.Join(parentDir, pathStr)

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func mergeMap(dest, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
