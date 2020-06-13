package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Job struct {
	Url string

	SavePath string
}

type JobListBuilder struct {
	url string

	savePath string

	startNumber int

	endNumber int

	step int

	useIncrementalCount bool
}

func (j *Job) ToDisplayString() string {
	return fmt.Sprintf("\"%s\" => \"%s\"", j.Url, j.SavePath)
}

func (b *JobListBuilder) SetUrl(url string) {
	b.url = url
}

func (b *JobListBuilder) SetSavePath(path string) {
	strings.TrimRight(path, string(os.PathSeparator))
	b.savePath = path
}

func (b *JobListBuilder) SetStart(number int) {
	b.startNumber = number
}

func (b *JobListBuilder) SetEnd(number int) {
	b.endNumber = number
}

func (b *JobListBuilder) SetStep(step int) {
	b.step = step
}

func (b *JobListBuilder) SetUserIncrementalCount(value bool) {
	b.useIncrementalCount = value
}

func (b *JobListBuilder) Build() ([]*Job, error) {
	jobs := make([]*Job, 0, 10)

	if b.savePath == "" {
		current, err := filepath.Abs(".")
		if err != nil {
			b.savePath = current
		}
	}

	count := 1
	if b.url != "" {
		for i := b.startNumber; i <= b.endNumber; i += b.step {
			url := fmt.Sprintf(b.url, i)

			savePath := filepath.Join(b.savePath, filepath.Base(url))
			if b.useIncrementalCount {
				savePath = filepath.Join(b.savePath, strconv.Itoa(count)+filepath.Ext(url))
			}

			jobs = append(jobs, &Job{
				Url:      url,
				SavePath: savePath,
			})

			count++
		}
	}

	return jobs, nil
}
