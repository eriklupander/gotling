build:
	GO111MODULE=on go build

fmt:
	go fmt

release:
	mkdir -p dist
	GO111MODULE=on GOOS=darwin go build -o dist/gotling-darwin-amd64
	GO111MODULE=on GOOS=linux go build -o dist/gotling-linux-amd64
	GO111MODULE=on GOOS=windows go build -o dist/gotling-windows-amd64