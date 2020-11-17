package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// BrowserStackLocal to start stop BS Local
type BrowserStackLocal exec.Cmd

// StartLocal is to start the BS Local from binary
func (bslocal *BrowserStackLocal) StartLocal() error {
	bslocalCmd := exec.Command("BrowserStackLocal", "--key", os.Getenv("BROWSERSTACK_ACCESSKEY"), "--local-identifier", "demo")
	// bslocalCmd := exec.Command("BrowserStackLocal", "--key", os.Getenv("BROWSERSTACK_ACCESSKEY"))
	bslocalCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	pr, pw := io.Pipe()
	bslocalCmd.Stdout = pw
	bslocalCmd.Stderr = pw

	err := bslocalCmd.Start()
	if err != nil {
		return err
	}
	var str, inp string
	for {
		fmt.Fscanln(pr, &inp)
		str += inp
		// fmt.Println(str)
		// fmt.Println(str, strings.Contains(str, "ERROR"), strings.Contains(str, "SUCCESS"))
		if strings.Contains(str, "ERROR") {
			// fmt.Println("Encountered error, process should be killed automatically")
			syscall.Kill(-bslocalCmd.Process.Pid, syscall.SIGKILL)
			return errors.New("Couldn't start BrowserStack Local. Some error has occured")
		} else if strings.Contains(str, "SUCCESS") {
			// time.Sleep(10 * time.Nanosecond)
			// fmt.Println("Process connected")
			break
		}
	}
	// return bslocalCmd
	*bslocal = BrowserStackLocal(*bslocalCmd)
	return nil
}

// StopLocal to stop local
func (bslocal *BrowserStackLocal) StopLocal() error {
	if bslocal == nil {
		return errors.New("BrowserStack Local is not started. Stop is illogical")
	}
	return syscall.Kill(-bslocal.Process.Pid, syscall.SIGKILL)
}

func main() {
	var bslocalCmd BrowserStackLocal
	err := bslocalCmd.StartLocal()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Minute)
	err = bslocalCmd.StopLocal()
	if err != nil {
		panic(err)
	}
}
