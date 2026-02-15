package upstream

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultManifestURL = "https://raw.githubusercontent.com/open-telemetry/opentelemetry-collector-releases/refs/heads/main/distributions/otelcol-contrib/manifest.yaml"

// UpstreamManifest represents the structure of the upstream manifest file
type UpstreamManifest struct {
	Dist struct {
		Version string `yaml:"version"`
	} `yaml:"dist"`
	Receivers  []Component `yaml:"receivers"`
	Processors []Component `yaml:"processors"`
	Exporters  []Component `yaml:"exporters"`
	Extensions []Component `yaml:"extensions"`
}

// Component represents a single OTel component
type Component struct {
	GoMod string `yaml:"gomod"`
}

// Fetch downloads and parses the upstream manifest from the given URL.
func Fetch(url string) (*UpstreamManifest, error) {
	if url == "" {
		url = DefaultManifestURL
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest UpstreamManifest
	err = yaml.Unmarshal(body, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

// PopulateComponents extracts component mappings from the manifest.
// It returns a map where keys are component types (e.g., "otlp") and values are full gomod lines.
func PopulateComponents(manifest *UpstreamManifest) (receivers, processors, exporters, extensions map[string]string) {
	receivers = make(map[string]string)
	processors = make(map[string]string)
	exporters = make(map[string]string)
	extensions = make(map[string]string)

	populateMap(manifest.Receivers, receivers, "receiver")
	populateMap(manifest.Processors, processors, "processor")
	populateMap(manifest.Exporters, exporters, "exporter")
	populateMap(manifest.Extensions, extensions, "extension")

	return
}

func populateMap(components []Component, targetMap map[string]string, suffix string) {
	for _, comp := range components {
		parts := strings.Split(comp.GoMod, " ")
		if len(parts) < 1 {
			continue
		}
		modPath := parts[0]
		base := path.Base(modPath)
		name := strings.TrimSuffix(base, suffix)
		targetMap[name] = comp.GoMod
	}
}
