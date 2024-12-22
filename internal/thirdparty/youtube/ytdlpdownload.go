package youtube

import (
	"context"
	"fmt"
	"log"

	"github.com/lrstanley/go-ytdlp"
)

func downloadWithYtdlp(url string) {
	ch := make(chan ytdlp.ProgressUpdate)
	defer close(ch)
	dl := ytdlp.New().
		ProgressFunc(1, func(update ytdlp.ProgressUpdate) {
			ch <- update
		}).
		ExtractAudio().
		Output("%(title)s.%(ext)s")

	go progressTracker(ch)
	_, err := dl.Run(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
}

func progressTracker(ch <-chan ytdlp.ProgressUpdate) {
	for update := range ch {
		if !update.Finished.IsZero() {
			fmt.Printf("%s is ready !!!!!\n", update.Filename)
			return
		}

		fmt.Printf("%s: ElapsedTime: %v  [%v%%]\n", update.Filename, update.ETA(), update.Percent())
	}
}
