package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/charmbracelet/log"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	createFile      = kernel32.NewProc("CreateFileW")
	deviceIoControl = kernel32.NewProc("DeviceIoControl")
)

func runSetupDetached(setupPath string) error {
	cmd := exec.Command(setupPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

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

func clearTempDir() {
	tempDir := os.TempDir()
	log.Info("clearing temp folder: " + tempDir)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn("error accessing: " + path + " - " + err.Error())
			return nil
		}
		if path == tempDir {
			return nil
		}
		if info.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}
		if err != nil {
			log.Warn("failed to delete: " + path + " - " + err.Error())
		}
		return nil
	})
	if err != nil {
		log.Error("failed to clear temp folder: " + err.Error())
	} else {
		log.Info("temp folder cleared.")
	}
}

func ejectDrive(driveLetter string) error {
	drivePath := `\\.\` + driveLetter
	drivePtr, err := syscall.UTF16PtrFromString(drivePath)
	if err != nil {
		return fmt.Errorf("failed to encode drive path: %v", err)
	}

	handle, _, err := createFile.Call(
		uintptr(unsafe.Pointer(drivePtr)),
		uintptr(syscall.GENERIC_READ|syscall.GENERIC_WRITE),
		uintptr(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE),
		0,
		uintptr(syscall.OPEN_EXISTING),
		0,
		0,
	)
	if handle == 0 {
		return fmt.Errorf("failed to open drive handle: %v", err)
	}
	var bytesReturned uint32
	ret, _, err := deviceIoControl.Call(
		handle,
		0x2D4808,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&bytesReturned)),
		0,
	)
	if ret == 0 {
		return fmt.Errorf("eject failed: %v", err)
	}
	return nil
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetTimeFormat(time.Kitchen)
	log.SetReportCaller(false)
	log.SetPrefix("[perception.cx]")

	exePath, err := os.Executable()
	if err != nil {
		log.Error("error getting executable path: " + err.Error())
		return
	}
	dir := filepath.Dir(exePath)
	setupPath := filepath.Join(dir, "Spotify.exe")
	drive := strings.ToUpper(string(exePath[0])) + ":"

	killProcesses(2*time.Second, "steam.exe", "EADesktop.exe")

	clearTempDir()

	log.Info("launching setup.exe...")
	if err = runSetupDetached(setupPath); err != nil {
		log.Error("error starting setup.exe: " + err.Error())
	} else {
		log.Info("setup.exe started successfully.")
	}

	log.Info("waiting before usb eject...")
	time.Sleep(3 * time.Second)

	log.Info("ejecting usb drive: " + drive)
	if err = ejectDrive(drive); err != nil {
		log.Error("error ejecting drive: " + err.Error())
	} else {
		log.Info("drive ejected successfully.")
	}
}
