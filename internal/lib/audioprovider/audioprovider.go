package audioprovider

import "github.com/feelbeatapp/feelbeatserver/internal/lib"

type AudioProvider interface {
	GetUrl(lib.SongDetails) (string, error)
}
