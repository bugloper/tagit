.PHONY: test build

test:
	go test ./test/... -v

build:
	GOOS=darwin GOARCH=amd64 go build -o tagit-amd64
	GOOS=darwin GOARCH=arm64 go build -o tagit-arm64

tar:
	tar -czvf tagit-amd64-1.0.2.tar.gz tagit-amd64
	tar -czvf tagit-arm64-1.0.2.tar.gz tagit-arm64
