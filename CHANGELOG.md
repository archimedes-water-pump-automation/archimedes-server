# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `MQTT_CONTRACT.md`, the single description of every event this system puts on the broker, mirrored in `tank-node`, `pump-ctl`, `scheduled-valve` and `archimedes-server`.
- `core/processor/domain.Envelope`, the `event`/`device`/`timestamp`/`uptime_s` fields every device event shares, with `At` resolving an event's record time.
- Tests for the processor domain types and for the reading, last-will, wrong-event-type and missing-timestamp paths of both processors.

### Changed

- **Breaking:** a pump event carries `tank_state` (the tank node's verdict: `full`, `not_full`, `unknown`) instead of `distance_cm`. The pump controller now receives the derived state on its own topic and never sees a distance, so it cannot report one. This server is the only subscriber to the tank's level topic.
- **Breaking:** tank and pump events are now the payloads the firmware actually publishes. A tank reading is `{"event":"level","device":…,"valid":…,"distance_cm":…}` rather than `{"tank_id":…,"event_type":…,"distance":…}`, and a pump event is `{"event":"pump","device":…,"state":"on"|"off","reason":…}` rather than `{"pump_id":…,"event_type":"start"|"stop","stop_reason":…}`. The events reaching this server never had the old shape, so nothing that used to be stored stops being stored.
- A tank's `dimensions` must be stored in centimetres, matching the `distance_cm` on the wire.

### Fixed

- Tank readings flagged invalid (`valid:false`, `distance_cm:null`) being stored as a volume computed from a distance of zero, which records a tank filled to the sensor at exactly the moment the level sensor cannot be read.
- Events arriving without a timestamp being stored with a zero time (year 1) instead of the time they were received. The publishing boards have no battery-backed RTC and omit the field until SNTP lands.
- A pump controller's last will (`state:"unknown"`) being treated as an unrecognized event rather than as an explicitly ignored one; it says the controller is unreachable, not that the pump stopped.

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
