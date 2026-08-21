// Package http exposes the read-only HTTP API for tanks and pumps: health
// check, list, get-by-id, and pump run history. It has no write endpoints —
// state changes only arrive through the MQTT stream consumers.
package http

import (
	"archimedes-server/core/log"
	pumpInterfaces "archimedes-server/core/pump/interfaces"
	tankInterfaces "archimedes-server/core/tank/interfaces"
	"net/http"
)

// Serve registers the tank, pump, and health check routes on a new
// http.ServeMux and starts listening on port in a background goroutine.
// It returns immediately; call the returned server's Close or Shutdown to
// stop it. Errors from ListenAndServe are not surfaced to the caller.
func Serve(
	port string,
	readTankRepository tankInterfaces.IReadTank,
	readPumpStatusRepository pumpInterfaces.IReadPumpStatus,
) *http.Server {
	mux := http.NewServeMux()

	tankAPI := &tankAPI{readTankRepository: readTankRepository}
	pumpAPI := &pumpAPI{readPumpStatusRepository: readPumpStatusRepository}

	mux.HandleFunc("GET /health", healthCheckHandler)

	mux.HandleFunc("GET /read/tank", tankAPI.GetTanksHandler)
	mux.HandleFunc("GET /read/tank/{id}", tankAPI.GetTankByIDHandler)

	mux.HandleFunc("GET /read/pump", pumpAPI.GetPumpsHandler)
	mux.HandleFunc("GET /read/pump/{id}", pumpAPI.GetPumpByIDHandler)
	mux.HandleFunc("GET /read/pump/{id}/historic", pumpAPI.GetPumpHistoricHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Log("Server is running on http://localhost:" + port)
	go server.ListenAndServe()

	return server
}
