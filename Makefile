default: all

all:
	go test -v ./...

single:
	go test -v -run TestSingle

parallel:
	go test -v -run TestParallel

local:
	go test -v -run TestLocal