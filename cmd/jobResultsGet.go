/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nhoag/sumo-search-job-cli/client"
	"github.com/spf13/cobra"
)

var (
	Limit        int32
	Offset       int32
	SleepSeconds int32
)

// jobResultsGetCmd represents the jobResultsGet command
var jobResultsGetCmd = &cobra.Command{
	Use:   "jobResultsGet JOB_ID",
	Short: "Fetch the results for a Sumo Logic Search Job",
	Long: `The jobResultsGet command will fetch the results for a Sumo Logic
	Search Job via the Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeJobResults(cmd, args)
	},
}

func validateJobResults() {
	// validation logic goes here
}

func executeJobResults(cmd *cobra.Command, args []string) {
	jobId := args[0]
	status := client.GetSearchJobStatus(jobId)
	all, _ := cmd.Flags().GetBool("all")
	messagesOnly, _ := cmd.Flags().GetBool("messages")
	recordsOnly, _ := cmd.Flags().GetBool("records")
	msgOffset := Offset
	recOffset := Offset
	if !recordsOnly && *status.MessageCount > int32(0) {
		for {
			messages := client.GetSearchJobMessages(jobId, Limit, msgOffset)
			messagesJson, _ := json.MarshalIndent(messages, "", "    ")
			fmt.Println(string(messagesJson))
			msgOffset = msgOffset + Limit
			if !all || msgOffset > *status.MessageCount {
				break
			}
			time.Sleep(time.Duration(SleepSeconds) * time.Second)
		}
	}
	if !messagesOnly && *status.RecordCount > int32(0) {
		for {
			records := client.GetSearchJobRecords(jobId, Limit, recOffset)
			recordsJson, _ := json.MarshalIndent(records, "", "    ")
			fmt.Println(string(recordsJson))
			recOffset = recOffset + Limit
			if !all || recOffset > *status.RecordCount {
				break
			}
			time.Sleep(time.Duration(SleepSeconds) * time.Second)
		}
	}
}

func init() {
	rootCmd.AddCommand(jobResultsGetCmd)
	jobResultsGetCmd.Flags().BoolP("records", "r", false, "Retrieve records only")
	jobResultsGetCmd.Flags().BoolP("messages", "m", false, "Retrieve messages only")
	jobResultsGetCmd.Flags().BoolP("all", "a", false, "Retrieve all paginated results (default is first page)")
	jobResultsGetCmd.Flags().Int32VarP(&Limit, "limit", "l", 50, "Specify pagination limit")
	jobResultsGetCmd.Flags().Int32VarP(&Offset, "offset", "o", 0, "Specify pagination offset")
	jobResultsGetCmd.Flags().Int32VarP(&SleepSeconds, "sleep", "Z", 1, "Specify sleep seconds")
	// @todo: Add format options
}
