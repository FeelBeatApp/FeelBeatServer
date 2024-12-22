package audiorepository

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/lib/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const commonTimeout = time.Millisecond * 100

type DownloaderStub struct {
	m *mock.Mock
}

func newDownloaderStub() DownloaderStub {
	return DownloaderStub{
		m: new(mock.Mock),
	}
}

func (d DownloaderStub) Download(ctx context.Context, option downloadRequest) error {
	args := d.m.Called(ctx, option)

	return args.Error(0)
}

func TestSingleListenerSuccess(t *testing.T) {
	assert := assert.New(t)
	downloader := newDownloaderStub()
	var wg sync.WaitGroup
	task := newAudioDownloadTask("single success", downloader)
	targetPath := "<path>"

	downloader.m.On("Download", mock.Anything, mock.Anything).Once().Return(nil)

	wg.Add(1)
	err := task.registerListener(&downloadListener{
		ctx: context.Background(),
		onCompleted: func(path string, err error) {
			assert.Nil(err)
			assert.Equal(targetPath, path)
			wg.Done()
		},
	})

	assert.Nil(err)

	task.startDownload(targetPath)

	assert.NoError(testutils.WaitTimeout(&wg, commonTimeout))
}

func TestSingleListenerExternalError(t *testing.T) {
	assert := assert.New(t)
	downloader := newDownloaderStub()
	var wg sync.WaitGroup
	task := newAudioDownloadTask("single fail", downloader)
	targetPath := "<path>"
	externalErr := errors.New("Some http error")

	downloader.m.On("Download", mock.Anything, mock.Anything).Once().Return(externalErr)

	wg.Add(1)
	err := task.registerListener(&downloadListener{
		ctx: context.Background(),
		onCompleted: func(path string, err error) {
			assert.Empty(path)
			assert.Equal(externalErr, err)
			wg.Done()
		},
	})

	assert.Nil(err)

	task.startDownload(targetPath)

	assert.NoError(testutils.WaitTimeout(&wg, commonTimeout))
}

func TestSingleListenerCancel(t *testing.T) {
	assert := assert.New(t)
	downloader := newDownloaderStub()
	var wg sync.WaitGroup
	task := newAudioDownloadTask("single success", downloader)
	targetPath := "<path>"
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	downloader.m.On("Download", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
		wg.Done()
	}).Once().Return(context.Canceled)

	err := task.registerListener(&downloadListener{
		ctx: ctx,
	})
	assert.Nil(err)

	task.startDownload(targetPath)

	cancel()

	assert.NoError(testutils.WaitTimeout(&wg, commonTimeout))
}

func TestMultiListenersPartialCancel(t *testing.T) {
	assert := assert.New(t)
	targetPath := "<path>"
	downloader := newDownloaderStub()
	successCount := 5
	cancelCount := 5
	successMap := make(map[int]bool)
	cancels := make(map[int]context.CancelFunc)
	var wg sync.WaitGroup

	downloader.m.On("Download", mock.Anything, mock.Anything).Once().After(time.Millisecond * 5).Return(nil)

	task := newAudioDownloadTask("multi cancel", downloader)

	task.startDownload(targetPath)

	wg.Add(successCount)
	for i := 0; i < cancelCount+successCount; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancels[i] = cancel
		err := task.registerListener(&downloadListener{
			ctx: ctx,
			onCompleted: func(path string, err error) {
				successMap[i] = true
				wg.Done()
			},
		})

		assert.Nil(err)
	}

	for i := 3; i < 3+cancelCount; i++ {
		cancels[i]()
	}

	assert.NoError(testutils.WaitTimeout(&wg, commonTimeout))

	for i := 0; i < 3; i++ {
		assert.True(successMap[i], "%v", i)
	}
	for i := 3; i < 3+cancelCount; i++ {
		_, ok := successMap[i]
		assert.False(ok, "%v", i)
	}
	for i := 3 + cancelCount; i < cancelCount+successCount; i++ {
		assert.True(successMap[i], "%v", i)
	}
}
