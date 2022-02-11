/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// jobProcessFullCmd represents the jobProcessFull command
// @todo: Make sure command descriptions are accurate.
var jobProcessFullCmd = &cobra.Command{
	Use:   "jobProcessFull",
	Short: "Perform the full process for a Sumo Logic Search Job",
	Long: `The jobProcessFull command will create a search job, poll for
	completion, fetch the results, and delete a Sumo Logic Search Job via the
	Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateProcessFull()
		executeProcessFull(cmd, args)
	},
}

// @todo: Each command should define it's own validate, which can then all be called here.
func validateProcessFull() {
	validateJobCreate()
	validateStatusCheck()
	validateJobResults()
	validateDelete()
}

func executeProcessFull(cmd *cobra.Command, args []string) {
	_, jobId := executeSearchJob(buildPayload(cmd, args))
	// Add Job ID as first arg for subsequent function calls.
	args = append([]string{jobId}, args...)
	executeStatusCheck(cmd, args)
	executeJobResults(cmd, args)
	executeDelete(cmd, args)
}

func init() {
	// @todo: Review all options
	rootCmd.AddCommand(jobProcessFullCmd)

	jobProcessFullCmd.Flags().StringVarP(&Job, "job", "j", "", "Search job definition")
	jobProcessFullCmd.Flags().StringVarP(&JobFile, "job-file", "J", "", "Path to file with full search job definition")

	jobProcessFullCmd.Flags().StringVarP(&Query, "query", "q", "", "Search query")
	jobProcessFullCmd.Flags().StringVarP(&QueryFile, "query-file", "Q", "", "Path to file with search query")

	// - relative lookback (-3h)
	jobProcessFullCmd.Flags().StringVarP(&Span, "span", "s", "", "Size of time span relative to now")
	// - from time: 2017-07-16T00:00:00
	// @todo: or milliseconds since epoch.
	jobProcessFullCmd.Flags().StringVarP(&FromTime, "from", "f", "", "Search window start time")
	// - to time: 2017-07-16T00:00:00
	// @todo: or milliseconds since epoch.
	jobProcessFullCmd.Flags().StringVarP(&ToTime, "to", "t", "", "Search window end time")
	// - timezone if ^^^ not millis
	jobProcessFullCmd.Flags().StringVarP(&TimeZone, "timezone", "z", "UTC", "Timezone to use for search window")
	jobProcessFullCmd.Flags().BoolP("by-receipt-time", "b", false, "Use receipt-time instead of log message timestamps")
	// - autoParsingMode (default is 'perfomance', 'intelligent' automatically runs field extraction rules)
	jobProcessFullCmd.Flags().StringVarP(&AutoParsingMode, "auto-parse", "A", "", "Specify auto-parsing mode to use")

	jobProcessFullCmd.Flags().BoolP("records", "r", false, "Retrieve records only")
	jobProcessFullCmd.Flags().BoolP("messages", "m", false, "Retrieve messages only")
	jobProcessFullCmd.Flags().BoolP("all", "a", true, "Retrieve all paginated results (default is first page)")
	jobProcessFullCmd.Flags().Int32VarP(&Limit, "limit", "l", 100, "Specify pagination limit")
	jobProcessFullCmd.Flags().Int32VarP(&Offset, "offset", "o", 0, "Specify pagination offset")
	jobProcessFullCmd.Flags().BoolP("poll", "p", true, "Poll for status until search job is complete")
	jobProcessFullCmd.Flags().Int32VarP(&SleepSeconds, "sleep", "Z", 1, "Specify sleep seconds")
}
