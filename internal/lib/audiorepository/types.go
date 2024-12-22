package audiorepository

import "context"

type downloadRequest struct {
	spotifyId string
	path      string
	onUpdate  func(progress float64)
}

type audioDownloader interface {
	Download(ctx context.Context, options downloadRequest) error
}
