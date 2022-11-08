/*
	Package commands contains github-audit shell commands related logic.
*/
package commands

import (
	"fmt"
	"os"
	"strconv"

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
	Run: func(cmd *cobra.Command, args []string) {
		err := stop()
		if err != nil {
			fmt.Printf("error[%v] while stopping github-audit", err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

// stop stops github-audit background process.
// it starts github-audit and saves it's process pid in git-audit.process file.
func stop() error {
	var pid int
	fileByte, err := os.ReadFile(pidFile)
	if err != nil {
		log.Errorf("error[%v] in reading github-audit.process pid file", err)
		return err
	}
	log.Debugf("github-audit pid as read from github.process is %v", string(fileByte))
	log.Infof("stopping github-audit process")
	pid, err = strconv.Atoi(string(fileByte))
	if err != nil {
		log.Errorf("error[%v] in converting pid string to integer")
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Errorf("error[%v] in finding github-audit process with pid %v", err, pid)
		return err
	}
	err = proc.Kill()
	if err != nil {
		log.Errorf("error[%v] in killing github-audit process with pid %v", err, pid)
		return err
	}
	// removing pidFile after stopping github-audit
	err = os.Remove(pidFile)
	if err != nil {
		log.Errorf("error[%v] in deleting github-audit pid file %v", err, pidFile)
	}
	return err
}

func init() {
	// adding stop sub command to root command.
	rootCmd.AddCommand(stopCmd)
}
