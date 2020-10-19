package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

func TestSingle(test *testing.T) {
	asserter := assert.New(test)
	caps := selenium.Capabilities{
		"bstack:options": map[string]interface{}{
			"os":              "Windows",
			"osVersion":       "10",
			"local":           "false",
			"seleniumVersion": "4.0.0-alpha-6",
			"projectName":     "BrowserStack",
			"buildName":       "Demo-GoLang",
			"sessionName":     "GoLang Firefox Test",
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
		test.Error(titleErr)
	}
	test.Log("Title Received:", title)
	asserter.Contains(title, "Google", "Title should contain google")
}

func TestParallel(test *testing.T) {
	// asserter := assert.New(test)
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
			nestedTest.Cleanup(func() { wd.Quit() })
			nestedTest.Parallel() // adding here to run tests in parallel,
			asserter := assert.New(nestedTest)
			wd.Get("https://google.com")
			title, titleErr := wd.Title()
			if titleErr != nil {
				nestedTest.Error(titleErr)
			}
			nestedTest.Logf("Title Recieved: %s", title)
			asserter.Contains(title, "Google", "Title should contain Google")
		})
	}
}
