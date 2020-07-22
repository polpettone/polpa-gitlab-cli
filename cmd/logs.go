package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func NewLogsCmd(apiClient APIClient) *cobra.Command {
	return &cobra.Command{
		Use:
		"logs <jobID> or logs -l",
		Short: "logs -l shows the log of the last failed job",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleLogsCommand(cmd, args, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(stdout)
		},
	}
}

func handleLogsCommand(cobraCommand *cobra.Command, args []string, apiClient APIClient) (string, error) {
	lastFailed, _ := cobraCommand.Flags().GetBool("lastFailed")
	if (len(args) < 1 || args[0] == "") && lastFailed == false {
		return "", errors.New(
			"you need to inform a job id or set -l flag to see the logs from the last failed job")
	}
	if lastFailed {
		return apiClient.getLastFailLog()
	} else {
		jobID := args[0]
		return apiClient.getLog(jobID)
	}
}

func init() {
	logsCmd := NewLogsCmd(gitlabAPIClient)
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolP(
		"lastFailed",
		"l",
		false,
		"Shows the logs from the last failed job")
}
