//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var replaceFileW = windows.NewLazySystemDLL("kernel32.dll").NewProc("ReplaceFileW")

// replaceOutput uses ReplaceFileW for an existing destination. MoveFileEx is
// only used when there is no destination to replace; if one appears meanwhile,
// retry the atomic replacement path.
func replaceOutput(temporary, target string) error {
	temporaryPath, err := windows.UTF16PtrFromString(temporary)
	if err != nil {
		return err
	}
	targetPath, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return err
	}

	if err := replaceExistingFile(targetPath, temporaryPath); err == nil {
		return nil
	} else if err != windows.ERROR_FILE_NOT_FOUND {
		return err
	}

	err = windows.MoveFileEx(temporaryPath, targetPath, windows.MOVEFILE_WRITE_THROUGH)
	if err != windows.ERROR_ALREADY_EXISTS {
		return err
	}
	return replaceExistingFile(targetPath, temporaryPath)
}

func replaceExistingFile(target, temporary *uint16) error {
	result, _, err := replaceFileW.Call(
		uintptr(unsafe.Pointer(target)),
		uintptr(unsafe.Pointer(temporary)),
		0,
		0,
		0,
		0,
	)
	if result == 0 {
		return err
	}
	return nil
}
