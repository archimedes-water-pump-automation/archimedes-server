package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/pump/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

// pumpAPI holds the handlers for the /read/pump routes.
type pumpAPI struct {
	readPumpStatusRepository interfaces.IReadPumpStatus
}

// GetPumpsHandler handles GET /read/pump, responding with
// {"pumps": [...]}.
func (api *pumpAPI) GetPumpsHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpsHandler")

	pumpNames, err := api.readPumpStatusRepository.GetPumps(r.Context())
	if err != nil {
		log.Log(fmt.Sprintf("Failed to get pumps: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(map[string]any{
		"pumps": pumpNames,
	})
	if err != nil {
		log.Log(fmt.Sprintf("Failed to marshal pumps: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetPumpByIDHandler handles GET /read/pump/{id}, responding with the
// pump's most recent run. It responds 404 both when the pump ID is unknown
// and when the pump exists but has no recorded runs yet, since the
// repository does not distinguish the two cases.
func (api *pumpAPI) GetPumpByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpByIDHandler")

	pumpID := r.PathValue("id")

	pumpStatusData, err := api.readPumpStatusRepository.GetPumpStatus(r.Context(), pumpID)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to get pump by ID: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if pumpStatusData == nil {
		log.Log(fmt.Sprintf("Pump not found: %q", pumpID))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(pumpStatusData)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to marshal pump: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetPumpHistoricHandler handles GET /read/pump/{id}/historic, responding
// with every recorded run for the pump, most recent first. Unlike
// GetPumpByIDHandler, an unknown pump ID yields an empty array rather than
// 404, since GetPumpStatusHistory never returns nil.
func (api *pumpAPI) GetPumpHistoricHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for pumpAPI.GetPumpHistoricHandler")

	pumpID := r.PathValue("id")

	pumpHistoricData, err := api.readPumpStatusRepository.GetPumpStatusHistory(r.Context(), pumpID)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to get pump historic: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if pumpHistoricData == nil {
		log.Log(fmt.Sprintf("Pump not found: %q", pumpID))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(pumpHistoricData)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to marshal pump: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
