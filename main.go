package main

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	createFile      = kernel32.NewProc("CreateFileW")
	deviceIoControl = kernel32.NewProc("DeviceIoControl")
)

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

	killProcesses(
		2*time.Second,
		// steam
		"steam.exe",
		"steamwebhelper.exe",
		// ea
		"EADesktop.exe",
		"EABackgroundService.exe",
		"EALocalHostSvc.exe",
		// riot/vanguard
		"Riot Client.exe",
		"RiotClientCrashHandler.exe",
		"RiotClientServices.exe",
		"vgc.exe",
		"vgtray.exe",
		// they bought badlion so just to make this future-proof
		"Lunar Client.exe",
		// epic games
		"EpicGamesLauncher.exe",
		"EpicWebHelper.exe",
		// eac
		"EACefSubProcess.exe",
	)

	clearTempDir()

	log.Info("launching " + setupPath + " ...")
	if err = runSetupDetached(setupPath); err != nil {
		log.Error("error starting " + setupPath + ": " + err.Error())
	} else {
		log.Info(setupPath + "  started successfully.")
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
