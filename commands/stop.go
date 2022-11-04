/*
Package commands contains github-audit shell commands related logic.
*/
package commands

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop commands stops github-audit process",
	Long: `stop commands stops github-audit process
To stop github-audit
Ex: github-audit stop`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := stop()
		return err
	},
}

// stop stops github-audit background process.
// it starts github-audit and saves it's process pid in git-audit.process file.
func stop() error {
	fileByte, err := os.ReadFile("git-audit.process")
	if err != nil {
		log.Errorf("error[%v] in reading github-audit.process pid file", err)
		return err
	}
	log.Debugf("github-audit pid as read from github.process is %v", string(fileByte))
	log.Infof("stopping github-audit process")
	// stopping github-audit using linux kill command
	cmd := exec.Command("kill", "-9", string(fileByte))
	err = cmd.Run()
	if err != nil {
		log.Errorf("error[%v] in stopping github-audit process", err)
		return err
	}
	return nil
}

func init() {
	// adding stop sub command to root command.
	rootCmd.AddCommand(stopCmd)
}
