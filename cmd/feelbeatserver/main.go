package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/ws"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/roomrepository"
	"github.com/feelbeatapp/feelbeatserver/internal/thirdparty/spotify"
	"github.com/knadh/koanf/v2"
)

const (
	envPrefix = "FEELBEAT_"
	tomlPath  = "config.toml"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	fblog.ColorizeLogger()

	config := koanf.New(".")
	err := setupConfig(config)
	if err != nil {
		fblog.Error(component.FeelBeatServer, "Reading config", "error", err)
		os.Exit(1)
	}

	fblog.Info(component.FeelBeatServer, "config loaded", "config", config.All())

	port := config.MustInt("server.port")
	path := config.MustString("websocket.path")

	roomRepo := roomrepository.NewInMemoryRoomRepository(spotify.SpotifyApi{}, ws.NewWSHub)

	ws := ws.NewWSHandler(path, roomRepo)
	ws.ServeWebsockets(path, auth.AuthorizeThroughSpotify)

	setupAPI(auth.AuthorizeThroughSpotify, roomRepo)

	fblog.Info(component.FeelBeatServer, "Server started", "port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
