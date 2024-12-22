package youtube

import (
	"context"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

type YTUrlExtractor struct {
	initialized bool
}

func NewYTUrlExtractor() *YTUrlExtractor {
	return &YTUrlExtractor{
		initialized: false,
	}
}

func (e *YTUrlExtractor) Init(ctx context.Context) {
	if e.initialized {
		return
	}

	ytdlp.MustInstall(ctx, nil)
	e.initialized = true
}

func (e *YTUrlExtractor) GetDirectUrl(ctx context.Context, url string) (string, error) {
	dl := ytdlp.New().GetURL()
	output, err := dl.Run(ctx, url)
	if err != nil {
		return "", err
	}
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	splits := strings.Split(output.Stdout, "https://")
	return "https://" + splits[2], nil
}
