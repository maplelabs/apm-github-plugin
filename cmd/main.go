/*
Package Main is the starting point of github-audit.
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/maplelabs/github-audit/commands"
	"github.com/maplelabs/github-audit/logger"
)

var (
	// log is used for adding logs
	log logger.Logger
)

func init() {
	// initialising the logger
	log = logger.GetLogger()
}

func main() {
	var wg sync.WaitGroup
	// this context is used throughout github-audit lifecycle
	mainCtx, mainCtxCancel := context.WithCancel(context.Background())

	// this logic handles graceful shutdown for github-audit
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sigInt := <-c
			log.Infof("stopping the github-audit as received signal %s", sigInt)
			mainCtxCancel()
			return
		}
	}()
	// this is the starting point of github-audit.
	commands.Execute(mainCtx)
	wg.Wait()
}
