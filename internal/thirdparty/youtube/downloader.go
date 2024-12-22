package youtube

import (
	"context"

	"github.com/lrstanley/go-ytdlp"
)


type YTDownloder struct {
	initialized bool
}

func NewYTDownloader() *YTDownloder {
	return &YTDownloder{
		initialized: false,
	}
}

type AudioDownload struct {
}

func (yt *YTDownloder) Init(ctx context.Context) {
	if yt.initialized {
		return
	}

	ytdlp.MustInstall(ctx, nil)
	yt.initialized = true
}

func (yt *YTDownloder) Download(ctx context.Context, sng songSelector) {
	
}
