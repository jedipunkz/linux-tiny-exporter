.PHONY: build build-linux-amd64 build-linux-arm64 build-darwin-arm64 build-win-amd64 build-win-arm64 clean

build: build-linux-amd64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-tiny-exporter-linux-amd64 main.go

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/linux-tiny-exporter-linux-arm64 main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/linux-tiny-exporter-darwin-arm64 main.go

build-win-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/linux-tiny-exporter-win-amd64.exe main.go

build-win-arm64:
	GOOS=windows GOARCH=arm64 go build -o bin/linux-tiny-exporter-win-arm64.exe main.go

clean:
	rm -rf ./bin
