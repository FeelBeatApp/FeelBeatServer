package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/ws"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/knadh/koanf/v2"
)

const (
	envPrefix = "FEELBEAT_"
	tomlPath  = "config.toml"
)

func main() {
	fblog.ColorizeLogger()

	config := koanf.New(".")
	err := setupConfig(config)
	if err != nil {
		fblog.Error(component.FeelBeatServer, "Reading config", "error", err)
		os.Exit(1)
	}

	fblog.Info(component.FeelBeatServer, "config loaded", "config", config.All())

	port := config.MustInt("websocket.port")
	path := config.MustString("websocket.path")

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWebsockets(hub, w, r)
	})

	setupAPI(auth.AuthorizeThroughSpotify)

	fblog.Info(component.FeelBeatServer, "Server started", "port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
