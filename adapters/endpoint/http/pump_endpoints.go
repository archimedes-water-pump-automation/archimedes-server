package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/pump/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

type pumpAPI struct {
	readPumpStatusRepository interfaces.IReadPumpStatus
}

func (api *pumpAPI) GetPumpsHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpsHandler")

	pumpNames, err := api.readPumpStatusRepository.GetPumps(r.Context())

	if err != nil {
		log.Log("Failed to get pumps: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(map[string]any{
		"pumps": pumpNames,
	})

	if err != nil {
		log.Log("Failed to marshal pumps: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func (api *pumpAPI) GetPumpByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpByIDHandler")

	pumpID := r.PathValue("id")

	pumpStatusData, err := api.readPumpStatusRepository.GetPumpStatus(r.Context(), pumpID)

	if err != nil {
		log.Log("Failed to get pump by ID: " + err.Error())
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
}

func (api *pumpAPI) GetPumpHistoricHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpHistoricHandler")

	pumpID := r.PathValue("id")

	pumpHistoricData, err := api.readPumpStatusRepository.GetPumpStatusHistory(r.Context(), pumpID)

	if err != nil {
		log.Log("Failed to get pump historic: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if pumpHistoricData == nil {
		log.Log("Pump not found: " + pumpID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(pumpHistoricData)

	if err != nil {
		log.Log("Failed to marshal pump: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}
