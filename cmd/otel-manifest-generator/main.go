package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"otel-manifest-generator/pkg/config"
	"otel-manifest-generator/pkg/generator"
	"otel-manifest-generator/pkg/upstream"
)

func main() {
	outputFile := flag.String("output", "builder-config.yaml", "Path to the generated builder configuration file")
	otelVersion := flag.String("otel-version", "0.119.0", "Version of the OTel Collector Distribution to use (overridden by upstream manifest)")
	flag.Parse()

	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		fmt.Println("Usage: otel-manifest-generator [flags] <config-file-1> <config-file-2> ...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. Fetch Upstream Manifest
	log.Println("Fetching upstream manifest...")
	upstreamManifest, err := upstream.Fetch(upstream.DefaultManifestURL)
	if err != nil {
		log.Fatalf("Error fetching upstream manifest: %v", err)
	}

	receivers, processors, exporters, extensions := upstream.PopulateComponents(upstreamManifest)
	log.Printf("Loaded %d receivers, %d processors, %d exporters, %d extensions from upstream.",
		len(receivers), len(processors), len(exporters), len(extensions))

	// 2. Parse and Merge OTel Configs
	log.Printf("Parsing configuration files: %v", inputFiles)
	otelConfig, err := config.Parse(inputFiles)
	if err != nil {
		log.Fatalf("Error parsing configuration files: %v", err)
	}

	// 3. Generate Manifest
	gen := generator.New(receivers, processors, exporters, extensions, upstreamManifest.Dist.Version)
	if *otelVersion != "0.119.0" && *otelVersion != upstreamManifest.Dist.Version {
		log.Printf("Warning: Using upstream version %s instead of requested %s", upstreamManifest.Dist.Version, *otelVersion)
	}

	err = gen.Generate(otelConfig, *outputFile)
	if err != nil {
		log.Fatalf("Error generating manifest: %v", err)
	}

	fmt.Printf("Successfully generated %s based on input configurations (OTel Version: %s)\n", *outputFile, upstreamManifest.Dist.Version)
}
