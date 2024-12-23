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
	downloader.m.AssertExpectations(t)
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
	downloader.m.AssertExpectations(t)
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
	downloader.m.AssertExpectations(t)
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

	downloader.m.AssertExpectations(t)
}

func TestMultiListenersUpdates(t *testing.T) {
	assert := assert.New(t)
	targetPath := "<path>"
	downloader := newDownloaderStub()
	updates1 := make([]float64, 0)
	updates2 := make([]float64, 0)
	updates3 := make([]float64, 0)

	task := newAudioDownloadTask("multi update", downloader)
	ctx1, cancel1 := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	goWrite := make(chan bool)
	goRead := make(chan bool)

	wg.Add(2)
	err := task.registerListener(&downloadListener{
		ctx: ctx1,
		onUpdate: func(progress float64) {
			updates1 = append(updates1, progress)
		},
		onCompleted: func(path string, err error) {
			wg.Done()
		},
	})
	assert.Nil(err)
	err = task.registerListener(&downloadListener{
		ctx: context.Background(),
		onUpdate: func(progress float64) {
			updates2 = append(updates2, progress)
		},
		onCompleted: func(path string, err error) {
			wg.Done()
		},
	})
	assert.Nil(err)
	err = task.registerListener(&downloadListener{
		ctx: context.Background(),
		onUpdate: func(progress float64) {
			updates3 = append(updates3, progress)
		},
		onCompleted: func(path string, err error) {
			wg.Done()
		},
	})
	assert.Nil(err)

	downloader.m.On("Download", mock.Anything, mock.Anything).Once().After(time.Millisecond * 5).
		Run(func(args mock.Arguments) {
			req := args.Get(1).(downloadRequest)
			assert.NotNil(req.onUpdate)

			<-goWrite
			req.onUpdate(10)
			goRead <- true
			<-goWrite
			req.onUpdate(25.24)
			goRead <- true
			<-goWrite
			req.onUpdate(80)
			goRead <- true
		}).Return(nil)

	task.startDownload(targetPath)

	assert.Empty(updates1)
	assert.Empty(updates2)
	assert.Empty(updates3)

	goWrite <- true
	<-goRead

	assert.Equal([]float64{10}, updates1)
	assert.Equal([]float64{10}, updates2)
	assert.Equal([]float64{10}, updates3)

	cancel1()
	time.Sleep(time.Millisecond * 5) // Wait for cancel being processed
	goWrite <- true
	<-goRead

	assert.Equal([]float64{10}, updates1)
	assert.Equal([]float64{10, 25.24}, updates2)
	assert.Equal([]float64{10, 25.24}, updates3)

	goWrite <- true
	<-goRead

	assert.Equal([]float64{10}, updates1)
	assert.Equal([]float64{10, 25.24, 80}, updates2)
	assert.Equal([]float64{10, 25.24, 80}, updates3)

	assert.NoError(testutils.WaitTimeout(&wg, commonTimeout))

	downloader.m.AssertExpectations(t)
}
