package http

import (
	"archimedes-server/core/log"
	"archimedes-server/core/tank"
	"encoding/json"
	"fmt"
	"net/http"
)

type tankAPI struct {
	readTankRepository tank.IReadTank
}

func (api *tankAPI) GetTanksHandler(w http.ResponseWriter, r *http.Request) {
	tankNames, err := api.readTankRepository.GetTanks(r.Context())

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
}

func (api *tankAPI) GetTankByIDHandler(w http.ResponseWriter, r *http.Request) {
	tankID := r.PathValue("id")

	tankData, err := api.readTankRepository.GetTankByID(r.Context(), tankID)

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
}
