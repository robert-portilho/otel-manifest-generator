package generator

import (
	"fmt"
	"log"
	"os"

	"otel-manifest-generator/pkg/config"

	"gopkg.in/yaml.v3"
)

// BuilderConfig represents the structure of the OCB builder-config.yaml
type BuilderConfig struct {
	Dist struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Output      string `yaml:"output_path"`
		OtelCol     string `yaml:"otelcol_version"`
	} `yaml:"dist"`
	Exporters  []Component `yaml:"exporters,omitempty"`
	Extensions []Component `yaml:"extensions,omitempty"`
	Processors []Component `yaml:"processors,omitempty"`
	Receivers  []Component `yaml:"receivers,omitempty"`
	Replaces   []string    `yaml:"replaces,omitempty"`
}

// Component represents a single OTel component in the builder config
type Component struct {
	GoMod string `yaml:"gomod"`
}

type Generator struct {
	ReceiverMap  map[string]string
	ProcessorMap map[string]string
	ExporterMap  map[string]string
	ExtensionMap map[string]string
	Version      string
}

func New(receivers, processors, exporters, extensions map[string]string, version string) *Generator {
	return &Generator{
		ReceiverMap:  receivers,
		ProcessorMap: processors,
		ExporterMap:  exporters,
		ExtensionMap: extensions,
		Version:      version,
	}
}

func (g *Generator) Generate(cfg *config.OTelConfig, outputPath string) error {
	builderConfig := BuilderConfig{}
	builderConfig.Dist.Name = "otel-custom-col"
	builderConfig.Dist.Description = "Custom OpenTelemetry Collector build from config"
	builderConfig.Dist.Output = "./otelcol-custom"
	builderConfig.Dist.OtelCol = g.Version

	builderConfig.Receivers = g.resolveComponents(cfg.Receivers, g.ReceiverMap)
	builderConfig.Processors = g.resolveComponents(cfg.Processors, g.ProcessorMap)
	builderConfig.Exporters = g.resolveComponents(cfg.Exporters, g.ExporterMap)
	builderConfig.Extensions = g.resolveComponents(cfg.Extensions, g.ExtensionMap)

	outData, err := yaml.Marshal(&builderConfig)
	if err != nil {
		return fmt.Errorf("error marshaling builder config: %w", err)
	}

	if err := os.WriteFile(outputPath, outData, 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}

func (g *Generator) resolveComponents(components map[string]interface{}, componentMap map[string]string) []Component {
	var manifestComponents []Component
	seen := make(map[string]bool)

	for name := range components {
		compType := name
		for i, r := range name {
			if r == '/' {
				compType = name[:i]
				break
			}
		}

		if seen[compType] {
			continue
		}
		seen[compType] = true

		gomod, ok := componentMap[compType]
		if !ok {
			log.Printf("Warning: Component type '%s' not found in upstream manifest. Skipping.", compType)
			continue
		}
		manifestComponents = append(manifestComponents, Component{
			GoMod: gomod,
		})
	}
	return manifestComponents
}
