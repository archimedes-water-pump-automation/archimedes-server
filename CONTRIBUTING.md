# Contributing to archimedes-server

Thank you for your interest in contributing!

## Prerequisites

- Go 1.25 or later
- A PostgreSQL instance and an MQTT broker for running the service end-to-end (not required to run the unit tests)

## Quick Start

```bash
# Clone the repository
git clone https://github.com/archimedes-water-pump-automation/archimedes-server.git
cd archimedes-server

# Build
go build -o archimedes-server .

# Run unit tests
go test ./...
```

## Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Run `go test ./...` and `go vet ./...`
6. Commit with a descriptive message
7. Push and open a Pull Request

## Code Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Add doc comments to all exported symbols
- Keep the hexagonal layout: `core/*` holds domain types, interfaces, and use cases with no dependency on a concrete adapter; `adapters/*` implements those interfaces against PostgreSQL, MQTT, HTTP, etc.
- Adding a new tank shape means implementing `core/volume/interfaces.IVolumeCalculator` and registering it in `adapters/database/postgresql.getVolumeType.GetVolumeFromShape`

## Reporting Issues

Use [GitHub Issues](https://github.com/archimedes-water-pump-automation/archimedes-server/issues). Include:

- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
