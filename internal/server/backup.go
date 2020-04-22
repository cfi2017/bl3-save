package server

import (
	"fmt"
	"io"
	"os"
	"time"
)

func backup(pwd, id string) {
	dumpDir := fmt.Sprintf("%s/backups/%s", pwd, id)
	err := os.MkdirAll(dumpDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	dumpPath := fmt.Sprintf("%s/%s-%v.sav", dumpDir, id, time.Now().Unix())

	_, err = copyFile(fmt.Sprintf("%s/%s.sav", pwd, id), dumpPath)
	if err != nil {
		panic(err)
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s ios not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	return io.Copy(destination, source)
}
