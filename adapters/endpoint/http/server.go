// Package http exposes the read-only HTTP API for tanks and pumps: health
// check, list, get-by-id, and pump run history. It has no write endpoints —
// state changes only arrive through the MQTT stream consumers.
package http

import (
	"archimedes-server/core/log"
	pumpInterfaces "archimedes-server/core/pump/interfaces"
	tankInterfaces "archimedes-server/core/tank/interfaces"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Serve registers the tank, pump, and health check routes on a new
// http.ServeMux, opens port, and serves on it in a background goroutine.
// It returns immediately once the port is open; call the returned server's
// Close or Shutdown to stop it.
//
// The port is opened before returning rather than inside the goroutine, so
// a port that cannot be opened is the caller's error to act on. Serving in
// the background and discarding that error would leave the process up with
// no read API and nothing said about it — the tank and pump consumers
// would go on recording events that nothing could then be read back from.
func Serve(
	port string,
	readTankRepository tankInterfaces.IReadTank,
	readPumpStatusRepository pumpInterfaces.IReadPumpStatus,
) (*http.Server, error) {
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

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to listen on port %s: %q", port, err.Error()))
		return nil, fmt.Errorf("listening on port %s: %w", port, err)
	}

	log.Log("Server is running on http://localhost:" + port)

	go func() {
		// A closed server is how Serve ends on shutdown, not a failure.
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Log(fmt.Sprintf("HTTP server stopped serving: %q", err.Error()))
		}
	}()

	return server, nil
}
