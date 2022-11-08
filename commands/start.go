/*
Package commands contains github-audit shell commands related logic.
*/
package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	// github-audit pid file
	pidFile = "github-audit.process"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start commands start github-audit as a background process",
	Long: `start commands starts github-audit as a background process
To start github-audit with default config.yaml 
Ex: github-audit start
	
To start github-audit with custom config.yaml location
Ex: github-audit start --config=<path to config.yaml>`,
	Run: func(cmd *cobra.Command, args []string) {
		err := start()
		if err != nil {
			fmt.Printf("error[%v] while starting github-audit", err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

// start starts github-audit as a background process.
// it starts github-audit and saves it's process pid in git-audit.process file.
func start() error {
	// checking if pidFile exists , only starting if pidFile does not exists
	_, err := os.Stat(pidFile)
	// pidFile exists , so github-audit already running
	if err == nil {
		log.Errorf("error[github-audit already running]")
		fileByte, err := os.ReadFile(pidFile)
		if err != nil {
			log.Errorf("error[%v] in reading github-audit.process pid file", err)
		}
		return fmt.Errorf("github-audit already running with pid %v , please stop first before starting", string(fileByte))
	}
	log.Info("starting github-audit background process")
	// getting github-audit binary path
	githubAuditPath, err := os.Executable()
	if err != nil {
		log.Errorf("error[%v] in getting github-audit binary path", err)
		return err
	}
	cmd := exec.Command(githubAuditPath, "--config", ConfigFile)
	// starting as deamon process
	err = cmd.Start()
	if err != nil {
		log.Errorf("error[%v] in starting github-audit background process", err)
		return err
	}
	log.Infof("github-audit process started successfully in background with pid %v", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(fmt.Sprintf("%v", cmd.Process.Pid)), 0644)
	if err != nil {
		log.Errorf("error[%v] in writing github-audit.process pid file", err)
		return err
	}
	return nil
}

func init() {
	// adding start sub command to root command.
	rootCmd.AddCommand(startCmd)
}
