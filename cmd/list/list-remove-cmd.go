package list

import (
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/list"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var removeCmd *cobra.Command

type removeCmdConfigs struct {
	ID     int
	Input  string
	Output string
	DryRun bool
}

var removeCmdConfig removeCmdConfigs

func initListRemoveCmd() {
	removeCmd = &cobra.Command{
		Use:   "remove [flags] list_file",
		Short: "Remove item from list of JSON data.",
		Args:  cobra.ExactArgs(1),
		Run:   removeCmdRun,
	}

	removeCmd.Flags().IntVarP(&removeCmdConfig.ID, "id", "i", 0, "ID of the item to delete.")
	removeCmd.Flags().StringVar(&removeCmdConfig.Input, "input", "", "Delete the item with the value specified in input.")
	removeCmd.Flags().StringVar(&removeCmdConfig.Output, "output", "", "Delete the item with the value specified in output.")
	removeCmd.Flags().BoolVarP(&removeCmdConfig.DryRun, "dry-run", "d", false, "Display the item to be deleted.")
}

func removeCmdRun(cmd *cobra.Command, args []string) {
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

	idSpecified := cmd.Flags().Changed("id")
	inputSpecified := cmd.Flags().Changed("input")
	outputSpecified := cmd.Flags().Changed("output")

	removeItems := make([]*list.DownloadListItem, 0, 4)

	for key, item := range listFile.List {
		removed := false

		if idSpecified && item.ID == removeCmdConfig.ID {
			removeItems = append(removeItems, item)
			removed = true
		} else if inputSpecified && item.Input == removeCmdConfig.Input {
			removeItems = append(removeItems, item)
			removed = true
		} else if outputSpecified && item.Output == removeCmdConfig.Output {
			removeItems = append(removeItems, item)
			removed = true
		}

		if removed {
			delete(listFile.List, key)
		}
	}

	if removeCmdConfig.DryRun {
		for _, removeItem := range removeItems {
			fmt.Println(removeItem.Format())
		}
		return
	}

	jsonBytes, err := listFile.ToJson()
	if err != nil {
		log.Fatalf("Failed to encode data into JSON.\n%v", err)
	}

	if err := ioutil.WriteFile(specifiedFilePath, jsonBytes, 0644); err != nil {
		log.Fatalf("Failed to write JSON data.\n%v", err)
	}
}
