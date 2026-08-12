package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/pump"
	"archimedes-server/core/tank"
	"net/http"
)

func Serve(port string, readTankRepository tank.IReadTank, readPumpStatusRepository pump.IReadPumpStatus) {
	mux := http.NewServeMux()

	tankAPI := &tankAPI{readTankRepository: readTankRepository}
	pumpAPI := &pumpAPI{readPumpStatusRepository: readPumpStatusRepository}

	mux.HandleFunc("GET /health", healthCheckHandler)

	mux.HandleFunc("GET /read/tank", tankAPI.GetTanksHandler)

	mux.HandleFunc("GET /read/pump", pumpAPI.GetPumpsHandler)
	mux.HandleFunc("GET /read/pump/{id}", pumpAPI.GetPumpByIDHandler)
	mux.HandleFunc("GET /read/pump/{id}/historic", pumpAPI.GetPumpHistoricHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Log("Server is running on http://localhost:" + port)
	server.ListenAndServe()
}
