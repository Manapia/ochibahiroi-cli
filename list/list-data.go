package list

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (dld *DownloadListData) Validation() (errorList []error) {
	errorList = make([]error, 0)

	idMap := make(map[int]bool, len(dld.List))
	outputMap := make(map[string]int, len(dld.List))

	for _, item := range dld.List {
		_, idExists := idMap[item.ID]
		if idExists {
			errorList = append(errorList, errors.New("duplicate id: 4"))
		} else {
			idMap[item.ID] = false
		}

		duplicateID, outputExists := outputMap[item.Output]
		if outputExists {
			errorList = append(errorList, fmt.Errorf("output \"%s\" with id: %d is duplicated with id: %d",
				item.Output, item.ID, duplicateID))
		} else {
			outputMap[item.Output] = item.ID
		}
	}

	return
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

func (dli *DownloadListItem) Format() string {
	return fmt.Sprintf("[%d] \"%s\" => \"%s\"", dli.ID, dli.Input, dli.Output)
}
