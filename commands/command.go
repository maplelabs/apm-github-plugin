/*
Package commands contains github-audit shell commands related logic.
*/
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/maplelabs/github-audit/input"
	"github.com/maplelabs/github-audit/internal/configurator"
	"github.com/maplelabs/github-audit/internal/taskmanager"
	"github.com/maplelabs/github-audit/logger"
	"github.com/spf13/cobra"
)

var (
	//ConfigFile stores config.yaml location
	ConfigFile string
	log        logger.Logger
	// GithubAuditBuildInfo contains build and version information for github-audit
	GithubAuditBuildInfo BuildInfo
)

// BuildInfo contains information related to github-audit build and version.
type BuildInfo struct {
	// version of github-audit.
	Version string

	// build time of github-audit.
	BuildTime string

	//commit hash of github-audit.
	Commit string
}

// rootCmd represents the base command which is executed in normal run of github-audit.
var rootCmd = &cobra.Command{
	Use:   "github-audit",
	Short: "github-audit monitors github repositories",
	Long: `github-audit monitors github commits , operations and sends to observability platform.
To run github-audit use subcommand "start" 
Ex: github-audit start --config=<path to config.yaml>

To stop github-audit use subcommand "stop"
Ex: github-audit stop
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// package main ctx received here
		ctx := cmd.Context()
		return StartGitAudit(ctx)
	},
}

// versionCmd represents the version command for printing github-audit build information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of github-audit",
	Long:  `All software has versions. This is github-audit's version with build information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Github-Audit\n")
		fmt.Printf("Version: %s\n", GithubAuditBuildInfo.Version)
		fmt.Printf("BuildTime: %s\n", GithubAuditBuildInfo.BuildTime)
		fmt.Printf("GitCommit: %s\n", GithubAuditBuildInfo.Commit)
		os.Exit(0)
	},
}

// Execute is the starting point to rootCmd.
func Execute(ctx context.Context) {
	// disabled default completion options.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Errorf("error[%v] in root command execution", err)
		os.Exit(1)
	}
}

func init() {
	//initialising config file flag in CLI args , default: config.yaml
	rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "config.yaml", "github-audit start --config=config.yaml")
	//setting config file for logging moodule
	logger.ConfigFile = ConfigFile
	//setting log file
	log = logger.GetLogger()
	// added version command
	rootCmd.AddCommand(versionCmd)
	// over writing help command to exit on call
	origHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		origHelpFunc(cmd, args)
		os.Exit(0)
	})
}

// StartGitAudit starts the github-audit process , creating config struct , starting task manager etc.
func StartGitAudit(ctx context.Context) error {
	log.Infof("starting github-audit process")
	config, err := input.InitConfig(ConfigFile)
	if err != nil {
		log.Errorf("error[%v] in initialising config file to Config struct", err)
		return err
	}
	log.Debugf("parsed config file %#v to config struct", config)
	tasks, err := configurator.StartProcessing(config)
	if err != nil {
		log.Errorf("error[%v] in configuring audit tasks", err)
		return err
	}
	//blocking call returns only if github-audit stops or github-audit crashes
	taskmanager.StartTasks(ctx, tasks)
	return nil
}
