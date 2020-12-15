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
	var buildName = "Demo-GoLang"
	if os.Getenv("JENKINS_ENV") != "" {
		buildName = os.Getenv("BROWSERSTACK_BUILD_NAME")
	}
	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"osVersion":    "13",
			"deviceName":   "iPhone XS",
			"realMobile":   "true",
			"projectName":  "BrowserStack GoLang",
			"buildName":    buildName,
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
			if os.Getenv("JENKINS_ENV") != "" {
				var options map[string]interface{}
				tempOptions, _ := capability["bstack:options"]
				options = tempOptions.(map[string]interface{})
				options["buildName"] = os.Getenv("BROWSERSTACK_BUILD_NAME")
				capability["bstack:options"] = options
			}
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
