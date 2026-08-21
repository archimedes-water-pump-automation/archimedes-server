# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- MQTT stream consumers that ingest water tank sensor readings and pump start/stop events.
- Volume calculation for cylindrical-cone tanks, with dimensions validated against a JSON Schema.
- Read-only HTTP API: health check, tank list/get-by-id, pump list/get-by-id, and pump run history.
- Graceful shutdown on `SIGINT`/`SIGTERM`, draining both MQTT consumers before exit.

### Fixed

- Pump and tank list endpoints returning an incorrect empty response.
- MQTT broker address no longer hard-coded; now read from `MQTT_BROKER_URL`.

[Unreleased]: https://github.com/archimedes-water-pump-automation/archimedes-server/commits/main
