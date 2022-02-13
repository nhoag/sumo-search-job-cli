package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// jobProcessFullCmd represents the jobProcessFull command
var jobProcessFullCmd = &cobra.Command{
	Use:   "jobProcessFull",
	Short: "Perform the full process for a Sumo Logic Search Job",
	Long: `The jobProcessFull command will create a search job, poll for
	completion, fetch the results, and delete a Sumo Logic Search Job via the
	Search Job API.`,
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobProcessFull\n", time.Now().UnixNano())
		}
		validateProcessFull()
		executeProcessFull(cmd, args)
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobProcessFull\n", time.Now().UnixNano())
		}
	},
}

func validateProcessFull() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobProcessFull::validateProcessFull()\n", time.Now().UnixNano())
	}
	validateJobCreate()
	validateStatusCheck()
	validateJobResults()
	validateDelete()
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobProcessFull::validateProcessFull()\n", time.Now().UnixNano())
	}
}

func executeProcessFull(cmd *cobra.Command, args []string) {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobProcessFull::executeProcessFull()\n", time.Now().UnixNano())
	}
	_, jobId := executeSearchJob(buildPayload(cmd, args))
	// Add Job ID as first arg for subsequent function calls.
	args = append([]string{jobId}, args...)
	executeJobResults(cmd, args)
	executeDelete(cmd, args)
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobProcessFull::executeProcessFull()\n", time.Now().UnixNano())
	}
}

func init() {
	rootCmd.AddCommand(jobProcessFullCmd)

	jobProcessFullCmd.Flags().StringVarP(&JobOpt, "job", "j", "", "Search job definition")
	jobProcessFullCmd.Flags().StringVarP(&JobFileOpt, "job-file", "J", "", "Path to file with full search job definition")

	jobProcessFullCmd.Flags().StringVarP(&QueryOpt, "query", "q", "", "Search query")
	jobProcessFullCmd.Flags().StringVarP(&QueryFileOpt, "query-file", "Q", "", "Path to file with search query")

	jobProcessFullCmd.Flags().StringVarP(&DurationOpt, "duration", "d", "", "Size of time span relative to now (e.g. -3h)")
	// @todo: Add FromTimeMillis option: milliseconds since epoch.
	jobProcessFullCmd.Flags().StringVarP(&FromTimeOpt, "from", "f", "", "Search window start time (e.g. 2017-07-16T00:00:00)")
	// @todo: Add ToTimeMillis option: milliseconds since epoch.
	jobProcessFullCmd.Flags().StringVarP(&ToTimeOpt, "to", "t", "", "Search window end time (e.g. 2017-07-16T00:00:00)")
	jobProcessFullCmd.Flags().StringVarP(&TimeZoneOpt, "timezone", "z", "UTC", "Timezone to use for search window")
	jobProcessFullCmd.Flags().BoolP("by-receipt-time", "b", false, "Use receipt-time instead of log message timestamps")
	jobProcessFullCmd.Flags().StringVarP(&AutoParsingModeOpt, "auto-parse", "A", "", "Specify auto-parsing mode to use ('intelligent' automatically runs field extraction rules)")

	jobProcessFullCmd.Flags().BoolP("records", "r", false, "Retrieve records only")
	jobProcessFullCmd.Flags().BoolP("messages", "m", false, "Retrieve messages only")
	jobProcessFullCmd.Flags().BoolP("all", "a", true, "Retrieve all paginated results")
	jobProcessFullCmd.Flags().Int32VarP(&LimitOpt, "limit", "l", 100, "Specify pagination limit")
	jobProcessFullCmd.Flags().Int32VarP(&OffsetOpt, "offset", "o", 0, "Specify pagination offset")
	jobProcessFullCmd.Flags().BoolP("poll", "p", true, "Poll for status until search job is complete")
	jobProcessFullCmd.Flags().Int32VarP(&SleepSecondsOpt, "sleep", "Z", 1, "Specify sleep seconds")
}
