# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-08-21

### Added

- MQTT stream consumers that ingest water tank sensor readings and pump start/stop events.
- Volume calculation for cylindrical-cone tanks, with dimensions validated against a JSON Schema.
- Read-only HTTP API: health check, tank list/get-by-id, pump list/get-by-id, and pump run history.
- Graceful shutdown on `SIGINT`/`SIGTERM`, draining both MQTT consumers before exit.
- Test coverage for the HTTP handlers, stream processors, volume calculator/schema, and logger.
- `README.md`, `CONTRIBUTING.md`, `LICENSE`, and `llms.txt`.
- Godoc comments across all packages.

### Fixed

- Pump and tank list endpoints returning an incorrect empty response.
- MQTT broker address no longer hard-coded; now read from `MQTT_BROKER_URL`.
- Cylindrical-cone volume calculator using the raw sensor distance instead of the submerged height when sizing the cone's fluid surface, producing an incorrect volume once the fluid reached into the conical section.
- Cylindrical-cone calculator panicking instead of returning an error when a tank's stored dimensions had the wrong type.
- `incline_angle` not being bounded to a valid 0-90 degree range by the dimensions schema.

### Changed

- Reorganized the codebase into `adapters/database`, `adapters/endpoint`, `adapters/stream`, and per-domain `core/*/domain`, `core/*/interfaces`, `core/*/usecases` packages.
- HTTP handlers now write JSON responses directly instead of converting them to a string first.

[1.0.0]: https://github.com/archimedes-water-pump-automation/archimedes-server/releases/tag/v1.0.0
