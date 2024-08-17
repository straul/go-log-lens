package cmd

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/straul/go-log-lens/internal"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

var (
	logFilePath     string   // 单个文件路径参数
	logFilePaths    []string // 多个文件路径参数
	logDir          string   // 日志目录路径参数
	includeKeywords []string // 日志包含的关键词
	excludeKeywords []string // 日志不包含的关键词
	regexPattern    string   // 正则表达式
	startTime       string   // 日志开始时间（依赖日志内容）
	endTime         string   // 日志结束时间（依赖日志内容）
	levels          []string // 日志级别
	jsonOutput      bool     // 输出格式是否为 JSON
	outputFilePath  string   // 输出文件路径参数
	concurrency     int      // 并发数
)

var rootCmd = &cobra.Command{
	Use:   "loglens",
	Short: "LogLens is a CLI tool to view and filter logs",
	Long:  `LogLens helps you to easily view, filter, and manage log files through a command-line interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			processStartTime = time.Now()
			outputWriter     *os.File
			err              error
		)

		// 判断是否写入文件
		if outputFilePath != "" {
			outputWriter, err = os.Create(outputFilePath)
			if err != nil {
				fmt.Printf("Error creating output file: %v\n", err)
				return
			}
			defer outputWriter.Close()
		}

		// 正则表达式编译
		var regex *regexp.Regexp
		if regexPattern != "" {
			regex, err = regexp.Compile(regexPattern)
			if err != nil {
				fmt.Printf("Invalid regex pattern: %v\n", err)
				return
			}
		}

		var processFunc = func(line string) {
			var start, end *time.Time
			if startTime != "" && endTime != "" {
				s, err := time.Parse(time.DateTime, startTime)
				if err != nil {
					fmt.Printf("Invalid start time format: %v\n", err)
					return
				}
				e, err := time.Parse(time.DateTime, endTime)
				if err != nil {
					fmt.Printf("Invalid end time format: %v\n", err)
					return
				}
				start = &s
				end = &e
			}

			// 应用过滤器
			if internal.FilterLine(line, includeKeywords, excludeKeywords, levels, start, end) {
				// 如果正则匹配存在，应用正则过滤
				if regex != nil && !regex.MatchString(line) {
					return
				}

				if jsonOutput {
					line = fmt.Sprintf("{\"log\": \"%s\"}\n", line)
				} else {
					line = fmt.Sprintf("%s\n", line)
				}

				if outputWriter != nil {
					outputWriter.WriteString(line) // 写入文件
				} else {
					fmt.Print(line) // 直接输出到终端
				}
			}
		}

		// 如果提供了日志目录
		if logDir != "" {
			logFilePaths, err = scanLogDirectory(logDir)
			if err != nil {
				fmt.Printf("Error scanning log directory: %v\n", err)
				return
			}
		}

		if logFilePath != "" {
			bar := createProgressBar(logFilePath)
			err := internal.StreamLogsWithProgress(logFilePath, processFunc, bar)
			if err != nil {
				fmt.Printf("Error reading log file: %v\n", err)
			}
		} else if len(logFilePaths) > 0 {
			bar := createProgressBarForMultipleFiles(logFilePaths)
			err := internal.StreamLogsConcurrentlyWithProgress(logFilePaths, processFunc, bar, concurrency)
			if err != nil {
				fmt.Printf("Error reading log files: %v\n", err)
			}
		} else {
			fmt.Println("Please provide a log file path using the --file flag, a log directory using the --log-dir flag, or a comma-separated list of files using the --files flag.")
		}

		duration := time.Since(processStartTime)
		fmt.Printf("Filtering completed in %d ms\n", duration.Milliseconds())
	},
}

func scanLogDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func createProgressBar(filePath string) *progressbar.ProgressBar {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil
	}
	return progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetDescription("Processing "+filePath),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func createProgressBarForMultipleFiles(filePaths []string) *progressbar.ProgressBar {
	var totalSize int64
	for _, filePath := range filePaths {
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			totalSize += fileInfo.Size()
		}
	}
	return progressbar.NewOptions64(
		totalSize,
		progressbar.OptionSetDescription("Processing multiple files"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func init() {
	rootCmd.Flags().StringVarP(&logFilePath, "file", "f", "", "Path to a single log file")
	rootCmd.Flags().StringSliceVarP(&logFilePaths, "files", "", nil, "Comma-separated list of log file paths")
	rootCmd.Flags().StringVarP(&logDir, "log-dir", "d", "", "Directory containing log files")
	rootCmd.Flags().StringSliceVarP(&includeKeywords, "keywords", "k", nil, "Comma-separated keywords to include logs")
	rootCmd.Flags().StringSliceVarP(&excludeKeywords, "exclude-keywords", "x", nil, "Comma-separated keywords to exclude logs")
	rootCmd.Flags().StringVarP(&regexPattern, "regex", "r", "", "Regex pattern to filter logs")
	rootCmd.Flags().StringVarP(&startTime, "start", "s", "", "Start time for log filtering (format: YYYY-MM-DD HH:MM:SS)")
	rootCmd.Flags().StringVarP(&endTime, "end", "e", "", "End time for log filtering (format: YYYY-MM-DD HH:MM:SS)")
	rootCmd.Flags().StringSliceVarP(&levels, "levels", "l", nil, "Comma-separated log levels to filter logs (e.g., ERROR,WARNING,INFO)")
	rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output logs in JSON format")
	rootCmd.Flags().StringVarP(&outputFilePath, "output-file", "o", "", "Path to output file for filtered logs")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", runtime.NumCPU(), "Number of concurrent goroutines")
}

func Execute() error {
	//runtime.GOMAXPROCS(concurrency)
	return rootCmd.Execute()
}
