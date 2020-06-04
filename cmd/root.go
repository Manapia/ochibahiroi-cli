package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "ochibahiroi",
		Short: "Ochibahiroi is a downloader that executes multiple files in the list in parallel",
		Run:   rootRun,
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func rootRun(cmd *cobra.Command, args []string) {

}
