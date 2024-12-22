package audiorepository

import "context"

type audioRequest struct {
	SpotifyId   string
	OnUpdate    func(progress float64)
	OnCompleted func(path string, err error)
}

type audioCache interface {
	// Returns the location where audio file should be stored for song with specified id
	GetOutputPath(spotifyId string) string

	// Returns path to audio file or nil if not exists
	GetAudio(spotifyId string) (string, bool)

	// Use to indicate that song's audio is available at specified path
	RegisterAudio(spotifyId string, path string)
}

type AudioRepository struct {
	downloader audioDownloader
	cache      audioCache

	// Map of downloads in progress
	downloads map[string]audioDownloadTask
}

type audioRepositoryConfig struct {
	Downloader audioDownloader
	Cache      audioCache
}

func NewAudioRepository(config audioRepositoryConfig) AudioRepository {
	return AudioRepository{
		downloader: config.Downloader,
		cache:      config.Cache,
		downloads:  make(map[string]audioDownloadTask),
	}
}

func (ar AudioRepository) RequestAudio(ctx context.Context, req audioRequest) {
	if path, audioExists := ar.cache.GetAudio(req.SpotifyId); !audioExists {
		go req.OnCompleted(path, nil)
		return
	}

	// if audioTask, inProgress := ar.downloads[req.SpotifyId]; !inProgress {
	// 	audioTask.registerListener(&downloadListener{
	// 		ctx:         ctx,
	// 		onUpdate:    req.OnUpdate,
	// 		onCompleted: req.OnCompleted,
	// 	})
	// 	return
	// }
	//
	// outputPath := ar.cache.GetOutputPath(req.SpotifyId)
	// ar.downloads[req.SpotifyId] = newAudioDownloadTask(req.SpotifyId, ar.downloader)
	//
	// ar.downloads[req.SpotifyId].startDownload(outputPath)
}
