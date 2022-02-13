package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nhoag/sumo-search-job-cli/client"
	"github.com/spf13/cobra"
)

var (
	LimitOpt        int32
	OffsetOpt       int32
	SleepSecondsOpt int32
)

// jobResultsGetCmd represents the jobResultsGet command
var jobResultsGetCmd = &cobra.Command{
	Use:   "jobResultsGet JOB_ID",
	Short: "Fetch the results for a Sumo Logic Search Job",
	Long: `The jobResultsGet command will fetch the results for a Sumo Logic
	Search Job via the Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobResultsGet\n", time.Now().UnixNano())
		}
		executeJobResults(cmd, args)
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobResultsGet\n", time.Now().UnixNano())
		}
	},
}

func validateJobResults() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobResultsGet::validateJobResults()\n", time.Now().UnixNano())
	}
	// validation logic goes here
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobResultsGet::validateJobResults()\n", time.Now().UnixNano())
	}
}

func executeJobResults(cmd *cobra.Command, args []string) {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobResultsGet::executeJobResults()\n", time.Now().UnixNano())
	}
	status := executeStatusCheck(cmd, args)

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No jobId specified!")
		os.Exit(1)
	}
	jobId := args[0]
	all, _ := cmd.Flags().GetBool("all")
	messagesOnly, _ := cmd.Flags().GetBool("messages")
	recordsOnly, _ := cmd.Flags().GetBool("records")
	msgOffset := OffsetOpt
	recOffset := OffsetOpt
	if !recordsOnly && *status.MessageCount > int32(0) {
		for {
			messages := client.GetSearchJobMessages(jobId, LimitOpt, msgOffset)
			messagesJson, _ := json.MarshalIndent(messages, "", "    ")
			fmt.Println(string(messagesJson))
			msgOffset = msgOffset + LimitOpt
			if !all || msgOffset > *status.MessageCount {
				break
			}
			time.Sleep(time.Duration(SleepSecondsOpt) * time.Second)
		}
	}
	if !messagesOnly && *status.RecordCount > int32(0) {
		for {
			records := client.GetSearchJobRecords(jobId, LimitOpt, recOffset)
			recordsJson, _ := json.MarshalIndent(records, "", "    ")
			fmt.Println(string(recordsJson))
			recOffset = recOffset + LimitOpt
			if !all || recOffset > *status.RecordCount {
				break
			}
			time.Sleep(time.Duration(SleepSecondsOpt) * time.Second)
		}
	}
	if !QuietOpt && *status.MessageCount == int32(0) && *status.RecordCount == int32(0) {
		fmt.Fprintf(os.Stderr, "No results for the specified search\n")
	}
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobResultsGet::executeJobResults()\n", time.Now().UnixNano())
	}
}

func init() {
	rootCmd.AddCommand(jobResultsGetCmd)
	jobResultsGetCmd.Flags().BoolP("records", "r", false, "Retrieve records only")
	jobResultsGetCmd.Flags().BoolP("messages", "m", false, "Retrieve messages only")
	jobResultsGetCmd.Flags().BoolP("all", "a", false, "Retrieve all paginated results (default is first page)")
	jobResultsGetCmd.Flags().Int32VarP(&LimitOpt, "limit", "l", 100, "Specify pagination limit")
	jobResultsGetCmd.Flags().Int32VarP(&OffsetOpt, "offset", "o", 0, "Specify pagination offset")
	jobResultsGetCmd.Flags().Int32VarP(&SleepSecondsOpt, "sleep", "Z", 1, "Specify sleep seconds")
	jobResultsGetCmd.Flags().BoolP("poll", "p", true, "Poll for status until search job is complete")
	// @todo: Add format options
}
