package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/networking"
	"github.com/knadh/koanf/v2"
)

const (
	ENV_PREFIX = "FEELBEAT_"
	TOML_PATH  = "config.toml"
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

	hub := networking.NewHub()
	go hub.Run()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		networking.ServeWebsockets(hub, w, r)
	})

	fblog.Info(component.FeelBeatServer, "Server started", "port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
