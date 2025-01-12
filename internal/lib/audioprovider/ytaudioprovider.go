package audioprovider

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/thirdparty/youtube"
)

type YTAudioProvider struct {
	urlExtractor *youtube.YTUrlExtractor
}

func NewYTAudioPRovider() YTAudioProvider {
	ex := youtube.NewYTUrlExtractor()
	ex.Init(context.Background())
	return YTAudioProvider{
		urlExtractor: ex,
	}
}

func (y YTAudioProvider) GetUrl(song lib.SongDetails) (string, error) {
	videoId, err := y.pickYoutubeVideo(song)
	if err != nil {
		fblog.Error(component.AudioProvider, "youtube search failed", "err", err)
		return "", err
	}

	url, err := y.urlExtractor.GetDirectUrl(context.Background(), videoId)
	if err != nil {
		fblog.Error(component.AudioProvider, "youtube url extraction failed", "err", err)
		return "", err
	}

	return url, nil
}

func (y YTAudioProvider) pickYoutubeVideo(song lib.SongDetails) (string, error) {
	songs, err := youtube.Search(fmt.Sprintf("%s - %s", song.Title, song.Artist))
	if err != nil {
		return "", err
	}
	if len(songs) == 0 {
		return "", fmt.Errorf("Empty search result")
	}

	deltas := make([]struct {
		delta  time.Duration
		result youtube.SearchResult
	}, 0, len(songs))

	for _, s := range songs {
		deltas = append(deltas, struct {
			delta  time.Duration
			result youtube.SearchResult
		}{
			delta:  (s.Duration - song.Duration).Abs(),
			result: s,
		})
	}

	sort.Slice(deltas, func(i, j int) bool {
		return deltas[i].delta < deltas[j].delta
	})

	return deltas[0].result.VideoId, nil
}
