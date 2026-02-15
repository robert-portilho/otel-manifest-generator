# OTel Manifest Generator

A Go-based tool to generate an [OpenTelemetry Collector Builder (OCB)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) manifest (`builder-config.yaml`) from existing OpenTelemetry Collector configuration files.

## Features

-   **Dynamic Component Mapping**: Automatically fetches the latest component versions and mappings from the [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) releases.
-   **Multiple Config Support**: Merges multiple OTel configuration files into a single build manifest.
-   **CLI Arguments**: Flexible command-line interface for specifying input and output files.

## Installation

Ensure you have Go installed (1.21+ recommended).

```bash
git clone <repository-url>
cd otel-manifest-generator
go mod tidy
```

## Usage

Run the tool using `go run`:

```bash
go run cmd/otel-manifest-generator/main.go [flags] <config-file-1> <config-file-2> ...
```

### Flags

-   `-output`: Path to the generated builder configuration file (default: `builder-config.yaml`).
-   `-otel-version`: Desired OTel Collector version (note: currently overridden by the upstream manifest version to ensure compatibility).

### Examples

**Generate from a single config:**

```bash
go run cmd/otel-manifest-generator/main.go -output builder-config.yaml otel-config.yaml
```

**Generate from multiple configs (e.g., base and overrides):**

```bash
go run cmd/otel-manifest-generator/main.go -output dist/builder.yaml base-config.yaml production-override.yaml
```

## Project Structure

-   `cmd/otel-manifest-generator`: Main entry point applications.
-   `pkg/config`: Configuration parsing and merging logic.
-   `pkg/upstream`: Upstream manifest fetching and component mapping.
-   `pkg/generator`: OCB builder manifest generation logic.
