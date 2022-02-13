package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/nhoag/sumo-search-job-cli/client"
	openapi "github.com/nhoag/sumologic-search-job-client-go"
)

var (
	JobOpt     string
	JobFileOpt string

	QueryOpt     string
	QueryFileOpt string

	FromTimeOpt string
	DurationOpt string
	ToTimeOpt   string
	TimeZoneOpt string

	AutoParsingModeOpt string
)

// jobCreateCmd represents the jobCreate command
var jobCreateCmd = &cobra.Command{
	Use:   "jobCreate",
	Short: "Create a Sumo Logic Search Job",
	Long: `The jobCreate command will initiate a Sumo Logic Search Job via the
	Search Job API. The returned Job ID can be used to query job status and
	fetch search results.`,
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobCreate\n", time.Now().UnixNano())
		}
		validateJobCreate()
		executeSearchJob(buildPayload(cmd, args))
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobCreate\n", time.Now().UnixNano())
		}
	},
}

type JobDefinition struct {
	Query           string `json:"query,omitempty"`
	From            string `json:"from,omitempty"`
	To              string `json:"to,omitempty"`
	Timezone        string `json:"timeZone,omitempty"`
	ByReceiptTime   bool   `json:"byReceiptTime,omitempty"`
	AutoParsingMode string `json:"autoParsingMode,omitempty"`
}

func buildPayload(cmd *cobra.Command, args []string) JobDefinition {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobCreate::buildPayload()\n", time.Now().UnixNano())
	}

	var jobDef JobDefinition
	jobDef.Timezone = "UTC"
	if len(JobOpt) > 0 {
		json.Unmarshal([]byte(JobOpt), &jobDef)
	}
	if len(JobFileOpt) > 0 {
		content, err := ioutil.ReadFile(JobFileOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		json.Unmarshal([]byte(content), &jobDef)
	}
	if len(QueryOpt) > 0 {
		jobDef.Query = QueryOpt
	}
	if len(QueryFileOpt) > 0 {
		content, err := ioutil.ReadFile(QueryFileOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		jobDef.Query = string(content)
	}
	if len(DurationOpt) > 0 {
		duration, _ := time.ParseDuration(DurationOpt)
		jobDef.To = time.Now().UTC().String()
		jobDef.From = time.Now().Add(-duration).UTC().String()
	}
	if len(ToTimeOpt) > 0 {
		jobDef.To = ToTimeOpt
	}
	if len(FromTimeOpt) > 0 {
		jobDef.From = FromTimeOpt
	}
	if len(TimeZoneOpt) > 0 {
		jobDef.Timezone = TimeZoneOpt
	}
	if len(JobOpt) == 0 && len(JobFileOpt) == 0 {
		jobDef.ByReceiptTime, _ = cmd.Flags().GetBool("by-receipt-time")
	}
	if len(AutoParsingModeOpt) > 0 {
		jobDef.AutoParsingMode = AutoParsingModeOpt
	}
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobCreate::buildPayload()\n", time.Now().UnixNano())
	}

	return jobDef
}

func executeSearchJob(jobDef JobDefinition) (*url.URL, string) {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobCreate::executeSearchJob()\n", time.Now().UnixNano())
	}

	searchJobDef := *openapi.NewSearchJobDefinition()
	searchJobDef.SetTo(jobDef.To)
	searchJobDef.SetFrom(jobDef.From)
	searchJobDef.SetTimeZone(jobDef.Timezone)
	searchJobDef.SetQuery(jobDef.Query)

	if VerboseOpt {
		defJson, err := json.Marshal(searchJobDef)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		fmt.Fprintf(os.Stderr, "SEARCH JOB: %s\n", string(defJson))
	}

	location, jobId := client.CreateSearchJob(searchJobDef)
	if !QuietOpt {
		fmt.Fprintf(os.Stderr, "Location:\t%s\nJob ID:\t\t%s\n", location, jobId)
	}
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobCreate::executeSearchJob()\n", time.Now().UnixNano())
	}

	return location, jobId

}

