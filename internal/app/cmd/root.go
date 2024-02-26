package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string

	debug bool

	startTime = time.Now()
)

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notion-cli",
		Short: "notion is a command line tool for interacting with the Notion API",

		// Needs to exist to make --version work
		Run: func(cmd *cobra.Command, args []string) {},

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", Version, Commit, Date)
				os.Exit(0)
			}
			if debug, _ = cmd.Flags().GetBool("debug"); debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
	}

	cmd.PersistentFlags().Bool("version", false, "Print the version of the CLI")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug logging")

	return cmd
}

func Execute() {
	logrus.SetOutput(os.Stderr)

	rootCmd := getRootCmd()
	rootCmd.AddCommand(getDBIssueCmd())
	rootCmd.AddCommand(getDBIssueDetailCmd())
	rootCmd.AddCommand(getMaintenanceCmd())

	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	logrus.Debugf("Execution took: %s", time.Since(startTime))
}
