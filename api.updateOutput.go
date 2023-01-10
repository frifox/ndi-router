package main

import (
	"net/http"
)

func (a *Api) handleUpdateOutput(w http.ResponseWriter, r *http.Request) {
	log := a.log.Named("handleUpdateOutput")

	query := r.URL.Query()
	outputID := query.Get("output")
	inputID := query.Get("input")

	log.Info("new request", "output", outputID, "input", inputID)

	err := router.UpdateOutput(outputID, inputID)
	if err != nil {
		log.Errorw("router.UpdateOutput", "err", err)
		w.WriteHeader(500)
		return
	}

	log.Infow("done")

	w.WriteHeader(200)

	go a.sendToClients(&RouterMatrix{
		Matrix: router.matrix,
	})
}
