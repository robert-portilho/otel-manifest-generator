package generator

import (
	"os"
	"path/filepath"
	"testing"

	"otel-manifest-generator/pkg/config"
)

func TestGenerate(t *testing.T) {
	// Setup temporary output directory
	tmpDir, err := os.MkdirTemp("", "generator-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "builder-config.yaml")

	// Sample Config
	cfg := &config.OTelConfig{
		Receivers:  map[string]interface{}{"otlp": nil},
		Processors: map[string]interface{}{"batch": nil},
		Exporters:  map[string]interface{}{"otlp": nil},
		Extensions: map[string]interface{}{},
	}

	// Sample Component Maps
	receivers := map[string]string{"otlp": "go.opentelemetry.io/collector/receiver/otlpreceiver v0.92.0"}
	processors := map[string]string{"batch": "go.opentelemetry.io/collector/processor/batchprocessor v0.92.0"}
	exporters := map[string]string{"otlp": "go.opentelemetry.io/collector/exporter/otlpexporter v0.92.0"}
	extensions := map[string]string{}

	gen := New(receivers, processors, exporters, extensions, "0.145.0")
	err = gen.Generate(cfg, outputPath)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Generate() did not create output file")
	}
}

func TestResolveComponents(t *testing.T) {
	gen := &Generator{
		ReceiverMap: map[string]string{
			"otlp": "go.opentelemetry.io/collector/receiver/otlpreceiver v0.92.0",
		},
	}

	components := map[string]interface{}{
		"otlp":   nil,
		"otlp/1": nil,
	}

	resolved := gen.resolveComponents(components, gen.ReceiverMap)

	if len(resolved) != 1 {
		t.Errorf("Expected 1 resolved component (deduplicated), got %d", len(resolved))
	}
}

func TestResolveComponents_Underscore(t *testing.T) {
	gen := &Generator{
		ProcessorMap: map[string]string{
			"memorylimiter": "go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.92.0",
		},
	}

	components := map[string]interface{}{
		"memory_limiter": nil,
	}

	resolved := gen.resolveComponents(components, gen.ProcessorMap)

	if len(resolved) != 1 {
		t.Errorf("Expected 1 resolved component for memory_limiter, got %d", len(resolved))
	}
}
