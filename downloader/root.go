package downloader

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"sync"
)

func Run(jobs []*Job, option DownloadOption) {
	client := grab.NewClient()

	requests := make([]*grab.Request, 0, len(jobs))

	ch := make(chan struct{}, option.Parallels)

	wg := sync.WaitGroup{}
	wg.Add(len(jobs))

	for _, job := range jobs {
		req, _ := grab.NewRequest(job.SavePath, job.Url)
		requests = append(requests, req)
		go download(ch, &wg, client, req)
	}

	wg.Wait()
}

func download(ch chan struct{}, wg *sync.WaitGroup, client *grab.Client, request *grab.Request) {
	defer wg.Done()
	ch <- struct{}{}

	fmt.Println("start " + request.URL().String())
	response := client.Do(request)

	<-response.Done

	<-ch
}
