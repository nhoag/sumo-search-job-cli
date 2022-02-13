package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	DurationMinutes int32
	IntervalSeconds int32
	RequestCount    int32
)

// jobKeepAliveCmd represents the jobKeepAlive command
var jobKeepAliveCmd = &cobra.Command{
	Use:   "jobKeepAlive JOB_ID",
	Short: "Issue periodic keep-alive job status request",
	Long:  `Keep a Search Job alive by issuing periodic status requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobKeepAlive\n", time.Now().UnixNano())
		}
		executeKeepAlive(cmd, args)
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobKeepAlive\n", time.Now().UnixNano())
		}
	},
}

func validateKeepAlive() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobKeepAlive::validateKeepAlive()\n", time.Now().UnixNano())
	}
	// Add validation logic here
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobKeepAlive::validateKeepAlive()\n", time.Now().UnixNano())
	}
}

func executeKeepAlive(cmd *cobra.Command, args []string) {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobKeepAlive::executeKeepAlive()\n", time.Now().UnixNano())
	}
	forever, _ := cmd.Flags().GetBool("forever")
	iterations := int32(1)
	start := time.Now().Unix()
	for {
		executeStatusCheck(cmd, args)
		if !forever &&
			(iterations >= RequestCount ||
				time.Now().Unix()-start > int64(DurationMinutes)*60) {
			break
		}
		iterations = iterations + int32(1)
		time.Sleep(time.Duration(IntervalSeconds) * time.Second)
	}

	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobKeepAlive::executeKeepAlive()\n", time.Now().UnixNano())
	}
}

func init() {
	rootCmd.AddCommand(jobKeepAliveCmd)
	jobKeepAliveCmd.Flags().Int32VarP(&IntervalSeconds, "interval", "i", 30, "Keep-alive interval in seconds")
	jobKeepAliveCmd.Flags().Int32VarP(&DurationMinutes, "duration", "k", 30, "Keep-alive duration in minutes")
	jobKeepAliveCmd.Flags().Int32VarP(&RequestCount, "count", "c", 10, "Keep-alive request count")
	jobKeepAliveCmd.Flags().BoolP("forever", "f", false, "Issue keep-alive requests indefinitely")
}