func validateJobCreate() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobCreate::validateJobCreate()\n", time.Now().UnixNano())
	}

	if len(JobOpt) > 0 {
		if len(JobFileOpt) > 0 {
			fmt.Fprintf(os.Stderr, "job-file is not compatible with job")
			os.Exit(1)
		}
		if len(QueryOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query is not compatible with job")
			os.Exit(1)
		}
		if len(QueryFileOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query-file is not compatible with job")
			os.Exit(1)
		}
		if len(FromTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "from is not compatible with job")
			os.Exit(1)
		}
		if len(DurationOpt) > 0 {
			fmt.Fprintf(os.Stderr, "span is not compatible with job")
			os.Exit(1)
		}
		if len(ToTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "to is not compatible with job")
			os.Exit(1)
		}
		if len(AutoParsingModeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "auto-parse is not compatible with job")
			os.Exit(1)
		}
	}
	if len(JobFileOpt) > 0 {
		if len(JobOpt) > 0 {
			fmt.Fprintf(os.Stderr, "job is not compatible with job-file")
			os.Exit(1)
		}
		if len(QueryOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query is not compatible with job")
			os.Exit(1)
		}
		if len(QueryFileOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query-file is not compatible with job")
			os.Exit(1)
		}
		if len(FromTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "from is not compatible with job")
			os.Exit(1)
		}
		if len(DurationOpt) > 0 {
			fmt.Fprintf(os.Stderr, "span is not compatible with job")
			os.Exit(1)
		}
		if len(ToTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "to is not compatible with job")
			os.Exit(1)
		}
		fmt.Printf("%+v", AutoParsingModeOpt)
		if len(AutoParsingModeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "auto-parse is not compatible with job-file")
			os.Exit(1)
		}
	}
	if len(QueryOpt) > 0 {
		if len(QueryFileOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query-file is not compatible with query")
			os.Exit(1)
		}
	}
	if len(QueryFileOpt) > 0 {
		if len(QueryOpt) > 0 {
			fmt.Fprintf(os.Stderr, "query is not compatible with query-file")
			os.Exit(1)
		}
	}
	if len(DurationOpt) > 0 {
		if len(FromTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "from is not compatible with span")
			os.Exit(1)
		}
		if len(ToTimeOpt) > 0 {
			fmt.Fprintf(os.Stderr, "to is not compatible with span")
			os.Exit(1)
		}
		if len(TimeZoneOpt) > 0 {
			fmt.Fprintf(os.Stderr, "timezone is not compatible with span")
			os.Exit(1)
		}
	}

	var toTime time.Time = time.Now().UTC()
	var fromTime time.Time = time.Now().UTC()
	var err error
	var duration time.Duration

	if len(DurationOpt) > 0 {
		duration, err = time.ParseDuration(DurationOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse the provided span: "+DurationOpt)
			os.Exit(1)
		}
		fromTime = time.Now().Add(-duration).UTC()
	}
	if len(ToTimeOpt) > 0 {
		toTime, err = time.Parse("2006-01-02T15:04:05", ToTimeOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse the provided to-time: "+ToTimeOpt)
			os.Exit(1)
		}
	}
	if len(FromTimeOpt) > 0 {
		fromTime, err = time.Parse("2006-01-02T15:04:05", FromTimeOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse the provided from-time: "+FromTimeOpt)
			os.Exit(1)
		}
	}
	if len(JobFileOpt) > 0 {
		var jobDef JobDefinition
		content, err := ioutil.ReadFile(JobFileOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		json.Unmarshal(content, &jobDef)
		toTime, err = time.Parse("2006-01-02T15:04:05", jobDef.To)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse the provided 'to' time: "+jobDef.To)
			os.Exit(1)
		}
		fromTime, err = time.Parse("2006-01-02T15:04:05", jobDef.From)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse the provided 'from' time: "+jobDef.From)
			os.Exit(1)
		}
	}
	if !fromTime.Before(toTime) {
		fmt.Fprintf(os.Stderr, "from "+fromTime.String()+" is not before to "+toTime.String())
		os.Exit(1)
	}

	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobCreate::validateJobCreate()\n", time.Now().UnixNano())
	}
}

func init() {
	rootCmd.AddCommand(jobCreateCmd)

	jobCreateCmd.Flags().StringVarP(&JobOpt, "job", "j", "", "Search job definition")
	jobCreateCmd.Flags().StringVarP(&JobFileOpt, "job-file", "J", "", "Path to file with full search job definition")

	jobCreateCmd.Flags().StringVarP(&QueryOpt, "query", "q", "", "Search query")
	jobCreateCmd.Flags().StringVarP(&QueryFileOpt, "query-file", "Q", "", "Path to file with search query")

	jobCreateCmd.Flags().StringVarP(&DurationOpt, "duration", "d", "", "Size of time span relative to now (e.g. -3h)")
	// @todo: Add FromTimeMillis option: milliseconds since epoch.
	jobCreateCmd.Flags().StringVarP(&FromTimeOpt, "from", "f", "", "Search window start time (e.g. 2017-07-16T00:00:00")
	// @todo: Add ToTimeMillis option: milliseconds since epoch.
	jobCreateCmd.Flags().StringVarP(&ToTimeOpt, "to", "t", "", "Search window end time (e.g. 2017-07-16T00:00:00)")
	jobCreateCmd.Flags().StringVarP(&TimeZoneOpt, "timezone", "z", "UTC", "Timezone to use for search window")
	jobCreateCmd.Flags().BoolP("by-receipt-time", "b", false, "Use receipt-time instead of log message timestamps")
	jobCreateCmd.Flags().StringVarP(&AutoParsingModeOpt, "auto-parse", "A", "", "Specify auto-parsing mode to use (['performance'] or 'intelligent' - automatically runs field extraction rules)")
}
