package ytsearch

import "time"

type SearchResult struct {
	VideoId  string
	Title    string
	Duration time.Duration
}
