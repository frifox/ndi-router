package main

import (
	"encoding/json"
)

func (a *Api) wsUpdateOutput(data []byte) {
	log := a.log.Named("wsUpdateOutput")

	var req UpdateOutput
	_ = json.Unmarshal(data, &req)
	log.Infow("new emit", "req", req)

	err := router.UpdateOutput(req.Output, req.Input)
	if err != nil {
		log.Errorw("router.UpdateOutput", "err", err)
		return
	}
	log.Info("updated", "output", req.Output, "input", req.Input)

	go a.sendToClients(&RouterMatrix{
		Matrix: router.matrix,
	})
}
