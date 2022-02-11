/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

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
		executeDelete(cmd, args)
		fmt.Fprintf(os.Stderr, "Successfully Deleted Job!\n")
	},
}

func validateDelete() {
	// Add validation logic here
}

func executeDelete(cmd *cobra.Command, args []string) {
	client.DeleteSearchJob(args[0])
}

func init() {
	rootCmd.AddCommand(jobDeleteCmd)
}
