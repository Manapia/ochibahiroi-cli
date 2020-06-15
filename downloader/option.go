package downloader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DownloadOption struct {
	Parallels int

	ShowProgress bool

	Header map[string]string
}

func ParseHeaderString(source io.Reader) (map[string]string, error) {
	result := make(map[string]string)

	scanner := bufio.NewScanner(source)

	rowCount := 1
	for scanner.Scan() {
		row := scanner.Text()

		var parsed []string
		if row[0] == ':' {
			continue
		}
		parsed = strings.SplitN(row, ":", 2)

		if len(parsed) == 1 {
			return nil, fmt.Errorf("no delimiter on line %d", rowCount)
		}

		result[parsed[0]] = strings.TrimSpace(parsed[1])

		rowCount++
	}

	return result, nil
}
