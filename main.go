package main

import (
	"net/http"
	"os"
)

var log *Logger
var conf *Config
var db *SQLite
var router *Router
var api *Api

func init() {
	log = &Logger{}
	log.Init()
	log.SwitchLevel("debug")

	conf = &Config{}
	if err := conf.Load("ndi-router.conf"); err != nil {
		log.Fatalw("conf.Load", "err", err)
		os.Exit(1)
	}

	db = &SQLite{}
	if err := db.Init(log.SugaredLogger); err != nil {
		log.Fatalw("db.Init", "err", err)
		os.Exit(1)
	}

	router = &Router{}
	if err := router.Init(log.SugaredLogger); err != nil {
		log.Fatalw("router.Init", "err", err)
		os.Exit(1)
	}

	api = &Api{}
	api.Init(log.SugaredLogger)
}

func main() {
	sources := router.GetSources()
	for _, source := range sources {
		log.Infow("ndi source found", "name", source.Name, "url", source.URL)
	}

	for id, source := range conf.Inputs {
		router.InitInput(id, source)
	}
	for id, chanName := range conf.Outputs {
		router.InitOutput(id, chanName)
	}

	err := http.ListenAndServe(conf.ApiAddr, api.mux)
	if err != nil {
		log.Fatalw("http.ListenAndServe", "err", err)
	}
}
