package list

import (
	"fmt"
	"github.com/Manapia/ochibahiroi-cli/downloader"
	"github.com/Manapia/ochibahiroi-cli/list"
	"github.com/Manapia/ochibahiroi-cli/sanitizer"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

type runCmdConfigs struct {
	DryRun        bool
	RelativePath  string
	parallels     int
	showProgress  bool
	makeOutputDir bool
	headerFile    string
}

var runCmd *cobra.Command
var runCmdConfig runCmdConfigs

func initRunCmd() {
	runCmd = &cobra.Command{
		Use:   "run [flags] list_file",
		Short: "Load and run the download list.",
		Args:  cobra.ExactArgs(1),
		Run:   runCmdRun,
	}

	runCmd.Flags().BoolVarP(&runCmdConfig.DryRun, "dry-run", "d", false, "Display a list of downloaded items.")
	runCmd.Flags().StringVar(&runCmdConfig.RelativePath, "relative-path", "", "Base path if the output is not an absolute path.")
	runCmd.Flags().IntVarP(&runCmdConfig.parallels, "parallels", "m", 2, "Number of concurrent downloads.")
	runCmd.Flags().BoolVarP(&runCmdConfig.showProgress, "progress", "p", true, "Show progress bars.")
	runCmd.Flags().BoolVar(&runCmdConfig.makeOutputDir, "make-output", false, "If the destination folder does not exist, it will be created.")
	runCmd.Flags().StringVar(&runCmdConfig.headerFile, "header-file", "", "Read the download request header from the file.")
}

func runCmdRun(_ *cobra.Command, args []string) {
	specifiedFilePath := args[0]
	listFile := list.DownloadListFile{
		FilePath: specifiedFilePath,
	}
	if !listFile.FileExists() {
		log.Fatalln("The specified download list file could not be found.")
	}
	if err := listFile.LoadFile(); err != nil {
		log.Fatalf("Failed to load download list file.\n%v", err)
	}
	if errList := listFile.Validation(); len(errList) != 0 {
		mes := fmt.Sprintf("%d issues found in download list.\n", len(errList))
		for _, err := range errList {
			mes += fmt.Sprintf("%v\n", err)
		}
		log.Fatalln(mes)
	}

	jobList := make([]*downloader.Job, 0, len(listFile.List))
	for _, item := range listFile.List {
		if runCmdConfig.RelativePath != "" && !filepath.IsAbs(item.Output) {
			item.Output = filepath.Join(runCmdConfig.RelativePath, item.Output)
		}

		newJob := &downloader.Job{
			Url:      item.Input,
			SavePath: sanitizer.SanitizeFilePath(item.Output),
		}
		jobList = append(jobList, newJob)
	}

	if runCmdConfig.DryRun {
		for _, job := range jobList {
			fmt.Println(job.ToDisplayString())
		}
		return
	}

	// running phase
	for _, job := range jobList {
		outputPathStat, err := os.Stat(filepath.Dir(job.SavePath))
		if err != nil {
			if runCmdConfig.makeOutputDir {
				if err := os.MkdirAll(filepath.Dir(job.SavePath), 0755); err != nil {
					log.Fatalln(err)
				}
			}
		} else if !outputPathStat.IsDir() {
			log.Fatalf("Output path %s is not directory", job.SavePath)
		}
	}

	option := downloader.DownloadOption{
		Parallels:    runCmdConfig.parallels,
		ShowProgress: runCmdConfig.showProgress,
	}

	if runCmdConfig.headerFile != "" {
		headerData, err := loadHeaderFile()
		if err != nil {
			log.Fatalf("Failed to process header file.\n%v", err)
		}
		option.Header = headerData
	}

	downloader.Run(jobList, option)
}

func loadHeaderFile() (map[string]string, error) {
	_, err := os.Stat(runCmdConfig.headerFile)
	if err != nil {
		return nil, fmt.Errorf("header file not found.\n%v", err)
	}

	f, err := os.Open(runCmdConfig.headerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open the header file.\n%v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close the header file.\n%v", err)
		}
	}()

	headerData, err := downloader.ParseHeaderString(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the file\n%v", err)
	}
	return headerData, nil
}
