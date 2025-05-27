package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
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

func ejectDrive(driveLetter string) error {
	drivePath := `\\.\` + driveLetter

	ptr, err := syscall.UTF16PtrFromString(drivePath)
	if err != nil {
		return fmt.Errorf("failed to convert drive path to UTF16: %v", err)
	}

	handle, _, err := createFile.Call(
		uintptr(unsafe.Pointer(ptr)),
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
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("error getting executable path:", err)
		return
	}
	dir := filepath.Dir(exePath)
	setupPath := filepath.Join(dir, "setup.exe")

	fmt.Println("launching setup.exe:", setupPath)
	if err := runSetupDetached(setupPath); err != nil {
		fmt.Println("error starting setup.exe:", err)
	} else {
		fmt.Println("setup.exe started successfully.")
	}

	time.Sleep(3 * time.Second)

	drive := strings.ToUpper(string(exePath[0])) + ":"
	fmt.Println("ejecting usb drive:", drive)
	if err := ejectDrive(drive); err != nil {
		fmt.Println("error ejecting drive:", err)
	} else {
		fmt.Println("drive ejected successfully.")
	}
}
