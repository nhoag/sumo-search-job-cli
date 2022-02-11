/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nhoag/sumo-search-job-cli/client"

	"github.com/spf13/cobra"
)

// jobStatusCheckCmd represents the jobStatusCheck command
var jobStatusCheckCmd = &cobra.Command{
	Use:   "jobStatusCheck JOB_ID",
	Short: "Check the status for a Sumo Logic Search Job",
	Long: `The jobStatusCheck command will check the status of a Sumo Logic
	Search Job via the Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeStatusCheck(cmd, args)
	},
}

func validateStatusCheck() {
	// Add validation logic here
}

func executeStatusCheck(cmd *cobra.Command, args []string) {
	poll, _ := cmd.Flags().GetBool("poll")
	for {
		status := client.GetSearchJobStatus(args[0])
		fmt.Fprintf(
			os.Stderr,
			"Status:\t\t%s\nMessage Count:\t%s\nRecord Count:\t%s\n",
			*status.State,
			strconv.FormatInt(int64(*status.MessageCount), 10),
			strconv.FormatInt(int64(*status.RecordCount), 10),
		)
		if !poll ||
			*status.State == "DONE GATHERING RESULTS" ||
			*status.State == "CANCELLED" ||
			*status.State == "FORCE PAUSED" {
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func init() {
	rootCmd.AddCommand(jobStatusCheckCmd)
	jobStatusCheckCmd.Flags().BoolP("poll", "p", false, "Poll for status until search job is complete")
}
