package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsUpgrade = websocket.Upgrader{}

func (a *Api) handleWS(w http.ResponseWriter, r *http.Request) {
	log := a.log.Named("handleWS")

	ws, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Errorw("wsUpgrade.Upgrade", "err", err)
		return
	}
	defer ws.Close()

	a.wsClients[ws] = true
	a.sendToClient(&RouterMatrix{
		Matrix: router.matrix,
	}, ws)

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Errorw("ws.ReadMessage", "err", err)
			break
		}

		var req Emit
		_ = json.Unmarshal(data, &req)
		switch req.Type {
		case "UpdateOutput":
			a.wsUpdateOutput(data)
		default:
			log.Errorw("unhandled ws request", "data", data)
		}
	}

	delete(a.wsClients, ws)
}
