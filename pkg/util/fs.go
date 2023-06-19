package util

import (
	"fmt"
	"os"
	"path"
)

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(os.Args) >= 2 {
		dirArg := os.Args[1]
		if path.IsAbs(dirArg) {
			dir = dirArg
		} else {
			dir = path.Clean(path.Join(dir, dirArg))
		}
	}

	return dir
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func ReadNumber(path string) (num int) {
	fd, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", path, err))
	}

	defer fd.Close()
	num = -1
	_, err = fmt.Fscanf(fd, "%d", &num)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func WriteNumber(path string, num int) error {
	fd, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", path, err))
	}

	defer fd.Close()
	_, err = fmt.Fprintf(fd, "%d", num)
	return err
}
