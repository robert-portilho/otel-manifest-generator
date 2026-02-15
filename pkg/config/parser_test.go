package config

import (
	"path/filepath"
	"testing"
)

func TestParse_Valid(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "valid.yaml"),
	}

	cfg, err := Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if _, ok := cfg.Receivers["otlp/1"]; !ok {
		t.Errorf("Expected receiver 'otlp/1' to be present")
	}
}

func TestParse_ExternalFile(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "external_main.yaml"),
	}

	cfg, err := Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check if processors from external file are loaded
	if _, ok := cfg.Processors["batch"]; !ok {
		t.Errorf("Expected processor 'batch' from processors.yaml")
	}
	if _, ok := cfg.Processors["memory_limiter"]; !ok {
		t.Errorf("Expected processor 'memory_limiter' from processors.yaml")
	}
}

func TestParse_Merge(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "valid.yaml"),
		filepath.Join("testdata", "valid2.yaml"),
	}

	cfg, err := Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check if elements from both files are present
	if _, ok := cfg.Receivers["otlp/1"]; !ok {
		t.Errorf("Expected receiver 'otlp/1' from valid.yaml")
	}
	if _, ok := cfg.Receivers["hostmetrics"]; !ok {
		t.Errorf("Expected receiver 'hostmetrics' from valid2.yaml")
	}
}

func TestParse_InvalidFile(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "non_existent.yaml"),
	}

	_, err := Parse(files)
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "invalid.yaml"),
	}

	_, err := Parse(files)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
}
