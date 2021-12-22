## Demo-GoLang-BrowserStack <img src="https://camo.githubusercontent.com/799a5c97a4d00394703cf20a5de308784c5454c05726b4c6ba559397644e58d2/68747470733a2f2f643938623874316e6e756c6b352e636c6f756466726f6e742e6e65742f70726f64756374696f6e2f696d616765732f6c61796f75742f6c6f676f2d6865616465722e706e673f31343639303034373830" height="22">

---

[![BrowserStack Status](https://automate.browserstack.com/badge.svg?badge_key=cTY2c2NMN2tPSzRJUUZNbFpXQ1doRGlhRVFrWG5KOXkzbmN1RTFMdjZNbz0tLUd6L1NDRkp4NHlhZ2UwYWphTytQWHc9PQ==--76e8689d99c05a6556cfaf5b48a2759865cfebd3)](https://automate.browserstack.com/public-build/cTY2c2NMN2tPSzRJUUZNbFpXQ1doRGlhRVFrWG5KOXkzbmN1RTFMdjZNbz0tLUd6L1NDRkp4NHlhZ2UwYWphTytQWHc9PQ==--76e8689d99c05a6556cfaf5b48a2759865cfebd3)

This repository contains sample tests to run on the BrowserStack Infrastructure using Selenium and GoLang.

### Setup

- Install the following necessary packages using command line (refer below for performing local testing):

  ```sh
  # install selenium client bindings for go-lang
  go get github.com/tebeka/selenium
  # install asserters
  go get github.com/stretchr/testify

  # OR install all the dependencies using

  go get -v ./...
  go install ./...
  ```

  > <small>Note: This step is entirely optional as go automatically downloads missing dependencies while running or building the programs</small>

- Export the environment variables for the Username and Access Key of your BrowserStack account

  ```sh
  export BROWSERSTACK_USERNAME=<browserstack-username> &&
  export BROWSERSTACK_ACCESS_KEY=<browserstack-access-key>
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
# run local test on desktop and mobile browsers
make local # it will start a local file server, serving the web pages hosted in website folder
# run tests single and parallel tests on mobile browsers
make mobile
# run all tests
make all
# run a failed test
make fail
```
