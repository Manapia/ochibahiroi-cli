package cmd

import (
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/cmd/list"
	"github.com/Manapia/ochibahiroi-cli/downloader"
	"github.com/Manapia/ochibahiroi-cli/sanitizer"
	"github.com/spf13/cobra"
	"log"
	"os"
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
	headerFile    string
}

var rootCmd *cobra.Command
var rootOptions = &rootOption{}

func init() {
	rootCmd = &cobra.Command{
		Use:   "ochibahiroi",
		Short: "Ochibahiroi is a downloader that executes multiple files in the list in parallel",
		Run:   rootRun,
	}

	rootCmd.Version = "1.1.3"

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
	rootCmd.Flags().StringVar(&rootOptions.headerFile, "header-file", "", "Read the download request header from the file.")

	rootCmd.AddCommand(list.Cmd)
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func rootRun(_ *cobra.Command, _ []string) {
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

	for _, job := range jobs {
		job.SavePath = sanitizer.SanitizeFilePath(job.SavePath)
	}

	if rootOptions.dryRun {
		showDryRun(jobs)
		return
	}

	outputPathStat, err := os.Stat(rootOptions.outputPath)
	if err != nil {
		if rootOptions.makeOutputDir {
			if err := os.MkdirAll(rootOptions.outputPath, 0755); err != nil {
				log.Fatalln(err)
			}
		}
	} else if !outputPathStat.IsDir() {
		log.Fatalf("Output path %s is not directory.", rootOptions.outputPath)
	}

	option := downloader.DownloadOption{
		Parallels:    rootOptions.parallels,
		ShowProgress: rootOptions.showProgress,
	}

	if rootOptions.headerFile != "" {
		headerData, err := loadHeaderFile()
		if err != nil {
			log.Fatalf("Failed to process header file.\n%v", err)
		}
		option.Header = headerData
	}

	downloader.Run(jobs, option)
}

func showDryRun(jobs []*downloader.Job) {
	for _, job := range jobs {
		fmt.Printf("%s => %s\n", job.Url, job.SavePath)
	}
}

func loadHeaderFile() (map[string]string, error) {
	_, err := os.Stat(rootOptions.headerFile)
	if err != nil {
		return nil, fmt.Errorf("header file not found.\n%v", err)
	}

	f, err := os.Open(rootOptions.headerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open the header file.\n%v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close the header file.\n%v", err)
		}
	}()

	headerData, err := downloader.ParseHeaderString(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the file\n%v", err)
	}
	return headerData, nil
}
