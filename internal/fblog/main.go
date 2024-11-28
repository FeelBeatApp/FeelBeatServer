package fblog

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/lmittmann/tint"
)

func ColorizeLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))
}

func messageWithComponent(component component.FeelBeatComponent, message string) string {
	return fmt.Sprintf("[%s] %s", component, message)
}

func Info(component component.FeelBeatComponent, message string, rest ...any) {
	slog.Info(messageWithComponent(component, message), rest...)
}

func Warn(component component.FeelBeatComponent, message string, rest ...any) {
	slog.Warn(messageWithComponent(component, message), rest...)
}

func Error(component component.FeelBeatComponent, message string, rest ...any) {
	slog.Error(messageWithComponent(component, message), rest...)
}

func Debug(component component.FeelBeatComponent, message string, rest ...any) {
	slog.Debug(messageWithComponent(component, message), rest...)
}
