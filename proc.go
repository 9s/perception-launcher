package main

import (
	"context"
	"github.com/charmbracelet/log"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func killProcesses(timeout time.Duration, names ...string) {
	for _, name := range names {
		log.Info("terminating " + strings.ToLower(name) + "...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		cmd := exec.CommandContext(ctx, "taskkill", "/f", "/im", name)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Run()
		cancel()

		if err != nil {
			log.Warn(strings.ToLower(name) + " not running or failed to terminate: " + err.Error())
		} else {
			log.Info(strings.ToLower(name) + " terminated.")
		}
	}
}

func runSetupDetached(setupPath string) error {
	cmd := exec.Command(setupPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	return cmd.Start()
}
