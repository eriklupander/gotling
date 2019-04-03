GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
TEST_RESULTS=/tmp/test-results

build:
	mkdir -p bin
	GO111MODULE=on go build -o bin/gotling cmd/gotling/main.go

test:
	mkdir -p ${TEST_RESULTS}
	@go test -coverprofile=${TEST_RESULTS}/unittest.out -v $(GOPACKAGES)
	@go tool cover -html=${TEST_RESULTS}/unittest.out -o ${TEST_RESULTS}/unittest-coverage.html
	rm -f ${TEST_RESULTS}/unittest.out

vet:
	go vet ./...

fmt:
	go fmt ./...

release:
	mkdir -p dist
	GO111MODULE=on GOOS=darwin go build -o dist/gotling-darwin-amd64 cmd/gotling/main.go
	GO111MODULE=on GOOS=linux go build -o dist/gotling-linux-amd64 cmd/gotling/main.go
	GO111MODULE=on GOOS=windows go build -o dist/gotling-windows-amd64 cmd/gotling/main.go