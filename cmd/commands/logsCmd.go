package commands

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func NewLogsCmd(apiClient *adapter.App) *cobra.Command {
	return &cobra.Command{
		Use:
		"logs -> interactive mode| logs <jobID> -> logs of job | logs -l -> logs of last failed job",
		Short: "when no job id " +
			"is given, interactive mode started to choose a pipeline, then a job to see logs or use flag -l to see log of the last failed job",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleLogsCommand(cmd, args, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(stdout)
		},
	}
}

func handleLogsCommand(cobraCommand *cobra.Command, args []string, apiClient *adapter.App) (string, error) {
	lastFailed, _ := cobraCommand.Flags().GetBool("lastFailed")
	toFile, _ := cobraCommand.Flags().GetString("toFile")

	var jobID string

	if (len(args) < 1 || args[0] == "") && lastFailed == false {

		pipelines, _ := apiClient.GetPipelines("",  20)
		pipeline, _ := adapter.ShowPipelineSelectionPrompt(pipelines)
		pipelineId := strconv.Itoa(pipeline.Id)

		jobs, _ := apiClient.GetJobs(pipelineId)
		job, _ := adapter.ShowJobSelectionPrompt(jobs)
		jobID = strconv.Itoa(job.Id)

		if toFile != "" {
			log, err := apiClient.GetLog(jobID)
			if err != nil {
				return "", err
			}
			err = writeToFile(log, toFile)
			if err != nil {
				return "", err
			}
			return "Written to " + toFile, nil
		}

		return apiClient.GetLog(jobID)
	} else {
		var log string
		var err error
		if lastFailed {
			log, err =  apiClient.GetLastFailLog()
		} else {
			jobID = args[0]
			log, err = apiClient.GetLog(jobID)
		}

		if toFile != "" {
			if err != nil {
				return "", err
			}
			err = writeToFile(log, toFile)
			if err != nil {
				return "", err
			}
			return "Written to " + toFile, nil
		}

		return log, err
	}
}

func init() {
	logsCmd := NewLogsCmd(adapter.NewApp())
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP(
		"lastFailed",
		"l",
		false,
		"Shows the logs from the last failed job")

	logsCmd.Flags().StringP(
		"toFile",
		"t",
		"",
		"Writes the logs to the given file")

}

func writeToFile(content string, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	return err
}
