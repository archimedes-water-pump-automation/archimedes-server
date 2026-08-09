package webserver

import (
	"archimedes-server/core/log"
	"archimedes-server/core/tank"
	"encoding/json"
	"fmt"
	"net/http"
)

func Serve(repository tank.IReadTank) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	mux.HandleFunc("GET /read", func(w http.ResponseWriter, r *http.Request) {
		tankNames, err := repository.GetTanks(r.Context())

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

	mux.HandleFunc("GET /read/{id}", func(w http.ResponseWriter, r *http.Request) {
		tankID := r.PathValue("id")

		tankData, err := repository.GetTankByID(r.Context(), tankID)

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

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Log("Server is running on http://localhost:8080")
	server.ListenAndServe()
}
