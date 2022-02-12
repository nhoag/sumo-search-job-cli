/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nhoag/sumo-search-job-cli/client"
	"github.com/spf13/cobra"
)

// jobDeleteCmd represents the jobDelete command
var jobDeleteCmd = &cobra.Command{
	Use:   "jobDelete JOB_ID",
	Short: "Delete a Sumo Logic Search Job",
	Long: `The jobDelete command will delete a Sumo Logic Search Job via the
	Search Job API.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		QuietOpt, _ = cmd.Flags().GetBool("quiet")
		VerboseOpt, _ = cmd.Flags().GetBool("verbose")
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tSTART\tjobDelete\n", time.Now().UnixNano())
		}
		executeDelete(cmd, args)
		if VerboseOpt {
			fmt.Fprintf(os.Stderr, "%d\tEND\tjobDelete\n", time.Now().UnixNano())
		}
	},
}

func validateDelete() {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobDelete::validateDelete()\n", time.Now().UnixNano())
	}
	// Add validation logic here
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobDelete::validateDelete()\n", time.Now().UnixNano())
	}
}

func executeDelete(cmd *cobra.Command, args []string) {
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tSTART\tjobDelete::executeDelete()\n", time.Now().UnixNano())
	}
	client.DeleteSearchJob(args[0])
	if !QuietOpt {
		fmt.Fprintf(os.Stderr, "Successfully Deleted Search Job!\n")
	}
	if VerboseOpt {
		fmt.Fprintf(os.Stderr, "%d\tEND\tjobDelete::executeDelete()\n", time.Now().UnixNano())
	}
}

func init() {
	rootCmd.AddCommand(jobDeleteCmd)
}
