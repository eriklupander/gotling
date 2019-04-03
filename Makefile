build:
	GO111MODULE=on go build cmd/gotling/main.go

fmt:
	go fmt ./...

release:
	mkdir -p dist
	GO111MODULE=on GOOS=darwin go build -o dist/gotling-darwin-amd64 cmd/gotling/main.go
	GO111MODULE=on GOOS=linux go build -o dist/gotling-linux-amd64 cmd/gotling/main.go
	GO111MODULE=on GOOS=windows go build -o dist/gotling-windows-amd64 cmd/gotling/main.go