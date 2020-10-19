# Demo-GoLang

---

This repository contains sample tests to run on the BrowserStack Infrastructure using Selenium and GoLang.

### Setup

Install the following packages in using command line:

```sh
# install selenium go-lang
go get github.com/tebeka/selenium
# install asserters
go get github.com/stretchr/testify
```

> NOTE: If you are not using \*nix based systems you may also need to install `make` command to run the commands given below from [here](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows). You can also refer to the [make file](Makefile) and directly copy test commands for eg command to run single web tests `go test -v -run TestSingle`

### Web

To run tests on a website run anyone of the following commands:

```sh
# run single test
make single
# run multiple tests in parallel
make parallel
# run all tests
make test
```
