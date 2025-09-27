package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func open(url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}

func browse(url string) error {
	// wait a little
	time.Sleep(100 * time.Millisecond)

	_, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(" couldn't open %s", url)
	}
	if err := open(url); err != nil {
		return fmt.Errorf(" couldn't open %s", url)
	}
	return nil
}
