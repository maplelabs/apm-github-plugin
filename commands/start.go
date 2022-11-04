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

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start commands start github-audit as a background process",
	Long: `start commands starts github-audit as a background process
To start github-audit with default config.yaml 
Ex: github-audit start
	
To start github-audit with custom config.yaml location
Ex: github-audit start --config=<path to config.yaml>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := start()
		return err
	},
}

// start starts github-audit as a background process.
// it starts github-audit and saves it's process pid in git-audit.process file.
func start() error {
	log.Info("starting github background process")
	cmd := exec.Command("./github-audit", "--config", ConfigFile)
	// starting as deamon process
	err := cmd.Start()
	if err != nil {
		log.Errorf("error[%v] in starting github-audit background process", err)
		return err
	}
	log.Infof("github-audit process started successfully in background with pid %v", cmd.Process.Pid)
	err = os.WriteFile("git-audit.process", []byte(fmt.Sprintf("%v", cmd.Process.Pid)), 0644)
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
