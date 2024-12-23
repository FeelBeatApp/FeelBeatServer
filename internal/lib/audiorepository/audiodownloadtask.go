package audiorepository

import (
	"context"
	"errors"
	"sync"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
)

type downloadState int

const (
	idle = iota
	inProgress
	finished
)

type downloadListener struct {
	ctx         context.Context
	onUpdate    func(progress float64)
	onCompleted func(path string, err error)
}

type audioDownloadTask struct {
	mu sync.RWMutex

	spotifyId  string
	state      downloadState
	downloader audioDownloader
	outputPath string

	listeners      map[*downloadListener]chan bool
	downloadCtx    context.Context
	cancelDownload context.CancelFunc
}

func newAudioDownloadTask(spotifyId string, downloader audioDownloader) *audioDownloadTask {
	return &audioDownloadTask{
		spotifyId:  spotifyId,
		state:      idle,
		downloader: downloader,
		listeners:  make(map[*downloadListener]chan bool),
	}
}

func (t *audioDownloadTask) startDownload(outputPath string) {
	t.mu.Lock()
	if t.state != idle {
		t.mu.Unlock()
		return
	}

	t.outputPath = outputPath
	t.state = inProgress
	ctx, cncl := context.WithCancel(context.Background())
	t.downloadCtx = ctx
	t.cancelDownload = cncl
	t.mu.Unlock()

	go t.downloadProcess()
}

func (t *audioDownloadTask) downloadProcess() {
	fblog.Info(component.AudioDownloadTask, "download started", "id", t.spotifyId)
	err := t.downloader.Download(t.downloadCtx, downloadRequest{
		spotifyId: t.spotifyId,
		path:      t.outputPath,
		onUpdate:  t.passProgressToListeners,
	})

	t.mu.Lock()
	t.state = finished
	t.mu.Unlock()

	t.handleFinish(err)
	fblog.Info(component.AudioDownloadTask, "download finished", "id", t.spotifyId)
}

func (t *audioDownloadTask) passProgressToListeners(progress float64) {
	fblog.Info(component.AudioDownloadTask, "passing update to listeners", "id", t.spotifyId, "progress", progress)
	for _, l := range t.pullListenersSafely() {
		if l.onUpdate != nil {
			l.onUpdate(progress)
		}
	}
}

func (t *audioDownloadTask) handleFinish(err error) {
	for _, l := range t.pullListenersSafely() {
		close(t.listeners[l])
		if l.onCompleted != nil {
			path := t.outputPath
			if err != nil {
				path = ""
			}
			l.onCompleted(path, err)
		}
	}
}

func (t *audioDownloadTask) handleCancellation(listener *downloadListener, end <-chan bool) {
	select {
	case <-end:
	case <-listener.ctx.Done():
		t.removeListener(listener)
		return
	}
}

func (t *audioDownloadTask) pullListenersSafely() []*downloadListener {
	// Copy to limit critical section
	t.mu.RLock()
	listeners := make([]*downloadListener, 0, len(t.listeners))
	for l := range t.listeners {
		listeners = append(listeners, l)
	}
	t.mu.RUnlock()

	return listeners
}

func (t *audioDownloadTask) registerListener(listener *downloadListener) error {
	fblog.Info(component.AudioDownloadTask, "registering listener", "id", t.spotifyId, "listener", listener)
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state == finished {
		return errors.New("Tried to register listeners on finished task")
	}

	t.listeners[listener] = make(chan bool)
	go t.handleCancellation(listener, t.listeners[listener])

	return nil
}

func (t *audioDownloadTask) removeListener(listener *downloadListener) {
	fblog.Info(component.AudioDownloadTask, "removing listener", "id", t.spotifyId, "listener", listener)
	t.mu.Lock()
	defer t.mu.Unlock()
	close(t.listeners[listener])
	delete(t.listeners, listener)

	if len(t.listeners) == 0 && t.state == inProgress {
		t.cancelDownload()
	}
}
