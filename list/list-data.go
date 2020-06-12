package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const DownloadListVersion = 1

type DownloadListData struct {
	Version int                       `json:"version"`
	List    map[int]*DownloadListItem `json:"list"`
}

type DownloadListFile struct {
	DownloadListData
	FilePath string
}

type DownloadListItem struct {
	ID     int    `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

func (dld *DownloadListData) ToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(dld)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode JSON data\n%v", err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, jsonBytes, "", "  "); err != nil {
		return nil, fmt.Errorf("Failed to format JSON data style\n%v", err)
	}

	return out.Bytes(), nil
}

func (dld *DownloadListData) Clear() {
	dld.List = make(map[int]*DownloadListItem)
}

func (dlf *DownloadListFile) FileExists() bool {
	_, err := os.Stat(dlf.FilePath)
	return err == nil
}

func (dlf *DownloadListFile) LoadFile() error {
	data, err := ioutil.ReadFile(dlf.FilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &dlf); err != nil {
		return err
	}

	idList := make(map[int]struct{})
	for _, item := range dlf.List {
		_, exists := idList[item.ID]
		if exists {
			return fmt.Errorf("duplicate id:%d in the download list", item.ID)
		}
	}

	return nil
}
