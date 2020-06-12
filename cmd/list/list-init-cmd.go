package list

import (
	"bufio"
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/list"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var initCmd *cobra.Command

func initInitCmd() {
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create an initialized download list file",
		Args:  cobra.ExactArgs(1),
		Run:   listInitCmdRun,
	}
}

func listInitCmdRun(_ *cobra.Command, args []string) {
	specifiedFileName := args[0]

	_, err := os.Stat(specifiedFileName)
	if err == nil {
		if !askOverWrite() {
			return
		}
	}

	newData := list.DownloadListData{
		Version: list.DownloadListVersion,
		List:    map[int]*list.DownloadListItem{},
	}

	out, err := newData.ToJson()
	if err != nil {
		log.Fatalf("Failed to encode data into JSON.\n%v", err)
	}

	if err := ioutil.WriteFile(specifiedFileName, out, 0644); err != nil {
		log.Fatalf("Failed to write JSON data.\n%v", err)
	}

	dstPath, err := filepath.Abs(specifiedFileName)
	if err != nil {
		log.Fatalf("Cannot get the path of the written file.")
	}
	fmt.Printf("Saved to file %s\n", dstPath)
}

func askOverWrite() bool {
	fmt.Println("File already exists. Do you want to overwrite and initialize it? (Y/n)")
	prompt := bufio.NewScanner(os.Stdin)
	prompt.Scan()
	return prompt.Text() == "Y"
}
