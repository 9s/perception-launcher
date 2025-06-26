package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

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
