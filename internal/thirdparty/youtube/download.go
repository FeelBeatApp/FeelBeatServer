package youtube

type youtubeDownloadUpdate struct {
	url      string
	progress float64
}

type youtubeDownloadResult struct {
	selector songSelector
	url      string
	path     string
}

type youtubeDownload struct {
	update chan youtubeDownloadUpdate
	done   chan youtubeDownloadResult
}
