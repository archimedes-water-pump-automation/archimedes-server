package http

import (
	"archimedes-server/core/log"
	pumpInterfaces "archimedes-server/core/pump/interfaces"
	tankInterfaces "archimedes-server/core/tank/interfaces"
	"net/http"
)

func Serve(port string, readTankRepository tankInterfaces.IReadTank, readPumpStatusRepository pumpInterfaces.IReadPumpStatus) *http.Server {
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
