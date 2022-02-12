/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nhoag/sumo-search-job-cli/client"
	openapi "github.com/nhoag/sumologic-search-job-client-go"

	"github.com/spf13/cobra"
)

// jobStatusCheckCmd represents the jobStatusCheck command
var jobStatusCheckCmd = &cobra.Command{
	Use:   "jobStatusCheck JOB_ID",
	Short: "Check the status for a Sumo Logic Search Job",
	Long: `The jobStatusCheck command will check the status of a Sumo Logic
	Search Job via the Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobStatusCheck\n", time.Now().UnixNano())
		}
		executeStatusCheck(cmd, args)
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobStatusCheck\n", time.Now().UnixNano())
		}
	},
}

func validateStatusCheck() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobStatusCheck::validateStatusCheck()\n", time.Now().UnixNano())
	}

	// Add validation logic here
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobStatusCheck::validateStatusCheck()\n", time.Now().UnixNano())
	}
}

func executeStatusCheck(cmd *cobra.Command, args []string) *openapi.SearchJobState {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobStatusCheck::executeStatusCheck()\n", time.Now().UnixNano())
	}
	poll, _ := cmd.Flags().GetBool("poll")
	var status *openapi.SearchJobState
	for {
		status = client.GetSearchJobStatus(args[0])
		if !QuietOpt {
			fmt.Fprintf(
				os.Stderr,
				"Status:\t\t%s\nMessage Count:\t%s\nRecord Count:\t%s\n",
				*status.State,
				strconv.FormatInt(int64(*status.MessageCount), 10),
				strconv.FormatInt(int64(*status.RecordCount), 10),
			)
		}
		if !poll ||
			*status.State == "DONE GATHERING RESULTS" ||
			*status.State == "CANCELLED" ||
			*status.State == "FORCE PAUSED" {
			break
		}
		if VerboseOpt {
			jsonStatus, err := json.Marshal(status)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}
			fmt.Fprintf(os.Stderr, "STATUS PAYLOAD: %s\n", string(jsonStatus))
			fmt.Fprintf(os.Stderr, "%d\tSLEEP SECONDS:\t%d\n", time.Now().UnixNano(), SleepSecondsOpt)
		}
		time.Sleep(time.Duration(SleepSecondsOpt) * time.Second)
	}
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobStatusCheck::executeStatusCheck()\n", time.Now().UnixNano())
	}
	return status
}

func init() {
	rootCmd.AddCommand(jobStatusCheckCmd)
	jobStatusCheckCmd.Flags().BoolP("poll", "p", false, "Poll for status until search job is complete")
	jobStatusCheckCmd.Flags().Int32VarP(&SleepSecondsOpt, "sleep", "Z", 1, "Specify sleep seconds")
}
