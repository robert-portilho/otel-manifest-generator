package upstream

import (
	"testing"
)

func TestPopulateComponents(t *testing.T) {
	manifest := &UpstreamManifest{
		Receivers: []Component{
			{GoMod: "go.opentelemetry.io/collector/receiver/otlpreceiver v0.92.0"},
			{GoMod: "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.92.0"},
		},
		Processors: []Component{
			{GoMod: "go.opentelemetry.io/collector/processor/batchprocessor v0.92.0"},
		},
		Exporters: []Component{
			{GoMod: "go.opentelemetry.io/collector/exporter/otlpexporter v0.92.0"},
		},
		Extensions: []Component{
			{GoMod: "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.92.0"},
		},
	}

	receivers, processors, exporters, extensions := PopulateComponents(manifest)

	// Check Receivers
	if _, ok := receivers["otlp"]; !ok {
		t.Error("Expected 'otlp' receiver")
	}
	if _, ok := receivers["hostmetrics"]; !ok {
		t.Error("Expected 'hostmetrics' receiver")
	}

	// Check Processors
	if _, ok := processors["batch"]; !ok {
		t.Error("Expected 'batch' processor")
	}

	// Check Exporters
	if _, ok := exporters["otlp"]; !ok {
		t.Error("Expected 'otlp' exporter")
	}

	// Check Extensions
	if _, ok := extensions["healthcheck"]; !ok {
		t.Error("Expected 'healthcheck' extension")
	}
}
