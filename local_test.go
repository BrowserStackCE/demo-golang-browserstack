package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

func TestLocal(test *testing.T) {
	test.Parallel()
	if os.Getenv("JENKINS_ENV") == "" {
		var bslocalCmd BrowserStackLocal
		err := bslocalCmd.StartLocal("demo") // defined in local.go
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

	test.Run("Desktop", func(webTest *testing.T) {
		var buildName = "Demo-GoLang"
		if os.Getenv("JENKINS_ENV") != "" {
			buildName = os.Getenv("BROWSERSTACK_BUILD_NAME")
		}
		caps := selenium.Capabilities{
			"bstack:options": map[string]interface{}{
				"os":              "Windows",
				"osVersion":       "10",
				"seleniumVersion": "4.0.0-alpha-6",
				"projectName":     "BrowserStack GoLang",
				"buildName":       buildName,
				"sessionName":     "GoLang Firefox Test Local",
				"local":           "true",
				"localIdentifier": os.Getenv("BROWSERSTACK_LOCAL_IDENTIFIER"),
			},
			"browserName":    "Firefox",
			"browserVersion": "latest",
		}
		wd, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESS_KEY")))
		if err != nil {
			webTest.Fatal(err)
		}
		webTest.Cleanup(func() {
			wd.Quit()
		})
		webTest.Parallel()

		asserter := assert.New(webTest)

		wd.Get("http://localhost:4000")
		time.Sleep(500 * time.Millisecond)
		osElement, err := wd.FindElement(selenium.ByCSSSelector, ".os .name")
		if err != nil {
			webTest.Fatal(err)
		}

		osVal, err := osElement.Text()
		if err != nil {
			webTest.Fatal(err)
		}

		asserter.Equal(osVal, "Windows", "OS for the local run should be Windows")
	})

	test.Run("Mobile", func(mobileTest *testing.T) {
		var buildName = "Demo-GoLang"
		if os.Getenv("JENKINS_ENV") != "" {
			buildName = os.Getenv("BROWSERSTACK_BUILD_NAME")
		}
		caps := selenium.Capabilities{
			"bstack:options": map[string]interface{}{
				"osVersion":       "13",
				"deviceName":      "iPhone XS",
				"realMobile":      "true",
				"projectName":     "BrowserStack GoLang",
				"buildName":       buildName,
				"sessionName":     "GoLang iPhone XS Test Local",
				"local":           "true",
				"localIdentifier": os.Getenv("BROWSERSTACK_LOCAL_IDENTIFIER"),
			},
			"browserName": "iPhone",
		}
		// time.Sleep(30 * time.Second)
		wd, err := selenium.NewRemote(caps, fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", os.Getenv("BROWSERSTACK_USERNAME"), os.Getenv("BROWSERSTACK_ACCESS_KEY")))
		if err != nil {
			mobileTest.Fatal(err)
		}
		mobileTest.Cleanup(func() {
			wd.Quit()
		})
		mobileTest.Parallel()

		asserter := assert.New(mobileTest)
		wd.Get("http://bs-local.com:4000")
		time.Sleep(500 * time.Millisecond)
		osElement, err := wd.FindElement(selenium.ByCSSSelector, ".os .name")
		if err != nil {
			mobileTest.Fatal(err)
		}
		osVal, err := osElement.Text()
		if err != nil {
			mobileTest.Fatal(err)
		}
		asserter.Equal(osVal, "iPhone", "OS for the local run should be Windows")
	})
}
