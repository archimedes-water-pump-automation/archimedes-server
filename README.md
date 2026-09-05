# archimedes-server

[![Go Version](https://img.shields.io/github/go-mod/go-version/archimedes-water-pump-automation/archimedes-server)](https://go.dev/) [![License](https://img.shields.io/github/license/archimedes-water-pump-automation/archimedes-server)](./LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/archimedes-water-pump-automation/archimedes-server)](https://goreportcard.com/report/github.com/archimedes-water-pump-automation/archimedes-server)

Archimedes is a backend service for monitoring water tanks and the pumps that fill them. It consumes sensor and pump-status events from an MQTT broker, computes each tank's current volume from a distance reading and the tank's geometry, persists the results to PostgreSQL, and exposes them through a read-only HTTP API.

## How it works

```
                 ┌─────────────────┐        ┌──────────────────────┐
  MQTT broker ──▶│  stream consumer │──────▶ │ tank / pump processor │──▶ PostgreSQL
 (tank + pump     │ (adapters/stream) │       │ (core/processor)      │
  topics)         └─────────────────┘        └──────────────────────┘

  PostgreSQL ──▶ read repositories ──▶ HTTP API (adapters/endpoint/http)
```

- A **tank** event carries a sensor `distance_cm` reading published by [`tank-node`](https://github.com/archimedes-water-pump-automation/tank-node), to this server and to no one else. The tank processor looks up the tank's registered shape (currently `cylindrical_cone`) and dimensions, converts the distance into a volume, and updates the tank's stored volume. A reading the node flagged invalid stores nothing — `distance_cm` is then `null`, and treating that as zero would record a tank filled to the sensor.
- A **pump** event is a `state` transition published by [`pump-ctl`](https://github.com/archimedes-water-pump-automation/pump-ctl). The pump processor opens a run on `"on"` and closes it on `"off"`, storing the event's `reason` as the stop reason.
- Both payloads are fixed by [MQTT_CONTRACT.md](MQTT_CONTRACT.md), which is mirrored in all four repositories of this system. Changing a field here means changing it in the firmware that publishes it.
- The HTTP API only reads what the processors have written — there are no write endpoints.

## 🚀 Getting Started

Requires Go 1.25+, a PostgreSQL database (with the schema described below), and an MQTT broker.

```bash
git clone https://github.com/archimedes-water-pump-automation/archimedes-server.git
cd archimedes-server
go build -o archimedes-server .
```

Configure the process through environment variables and run it:

```bash
export DB_CONN_STRING="postgres://user:pass@localhost:5432/archimedes"
export DB_TLS_ENABLED=false
export MQTT_BROKER_URL="tcp://localhost:1883"
export WATER_TANK_TOPIC="watertank/tank-01/level"
export PUMP_STATUS_TOPIC="watertank/pump-01/pump"
export LOG_FILE="./archimedes.log"

./archimedes-server
```

The server starts an HTTP server on port `8080` and blocks until it receives `SIGINT` or `SIGTERM`, at which point it drains both MQTT consumers, closes the database pool, and exits.

### Configuration reference

| Variable | Description |
| --- | --- |
| `DB_CONN_STRING` | PostgreSQL connection string (pgx format). |
| `DB_TLS_ENABLED` | `true` to connect over TLS with certificate verification skipped (for providers with managed certs the client can't validate). |
| `MQTT_BROKER_URL` | Broker URL, e.g. `tcp://host:1883`. |
| `WATER_TANK_TOPIC` | Topic the tank stream consumer subscribes to, e.g. `watertank/tank-01/level`. |
| `PUMP_STATUS_TOPIC` | Topic the pump stream consumer subscribes to, e.g. `watertank/pump-01/pump`. |
| `LOG_FILE` | Path to the file the process appends log lines to. |

## ✨ Features

### MQTT ingestion

Two independent consumers, each on its own goroutine, subscribe at QoS 1 to the tank and pump topics and hand every message off to a use case. Every message on every topic shares an envelope — `event`, `device`, and an optional `timestamp` and `uptime_s` — followed by the fields of the event itself:

```jsonc
// Water tank level reading, from tank-node
{ "event": "level", "device": "tank-01", "timestamp": "2026-09-05T03:10:12Z",
  "valid": true, "distance_cm": 62.5, "uptime_s": 360 }

// Sensor unreadable, or the node's last will: stored as nothing at all
{ "event": "level", "device": "tank-01", "valid": false,
  "distance_cm": null, "reason": "sensor_unreadable" }

// Pump transitions, from pump-ctl
{ "event": "pump", "device": "pump-01", "timestamp": "2026-09-05T03:10:12Z",
  "state": "on", "reason": "flow_confirmed", "flow_lpm": 11.40,
  "tank_state": "not_full", "uptime_s": 338 }
{ "event": "pump", "device": "pump-01", "timestamp": "2026-09-05T03:14:41Z",
  "state": "off", "reason": "tank_full", "flow_lpm": 0.0,
  "tank_state": "full", "uptime_s": 607 }
```

`device` is the key each event is stored against: it must match the tank's or pump's `id` in the database.

`tank_state` (`full`, `not_full` or `unknown`) is what the tank node told the controller, not what the controller measured: the distance stays between `tank-node` and this server, while `pump-ctl` receives only the derived state on a separate topic, where it can stop a pump but never start one. This server is the only subscriber to the level topic.

`timestamp` is UTC and optional, because the publishing boards have no battery-backed RTC and a last will is published by the broker rather than by the device. When it is absent the server records the moment it received the message, which is the closest true answer available.

Retained messages are ignored, on both topics. Retention exists for dashboards: this server connects with a clean session, so the broker replays a retained pump event on every reconnect, and processing it again would record a second run for a start that happened once. The cost is that a server starting mid-run does not learn the pump is already running until its next transition — the controller also announces its real state with `reason: "boot"` on its first connection after a restart, which closes a run left open by a crash.

A pump event with `"state": "unknown"` is the controller's last will. It is logged and stored as nothing: an unreachable controller is not a stopped pump, and closing a run on it would put a fabricated stop time in the history.

### Volume calculation

Each tank has a `tank_shape` and a `dimensions` JSON object stored in PostgreSQL. `core/volume` resolves the shape to an `IVolumeCalculator` implementation; today that's `cylindrical_cone`, which models a vertical cylinder with a partial cone at the bottom:

- `bigger_radius`, `incline_angle`, `cylindrical_height`, `conical_height` — validated against a JSON Schema before use.
- **Stored in centimetres**, matching the `distance_cm` on the wire: the calculator works in whatever unit it is given, so a tank measured in metres and a sensor reporting centimetres would produce a volume that is wrong by six orders of magnitude.
- The reported volume is the sum of the fluid held in the cylindrical section and whatever portion of the cone is submerged.

Adding a new shape means implementing `core/volume/interfaces.IVolumeCalculator` and registering it in `adapters/database/postgresql.getVolumeType.GetVolumeFromShape`.

### Read API

| Method & Path | Description |
| --- | --- |
| `GET /health` | Liveness check. |
| `GET /read/tank` | List all tanks. |
| `GET /read/tank/{id}` | Latest status (capacity, volume, timestamps) for one tank. |
| `GET /read/pump` | List all pumps. |
| `GET /read/pump/{id}` | Most recent run for one pump. |
| `GET /read/pump/{id}/historic` | Full run history for one pump, most recent first. |

```bash
curl http://localhost:8080/read/tank
curl http://localhost:8080/read/pump/{id}/historic
```

## 🤝 Contributing

Please read the [contributing guide](CONTRIBUTING.md) before submitting a PR.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
