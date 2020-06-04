package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"ochibahiroi-cli/downloader"
)

var rootCmd *cobra.Command
var rootOption = &downloader.DownloadOption{}

func init() {
	rootCmd = &cobra.Command{
		Use:   "ochibahiroi",
		Short: "Ochibahiroi is a downloader that executes multiple files in the list in parallel",
		Run:   rootRun,
	}

	rootCmd.Flags().IntVarP(&rootOption.Parallels, "parallels", "m", 2, "Number of concurrent downloads")
	rootCmd.Flags().BoolVarP(&rootOption.ShowProgress, "progress", "p", true, "Show progress")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	jobs := make([]*downloader.Job, 0, 10)

	option := *rootOption

	downloader.Run(jobs, option)
}
