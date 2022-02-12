/*
Copyright © 2022 Nathaniel Hoag <info@nathanielhoag.com>

*/
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
	Job     string
	JobFile string

	Query     string
	QueryFile string

	FromTime string
	Span     string
	ToTime   string
	TimeZone string

	AutoParsingMode string
)

// jobCreateCmd represents the jobCreate command
var jobCreateCmd = &cobra.Command{
	Use:   "jobCreate",
	Short: "Create a Sumo Logic Search Job",
	Long: `The jobCreate command will initiate a Sumo Logic Search Job via the
	Search Job API. The returned Job ID can be used to query job status and
	fetch search results.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateJobCreate()
		location, jobId := executeSearchJob(buildPayload(cmd, args))
		fmt.Fprintf(os.Stderr, "Location:\t%s\nJob ID:\t\t%s\n", location, jobId)
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
	var jobDef JobDefinition
	jobDef.Timezone = "UTC"
	if len(Job) > 0 {
		json.Unmarshal([]byte(Job), &jobDef)
	}
	if len(JobFile) > 0 {
		content, err := ioutil.ReadFile(JobFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		json.Unmarshal([]byte(content), &jobDef)
	}
	if len(Query) > 0 {
		jobDef.Query = Query
	}
	if len(QueryFile) > 0 {
		content, err := ioutil.ReadFile(QueryFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		jobDef.Query = string(content)
	}
	if len(Span) > 0 {
		duration, _ := time.ParseDuration(Span)
		jobDef.To = time.Now().UTC().String()
		jobDef.From = time.Now().Add(-duration).UTC().String()
	}
	if len(ToTime) > 0 {
		jobDef.To = ToTime
	}
	if len(FromTime) > 0 {
		jobDef.From = FromTime
	}
	if len(TimeZone) > 0 {
		jobDef.Timezone = TimeZone
	}
	if len(Job) == 0 && len(JobFile) == 0 {
		jobDef.ByReceiptTime, _ = cmd.Flags().GetBool("by-receipt-time")
	}
	if len(AutoParsingMode) > 0 {
		jobDef.AutoParsingMode = AutoParsingMode
	}

	return jobDef
}

func executeSearchJob(jobDef JobDefinition) (*url.URL, string) {
	searchJobDef := *openapi.NewSearchJobDefinition()
	searchJobDef.SetTo(jobDef.To)
	searchJobDef.SetFrom(jobDef.From)
	searchJobDef.SetTimeZone(jobDef.Timezone)
	searchJobDef.SetQuery(jobDef.Query)

	return client.CreateSearchJob(searchJobDef)
}

func validateJobCreate() {
	if len(Job) > 0 {
		if len(JobFile) > 0 {
			fmt.Println("job-file is not compatible with job")
			os.Exit(1)
		}
		if len(Query) > 0 {
			fmt.Println("query is not compatible with job")
			os.Exit(1)
		}
		if len(QueryFile) > 0 {
			fmt.Println("query-file is not compatible with job")
			os.Exit(1)
		}
		if len(FromTime) > 0 {
			fmt.Println("from is not compatible with job")
			os.Exit(1)
		}
		if len(Span) > 0 {
			fmt.Println("span is not compatible with job")
			os.Exit(1)
		}
		if len(ToTime) > 0 {
			fmt.Println("to is not compatible with job")
			os.Exit(1)
		}
		if len(AutoParsingMode) > 0 {
			fmt.Println("auto-parse is not compatible with job")
			os.Exit(1)
		}
	}
	if len(JobFile) > 0 {
		// @todo: check file is readable
		if len(Job) > 0 {
			fmt.Println("job is not compatible with job-file")
			os.Exit(1)
		}
		if len(Query) > 0 {
			fmt.Println("query is not compatible with job")
			os.Exit(1)
		}
		if len(QueryFile) > 0 {
			fmt.Println("query-file is not compatible with job")
			os.Exit(1)
		}
		if len(FromTime) > 0 {
			fmt.Println("from is not compatible with job")
			os.Exit(1)
		}
		if len(Span) > 0 {
			fmt.Println("span is not compatible with job")
			os.Exit(1)
		}
		if len(ToTime) > 0 {
			fmt.Println("to is not compatible with job")
			os.Exit(1)
		}
		if len(AutoParsingMode) > 0 {
			fmt.Println("auto-parse is not compatible with job")
			os.Exit(1)
		}
	}
	if len(Query) > 0 {
		if len(QueryFile) > 0 {
			fmt.Println("query-file is not compatible with query")
			os.Exit(1)
		}
	}
	if len(QueryFile) > 0 {
		// @todo: check file is readable
		if len(Query) > 0 {
			fmt.Println("query is not compatible with query-file")
			os.Exit(1)
		}
		fileInfo, err := os.Stat(QueryFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		mode := fileInfo.Mode()
		fmt.Println(mode.Perm())
		os.Exit(1)
	}
	if len(Span) > 0 {
		if len(FromTime) > 0 {
			fmt.Println("from is not compatible with span")
			os.Exit(1)
		}
		if len(ToTime) > 0 {
			fmt.Println("to is not compatible with span")
			os.Exit(1)
		}
		if len(TimeZone) > 0 {
			fmt.Println("timezone is not compatible with span")
			os.Exit(1)
		}
	}

	var toTime time.Time = time.Now().UTC()
	var fromTime time.Time = time.Now().UTC()
	var err error
	var duration time.Duration

	if len(Span) > 0 {
		duration, err = time.ParseDuration(Span)
		if err != nil {
			fmt.Println("Unable to parse the provided span: " + Span)
			os.Exit(1)
		}
		fromTime = time.Now().Add(-duration).UTC()
	}
	if len(ToTime) > 0 {
		toTime, err = time.Parse("2006-01-02T15:04:05", ToTime)
		if err != nil {
			fmt.Println("Unable to parse the provided to-time: " + ToTime)
			os.Exit(1)
		}
	}
	if len(FromTime) > 0 {
		fromTime, err = time.Parse("2006-01-02T15:04:05", FromTime)
		if err != nil {
			fmt.Println("Unable to parse the provided from-time: " + FromTime)
			os.Exit(1)
		}
	}
	if len(JobFile) > 0 {
		var jobDef JobDefinition
		content, err := ioutil.ReadFile(JobFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		json.Unmarshal(content, &jobDef)
		toTime, err = time.Parse("2006-01-02T15:04:05", jobDef.To)
		if err != nil {
			fmt.Println("Unable to parse the provided 'to' time: " + jobDef.To)
			os.Exit(1)
		}
		fromTime, err = time.Parse("2006-01-02T15:04:05", jobDef.From)
		if err != nil {
			fmt.Println("Unable to parse the provided 'from' time: " + jobDef.From)
			os.Exit(1)
		}
	}
	if !fromTime.Before(toTime) {
		fmt.Println("from " + fromTime.String() + " is not before to " + toTime.String())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(jobCreateCmd)

	jobCreateCmd.Flags().StringVarP(&Job, "job", "j", "", "Search job definition")
	jobCreateCmd.Flags().StringVarP(&JobFile, "job-file", "J", "", "Path to file with full search job definition")

	jobCreateCmd.Flags().StringVarP(&Query, "query", "q", "", "Search query")
	jobCreateCmd.Flags().StringVarP(&QueryFile, "query-file", "Q", "", "Path to file with search query")

	// - relative lookback (-3h)
	jobCreateCmd.Flags().StringVarP(&Span, "span", "s", "", "Size of time span relative to now")
	// - from time: 2017-07-16T00:00:00
	// @todo: or milliseconds since epoch.
	jobCreateCmd.Flags().StringVarP(&FromTime, "from", "f", "", "Search window start time")
	// - to time: 2017-07-16T00:00:00
	// @todo: or milliseconds since epoch.
	jobCreateCmd.Flags().StringVarP(&ToTime, "to", "t", "", "Search window end time")
	// - timezone if ^^^ not millis
	jobCreateCmd.Flags().StringVarP(&TimeZone, "timezone", "z", "UTC", "Timezone to use for search window")
	jobCreateCmd.Flags().BoolP("by-receipt-time", "b", false, "Use receipt-time instead of log message timestamps")
	// - autoParsingMode (default is 'perfomance', 'intelligent' automatically runs field extraction rules)
	jobCreateCmd.Flags().StringVarP(&AutoParsingMode, "auto-parse", "A", "", "Specify auto-parsing mode to use")
}
