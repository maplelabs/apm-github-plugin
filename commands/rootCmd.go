/*
Package commands contains github-audit shell commands related logic.
*/
package commands

import (
	"context"
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
)

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

// Execute is the starting point to rootCmd.
func Execute(ctx context.Context) {
	// disabled default completion options.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Errorf("error[%v] in root command execution", err)
		os.Exit(0)
	}
}

func init() {
	//initialising config file flag in CLI args , default: config.yaml
	rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "config.yaml", "github-audit start --config=config.yaml")
	//setting config file for logging moodule
	logger.ConfigFile = ConfigFile
	//setting log file
	log = logger.GetLogger()
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
