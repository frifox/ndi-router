package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type Api struct {
	mux *http.ServeMux
	log *zap.SugaredLogger

	wsClients map[*websocket.Conn]bool
}

func (a *Api) Init(logger *zap.SugaredLogger) {
	a.log = logger.Named("Api")
	a.wsClients = make(map[*websocket.Conn]bool)

	a.mux = http.NewServeMux()
	a.mux.HandleFunc("/", a.handleWS)
	a.mux.HandleFunc("/updateOutput", a.handleUpdateOutput)
}

func (a *Api) sendToClients(emit interface{}) {
	log := a.log.Named("sendToClients")
	log.Infow("sending to all clients", "emit", fmt.Sprintf("%+v", emit))

	for ws, _ := range a.wsClients {
		a.sendToClient(emit, ws)
	}
}
func (a *Api) sendToClient(emit interface{}, client *websocket.Conn) {
	log := a.log.Named("sendToClient").With("RemoteAddr", client.RemoteAddr())

	data, _ := json.Marshal(emit)

	// TODO undo hack
	if emit, ok := emit.(*RouterMatrix); ok {
		data, _ = json.Marshal(emit.Matrix)
	}

	log.Infow("sending", "data", string(data))

	err := client.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Errorw("ws.WriteMessage", "err", err)
		return
	}
}
