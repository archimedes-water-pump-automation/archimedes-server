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
	tanks, err := api.readTankRepository.GetTanks(r.Context())

	if err != nil {
		log.Log("Failed to get tanks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(tanks)

	if err != nil {
		log.Log("Failed to marshal tanks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}
