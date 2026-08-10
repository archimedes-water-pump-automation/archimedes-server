package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/pump"
	"archimedes-server/core/tank"
	"encoding/json"
	"fmt"
	"net/http"
)

func Serve(readTankRepository tank.IReadTank, readPumpStatusRepository pump.IReadPumpStatus) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	mux.HandleFunc("GET /read/tank", func(w http.ResponseWriter, r *http.Request) {
		tankNames, err := readTankRepository.GetTanks(r.Context())

		if err != nil {
			log.Log("Failed to get tanks: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		response, err := json.Marshal(tankNames)

		if err != nil {
			log.Log("Failed to marshal tanks: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	})

	mux.HandleFunc("GET /read/tank/{id}", func(w http.ResponseWriter, r *http.Request) {
		tankID := r.PathValue("id")

		tankData, err := readTankRepository.GetTankByID(r.Context(), tankID)

		if err != nil {
			log.Log("Failed to get tanks: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if tankData == nil {
			log.Log("Tank not found: " + tankID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response, err := json.Marshal(tankData)

		if err != nil {
			log.Log("Failed to marshal tank: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	})

	mux.HandleFunc("GET /read/pump", func(w http.ResponseWriter, r *http.Request) {
		pumpNames, err := readPumpStatusRepository.GetPumps(r.Context())

		if err != nil {
			log.Log("Failed to get pumps: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		response, err := json.Marshal(pumpNames)

		if err != nil {
			log.Log("Failed to marshal pumps: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	})

	mux.HandleFunc("GET /read/pump/{id}", func(w http.ResponseWriter, r *http.Request) {
		pumpID := r.PathValue("id")

		pumpStatusData, err := readPumpStatusRepository.GetPumpStatus(r.Context(), pumpID)

		if err != nil {
			log.Log("Failed to get pumps: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if pumpStatusData == nil {
			log.Log("Pump not found: " + pumpID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response, err := json.Marshal(pumpStatusData)

		if err != nil {
			log.Log("Failed to marshal pump: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Log("Server is running on http://localhost:8080")
	server.ListenAndServe()
}
