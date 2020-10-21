## Demo-GoLang

---

This repository contains sample tests to run on the BrowserStack Infrastructure using Selenium and GoLang.

### Setup

Install the following necessary packages using command line (refer below for performing local testing):

```sh
# install selenium client bindings for go-lang
go get github.com/tebeka/selenium
# install asserters
go get github.com/stretchr/testify
```

or install all the packages using

```sh
go install .
```

> NOTE: If you are not using \*nix based systems you may also need to install `make` command to run the commands given below from [here](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows). You can also refer to the [make file](Makefile) and directly copy test commands for eg command to run single web tests `go test -v -run TestSingle`

#### Local Support

<small> You will need to download the local binary from [here](https://www.browserstack.com/local-testing/releases), install it and ensure it is the PATH on your machine as there is no local binding specifically for golang. </small>

### Web

To run tests on a website run anyone of the following commands:

```sh
# run single test
make single
# run multiple tests in parallel
make parallel
# run local test
make local # it will start a local file server, serving the web pages hosted in website folder
# run all tests
make test
# run a test and mark it's status as pass or fail
make test-status
# run all tests and generate tests reports in junit-reporter
make testReport
```
