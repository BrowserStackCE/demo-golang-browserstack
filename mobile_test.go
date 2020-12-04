package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

func TestSingleMobile(test *testing.T) {
	test.Parallel()
	asserter := assert.New(test)
	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"osVersion":    "13",
			"deviceName":   "iPhone XS",
			"realMobile":   "true",
			"projectName":  "BrowserStack",
			"buildName":    "Demo-GoLang",
			"sessionName":  "GoLang iPhone XS Test Single",
			"local":        "false",
			"debug":        "true",
			"networkLogs":  "true",
			"consoleLogs":  "verbose",
			"maskCommands": "setValues,getValues, setCookies,getCookies",
		},
		"browserName": "iPhone",
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESSKEY")))
	if err != nil {
		panic(err)
	}
	test.Cleanup(func() { wd.Quit() })
	wd.Get("https://google.com")
	title, titleErr := wd.Title()
	if titleErr != nil {
		test.Fatal(titleErr)
	}
	// test.Log("Title Received:", title)
	asserter.Contains(title, "Google", "Title should contain google")
}

func TestParallelMobile(test *testing.T) {
	// asserter := assert.New(test)
	test.Parallel()
	var capabilities []map[string]interface{}
	fileData, _ := ioutil.ReadFile("./config/devices.json")
	json.Unmarshal(fileData, &capabilities)
	var remoteServerURL = fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESSKEY"))
	for _, capability := range capabilities {
		test.Run(fmt.Sprintf("Running on %s", capability["browserName"]), func(nestedTest *testing.T) {
			// nestedTest.Parallel() // when enabled this it runs all tests in parallel but always run for the last capability
			wd, err := selenium.NewRemote(capability, remoteServerURL)
			if err != nil {
				panic(err)
			}
			nestedTest.Cleanup(func() {
				sessionID := wd.SessionID()
				wd.Quit()
				test.Log(sessionID)
				var req *http.Request
				if test.Failed() {
					req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.browserstack.com/automate/sessions/%s.json", sessionID), strings.NewReader(`{"status":"failed", "reason":"failed all tests"}`))
					if err != nil {
						test.Fatal(err)
					}
				} else {
					req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.browserstack.com/automate/sessions/%s.json", sessionID), strings.NewReader(`{"status":"passed", "reason":"passed all tests"}`))
					if err != nil {
						test.Fatal(err)
					}
				}
				req.SetBasicAuth(os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESSKEY"))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				_, err := client.Do(req)
				if err != nil {
					test.Fatal(err)
				}
			})
			nestedTest.Parallel() // adding here to run tests in parallel,
			asserter := assert.New(nestedTest)
			wd.Get("https://google.com")
			title, titleErr := wd.Title()
			if titleErr != nil {
				nestedTest.Fatal(titleErr)
			}
			// nestedTest.Logf("Title Recieved: %s", title)
			asserter.Contains(title, "Google", "Title should contain Google")
		})
	}
}

func TestLocalMobile(test *testing.T) {
	test.Parallel()
	if os.Getenv("JENKINS_ENV") == "" {
		var bslocalCmd BrowserStackLocal
		err := bslocalCmd.StartLocal() // defined in local.go
		if err != nil {
			test.Fatal(err.Error())
		}
		test.Cleanup(func() {
			bslocalCmd.StopLocal()
		})
		os.Setenv("BROWSERSTACK_LOCAL_IDENTIFIER", "demo")
	}
	// Starting local binary

	fileServer := &http.Server{
		Addr:    ":4000",
		Handler: http.FileServer(http.Dir("./website")),
	}
	go fileServer.ListenAndServe()
	test.Cleanup(func() { fileServer.Close() })

	test.Log("Server started")

	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"osVersion":       "13",
			"deviceName":      "iPhone XS",
			"realMobile":      "true",
			"projectName":     "BrowserStack",
			"buildName":       "Demo-GoLang",
			"sessionName":     "GoLang iPhone XS Test Single",
			"local":           "true",
			"localIdentifier": os.Getenv("BROWSERSTACK_LOCAL_IDENTIFIER"),
		},
		"browserName": "iPhone",
	}
	// time.Sleep(30 * time.Second)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESS_KEY")))
	if err != nil {
		test.Fatal(err)
	}
	test.Cleanup(func() {
		wd.Quit()
	})

	asserter := assert.New(test)
	wd.Get("http://bs-local.com:4000")
	// time.Sleep(5 * time.Second)
	osElement, err := wd.FindElement(selenium.ByCSSSelector, ".os .name")
	if err != nil {
		test.Fatal(err)
	}
	osVal, err := osElement.Text()
	if err != nil {
		test.Fatal(err)
	}
	asserter.Equal(osVal, "iPhone", "OS for the local run should be Windows")
}
