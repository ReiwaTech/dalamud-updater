package util

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetDLLVersion(filePath string) (string, error) {
	var zero windows.Handle
	size, err := windows.GetFileVersionInfoSize(filePath, &zero)

	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", err
	}

	if size == 0 {
		return "", fmt.Errorf("size is 0")
	}

	versionInfo := make([]byte, size)
	if err = windows.GetFileVersionInfo(filePath, 0, size, unsafe.Pointer(&versionInfo[0])); err != nil {
		return "", err
	}

	var fixedInfo *windows.VS_FIXEDFILEINFO
	fixedInfoLen := uint32(unsafe.Sizeof(*fixedInfo))
	err = windows.VerQueryValue(unsafe.Pointer(&versionInfo[0]), `\`, (unsafe.Pointer)(&fixedInfo), &fixedInfoLen)
	if err != nil {
		return "", err
	}

	major := fixedInfo.FileVersionMS >> 16
	minor := fixedInfo.FileVersionMS & 0xFFFF
	patch := fixedInfo.FileVersionLS >> 16
	build := fixedInfo.FileVersionLS & 0xFFFF

	return fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build), nil
}
