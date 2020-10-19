test:
	go test -v ./...

testReport:
	go test -v ./... | go-junit-report > report.xml

single:
	go test -v -run TestSingle

parallel:
	go test -v -run TestParallel