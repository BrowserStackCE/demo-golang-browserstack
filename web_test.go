package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

func TestSingle(test *testing.T) {
	test.Parallel()
	asserter := assert.New(test)
	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"os":              "Windows",
			"osVersion":       "10",
			"local":           "false",
			"seleniumVersion": "4.0.0-alpha-6",
			"projectName":     "BrowserStack GoLang",
			"buildName":       "Demo-GoLang",
			"sessionName":     "GoLang Firefox Test Single",
			"debug":           "true",
			"networkLogs":     "true",
			"consoleLogs":     "verbose",
		},
		"browserName":    "Firefox",
		"browserVersion": "latest",
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
	asserter.Contains(title, "Google", "Title should contain google")
}

func TestParallel(test *testing.T) {
	test.Parallel()
	var capabilities []map[string]interface{}
	fileData, _ := ioutil.ReadFile("./config/browsers.json")
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
				nestedTest.Log(sessionID)
				var req *http.Request
				if nestedTest.Failed() {
					req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.browserstack.com/automate/sessions/%s.json", sessionID), strings.NewReader(`{"status":"failed", "reason":"failed all tests"}`))
					if err != nil {
						nestedTest.Fatal(err)
					}
				} else {
					req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.browserstack.com/automate/sessions/%s.json", sessionID), strings.NewReader(`{"status":"passed", "reason":"passed all tests"}`))
					if err != nil {
						nestedTest.Fatal(err)
					}
				}
				req.SetBasicAuth(os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESSKEY"))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				_, err := client.Do(req)
				if err != nil {
					nestedTest.Fatal(err)
				}
			})
			nestedTest.Parallel() // adding here to run tests in parallel,
			asserter := assert.New(nestedTest)
			wd.Get("https://google.com")
			title, titleErr := wd.Title()
			if titleErr != nil {
				nestedTest.Fatal(titleErr)
			}
			nestedTest.Logf("Title Recieved: %s", title)
			asserter.Contains(title, "Google", "Title should contain Google")
		})
	}
}

func TestFail(test *testing.T) {
	if os.Getenv("FAIL_TEST") == "" {
		test.SkipNow()
	}
	test.Parallel()
	asserter := assert.New(test)
	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"os":              "Windows",
			"osVersion":       "10",
			"local":           "false",
			"seleniumVersion": "4.0.0-alpha-6",
			"projectName":     "BrowserStack GoLang",
			"buildName":       "Demo-GoLang",
			"sessionName":     "GoLang Firefox Test Fail",
			"debug":           "true",
			"networkLogs":     "true",
			"consoleLogs":     "verbose",
		},
		"browserName":    "Firefox",
		"browserVersion": "latest",
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESSKEY")))
	if err != nil {
		panic(err)
	}
	test.Cleanup(func() {
		if test.Failed() {
			wd.ExecuteScript("browserstack_executor: {\"action\": \"setSessionStatus\", \"arguments\": {\"status\":\"failed\"}}", nil)
		} else {
			wd.ExecuteScript("browserstack_executor: {\"action\": \"setSessionStatus\", \"arguments\": {\"status\":\"passed\"}}", nil)
		}
		wd.Quit()
	})
	wd.Get("https://google.com")
	title, titleErr := wd.Title()
	if titleErr != nil {
		test.Fatal(titleErr)
	}
	asserter.Equal("Microsoft", title, "Title should have been Google")
}
