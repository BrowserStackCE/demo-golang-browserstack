package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// StartLocal is to start the BS Local from binary
func StartLocal() *exec.Cmd {
	bslocalCmd := exec.Command("BrowserStackLocal", "--key", os.Getenv("BROWSERSTACK_ACCESSKEY"), "--disable-dashboard")
	bslocalCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	pr, pw := io.Pipe()
	bslocalCmd.Stdout = pw
	bslocalCmd.Stderr = pw

	err := bslocalCmd.Start()
	if err != nil {
		panic(err)
	}
	var str, inp string
	for {
		fmt.Fscanln(pr, &inp)
		str += inp
		// fmt.Println(str, strings.Contains(str, "ERROR"), strings.Contains(str, "SUCCESS"))
		if strings.Contains(str, "ERROR") {
			// fmt.Println("Encountered error, process should be killed automatically")
			syscall.Kill(-bslocalCmd.Process.Pid, syscall.SIGKILL)
			return nil
		} else if strings.Contains(str, "SUCCESS") {
			// fmt.Println("Process connected")
			break
		}
	}
	return bslocalCmd
}

func main() {
	bslocalCmd := StartLocal()
	time.Sleep(1 * time.Minute)
	syscall.Kill(-bslocalCmd.Process.Pid, syscall.SIGKILL)
}
