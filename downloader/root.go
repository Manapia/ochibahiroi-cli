package downloader

import (
	"github.com/cavaliercoder/grab"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"path/filepath"
	"sync"
	"time"
)

func Run(jobs []*Job, option DownloadOption) {
	client := grab.NewClient()

	requests := make([]*grab.Request, 0, len(jobs))

	progressBar := mpb.New(mpb.WithWidth(64))

	ch := make(chan struct{}, option.Parallels)

	wg := sync.WaitGroup{}
	wg.Add(len(jobs))

	for _, job := range jobs {
		req, _ := grab.NewRequest(job.SavePath, job.Url)
		requests = append(requests, req)
		go download(ch, &wg, progressBar, client, req, option)
	}

	wg.Wait()
}

func download(ch chan struct{}, wg *sync.WaitGroup, p *mpb.Progress, client *grab.Client, request *grab.Request, option DownloadOption) {
	defer wg.Done()
	ch <- struct{}{}

	response := client.Do(request)

	before := response.BytesComplete()
	var bar *mpb.Bar
	if option.ShowProgress {
		bar = p.AddBar(
			response.Size,
			mpb.PrependDecorators(
				decor.Name(filepath.Base(request.Filename)),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			))
	}

	for {
		if bar != nil {
			bar.IncrBy(int(response.BytesComplete() - before))
			before = response.BytesComplete()
		}

		if response.IsComplete() {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	<-ch
}
