package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"com.github.feelbeatapp.feelbeatserver/internal/networking"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/lmittmann/tint"
)

const (
	ENV_PREFIX = "FEELBEAT_"
	TOML_PATH  = "config.toml"
)

func colorizeLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))
}

func setupConfig(config *koanf.Koanf) error {
	err := config.Load(file.Provider(TOML_PATH), toml.Parser())

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		slog.Info("config file not detected")
	} else {
		slog.Info("config file loaded")
	}

	envProvider := env.Provider(ENV_PREFIX, ".", func(key string) string {
		envvar := strings.ToLower(strings.TrimPrefix(key, ENV_PREFIX))

		return strings.ReplaceAll(envvar, "_", ".")
	})

	return config.Load(envProvider, nil)
}

func main() {
	colorizeLogger()

	config := koanf.New(".")
	err := setupConfig(config)
	if err != nil {
		slog.Error("Reading config", "error", err)
		os.Exit(1)
	}

	slog.Info("config", "config", config.All())

	port := config.MustInt("websocket.port")
	path := config.MustString("websocket.path")

	hub := networking.NewHub()
	go hub.Run()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		networking.ServeWebsockets(hub, w, r)
	})

	slog.Info("Server started", "port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
