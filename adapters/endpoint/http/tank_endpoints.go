package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/tank/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

// tankAPI holds the handlers for the /read/tank routes.
type tankAPI struct {
	readTankRepository interfaces.IReadTank
}

// GetTanksHandler handles GET /read/tank, responding with
// {"tanks": [...]}.
func (api *tankAPI) GetTanksHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for tankAPI.GetTanksHandler")

	tanks, err := api.readTankRepository.GetTanks(r.Context())
	if err != nil {
		log.Log(fmt.Sprintf("Failed to get tanks: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(map[string]any{
		"tanks": tanks,
	})
	if err != nil {
		log.Log(fmt.Sprintf("Failed to marshal tanks: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetTankByIDHandler handles GET /read/tank/{id}, responding with the
// tank's latest status, or 404 if the tank does not exist.
func (api *tankAPI) GetTankByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("Received request for tankAPI.GetTankByIDHandler")

	tankID := r.PathValue("id")

	tankStatusData, err := api.readTankRepository.GetTankByID(r.Context(), tankID)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to get tank by ID: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tankStatusData == nil {
		log.Log(fmt.Sprintf("Tank not found: %q", tankID))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(tankStatusData)
	if err != nil {
		log.Log(fmt.Sprintf("Failed to marshal tank: %q", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
