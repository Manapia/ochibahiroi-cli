package cmd

import (
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/downloader"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type rootOption struct {
	parallels     int
	showProgress  bool
	dryRun        bool
	url           string
	numberStart   int
	numberEnd     int
	step          int
	resetNumber   bool
	outputPath    string
	makeOutputDir bool
	openOutputDir bool
}

var rootCmd *cobra.Command
var rootOptions = &rootOption{}

func init() {
	rootCmd = &cobra.Command{
		Use:   "ochibahiroi",
		Short: "Ochibahiroi is a downloader that executes multiple files in the list in parallel",
		Run:   rootRun,
	}

	rootCmd.Flags().IntVarP(&rootOptions.parallels, "parallels", "m", 2, "Number of concurrent downloads.")
	rootCmd.Flags().BoolVarP(&rootOptions.showProgress, "progress", "p", true, "Show progress bars.")
	rootCmd.Flags().BoolVarP(&rootOptions.dryRun, "dry-run", "d", false, "Show a list of files to download.")
	rootCmd.Flags().StringVarP(&rootOptions.url, "url", "u", "", "URL of the download source.")
	rootCmd.Flags().IntVar(&rootOptions.numberStart, "start", 0, "The first number in the sequence.")
	rootCmd.Flags().IntVar(&rootOptions.numberEnd, "end", 0, "The last number in the sequence.")
	rootCmd.Flags().IntVar(&rootOptions.step, "step", 1, "Number of steps in a sequence.")
	rootCmd.Flags().BoolVarP(&rootOptions.resetNumber, "reset-number", "r", false, "Use numbers that start with 1 instead of the original filename.")
	rootCmd.Flags().StringVarP(&rootOptions.outputPath, "output-path", "o", "./", "The output destination of the downloaded files.")
	rootCmd.Flags().BoolVar(&rootOptions.makeOutputDir, "make-output", false, "If the destination folder does not exist, it will be created.")
	rootCmd.Flags().BoolVar(&rootOptions.openOutputDir, "open", false, "After the download is complete, open the destination folder.")
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	builder := downloader.JobListBuilder{}
	builder.SetUrl(rootOptions.url)
	builder.SetSavePath(rootOptions.outputPath)
	builder.SetStart(rootOptions.numberStart)
	builder.SetEnd(rootOptions.numberEnd)
	builder.SetStep(rootOptions.step)
	builder.SetUserIncrementalCount(rootOptions.resetNumber)

	jobs, err := builder.Build()
	if err != nil {
		log.Fatalln(err)
	}

	if rootOptions.dryRun {
		showDryRun(jobs)
		return
	}

	outputPathStat, err := os.Stat(rootOptions.outputPath)
	if err != nil {
		if rootOptions.makeOutputDir {
			if err := os.MkdirAll(rootOptions.outputPath, 0666); err != nil {
				log.Fatalln(err)
			}
		}
	} else if !outputPathStat.IsDir() {
		log.Fatalf("output path %s is not directory", rootOptions.outputPath)
	}

	option := downloader.DownloadOption{
		Parallels:    rootOptions.parallels,
		ShowProgress: rootOptions.showProgress,
	}

	downloader.Run(jobs, option)

	if rootOptions.openOutputDir {
		showOutputDirectory()
	}
}

func showDryRun(jobs []*downloader.Job) {
	for _, job := range jobs {
		fmt.Printf("%s => %s\n", job.Url, job.SavePath)
	}
}

func showOutputDirectory() {
	var err error

	abs, err := filepath.Abs(rootOptions.outputPath)

	if err == nil {
		switch runtime.GOOS {
		case "darwin":
			err = exec.Command("open", abs).Start()
		case "linux":
			err = exec.Command("xdg-open", abs).Start()
		case "windows":
			cmd := exec.Command(`explorer`, `/select,`, abs)
			err = cmd.Run()
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
