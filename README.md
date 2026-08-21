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

- A **tank** event carries a raw sensor `distance` reading. The tank processor looks up the tank's registered shape (currently `cylindrical_cone`) and dimensions, converts the distance into a volume, and updates the tank's stored volume.
- A **pump** event is a `start` or `stop` notification. The pump processor opens or closes a run in the pump's status history.
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
export WATER_TANK_TOPIC="archimedes/tank"
export PUMP_STATUS_TOPIC="archimedes/pump"
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
| `WATER_TANK_TOPIC` | Topic the tank stream consumer subscribes to. |
| `PUMP_STATUS_TOPIC` | Topic the pump stream consumer subscribes to. |
| `LOG_FILE` | Path to the file the process appends log lines to. |

## ✨ Features

### MQTT ingestion

Two independent consumers, each on its own goroutine, subscribe at QoS 1 to the tank and pump topics and hand every message off to a use case:

```jsonc
// Water tank event
{ "tank_id": "…", "event_type": "reading", "distance": 1.42, "timestamp": "2026-08-21T10:00:00Z" }

// Pump event
{ "pump_id": "…", "event_type": "start", "timestamp": "2026-08-21T10:00:00Z" }
{ "pump_id": "…", "event_type": "stop", "stop_reason": "tank full", "timestamp": "2026-08-21T10:05:00Z" }
```

### Volume calculation

Each tank has a `tank_shape` and a `dimensions` JSON object stored in PostgreSQL. `core/volume` resolves the shape to an `IVolumeCalculator` implementation; today that's `cylindrical_cone`, which models a vertical cylinder with a partial cone at the bottom:

- `bigger_radius`, `incline_angle`, `cylindrical_height`, `conical_height` — validated against a JSON Schema before use.
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
