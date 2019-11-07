package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func GoFmtPath(path string) (err error) {
	fmt.Println("fmt path " + path)
	cmd := exec.Command("gofmt", "-w", path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func CurrentPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	return
}
