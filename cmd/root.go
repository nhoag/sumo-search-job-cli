package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	DeploymentOpt string
	QuietOpt      bool
	VerboseOpt    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sumo-search-job-cli",
	Short: "Sumo Logic Search Job CLI",
	Long:  `Command line interface to the Sumo Logic Search Job API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sumo-search-job-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&DeploymentOpt, "deployment", "us1", "Deployment of Sumo Logic instance (au, ca, de, eu, fed, in, jp, us1, us2)")
	viper.BindPFlag("deployment", rootCmd.PersistentFlags().Lookup("deployment"))
	rootCmd.PersistentFlags().BoolP("quiet", "S", false, "Don't display status updates")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose information")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".test-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sumo-search-job-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// React to config file read success here
	}
}
