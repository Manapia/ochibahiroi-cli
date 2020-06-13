package list

import (
	"github.com/Manapia/ochibahiroi-cli/list"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var reformatCmd *cobra.Command

func initReformatCmd() {
	reformatCmd = &cobra.Command{
		Use:   "reformat list_file",
		Short: "Reformat style of list file.",
		Args:  cobra.ExactArgs(1),
		Run:   reformatCmdRun,
	}
}

func reformatCmdRun(_ *cobra.Command, args []string) {
	specifiedFilePath := args[0]
	listFile := list.DownloadListFile{
		FilePath: specifiedFilePath,
	}
	if !listFile.FileExists() {
		log.Fatalln("The specified file could not be found.")
	}
	if err := listFile.LoadFile(); err != nil {
		log.Fatalf("Failed to load file.\n%v", err)
	}

	jsonBytes, err := listFile.ToJson()
	if err != nil {
		log.Fatalf("Failed to encode data into JSON.\n%v", err)
	}

	if err := ioutil.WriteFile(specifiedFilePath, jsonBytes, 0644); err != nil {
		log.Fatalf("Failed to write JSON data.\n%v", err)
	}
}
