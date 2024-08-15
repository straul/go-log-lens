package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/straul/go-log-lens/internal"
	"strings"
)

var (
	logFilePath  string
	keywords     string
	regexPattern string
)

var rootCmd = &cobra.Command{
	Use:   "log-lens",
	Short: "LogLens is a CLI tool to view and filter logs",
	Long:  `LogLens helps you to easily view, filter, and manage log files through a command-line interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		if logFilePath == "" {
			fmt.Println("Please provide a log file path using the --file flat.")
			return
		}

		logs, err := internal.ReadLogs(logFilePath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		if keywords != "" {
			logs = internal.FilterLogs(logs, strings.Split(keywords, ","))
		}

		if regexPattern != "" {
			logs, err = internal.FilterLogsByRegex(logs, regexPattern)
			if err != nil {
				fmt.Printf("Error applying regex: %v\n", err)
				return
			}
		}

		for _, line := range logs {
			fmt.Println(line)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&logFilePath, "file", "f", "", "Path to the log file")
	rootCmd.Flags().StringVarP(&keywords, "keywords", "k", "", "Comma-separated keywords to filter logs")
	rootCmd.Flags().StringVarP(&regexPattern, "regex", "r", "", "Regex pattern to filter logs")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
