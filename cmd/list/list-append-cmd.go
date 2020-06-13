package list

import (
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/list"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var appendCmd *cobra.Command

func initAppendCmd() {
	appendCmd = &cobra.Command{
		Use:   "append input_file [flags] [input_file output_file]...",
		Short: "Append to download list file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 && len(args)%2 == 0 {
				return fmt.Errorf("arguments must be at least 3 and odd")
			}
			return nil
		},
		Run: initAppendCmdRun,
	}
}

func initAppendCmdRun(_ *cobra.Command, args []string) {
	specifiedFilePath := args[0]
	fileData := list.DownloadListFile{
		FilePath: specifiedFilePath,
	}
	if !fileData.FileExists() {
		log.Fatalln("The specified file could not be found.")
	}

	if err := fileData.LoadFile(); err != nil {
		log.Fatalf("Failed to load file.\n%v", err)
	}

	startID := 0
	for _, item := range fileData.List {
		if startID < item.ID {
			startID = item.ID
		}
	}
	if startID > startID+len(fileData.List) {
		log.Fatalf("id overflow")
	}
	currentID := startID + 1

	for i := 1; i < len(args)-1; i += 2 {
		input := args[i]
		output := args[i+1]

		fileData.List[currentID] = &list.DownloadListItem{
			ID:     currentID,
			Input:  input,
			Output: output,
		}

		currentID++
	}

	jsonBytes, err := fileData.ToJson()
	if err != nil {
		log.Fatalf("Failed to encode data into JSON.\n%v", err)
	}

	if err := ioutil.WriteFile(specifiedFilePath, jsonBytes, 0644); err != nil {
		log.Fatalf("Failed to write JSON data.\n%v", err)
	}
}
