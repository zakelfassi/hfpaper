VERSION=1.0.0
BINARY=hfpaper

build:
	go build -o $(BINARY) main.go

install:
	go install

test:
	go test ./...
	go vet ./...

release:
	# Simulating release with local cross-compile for basic targets
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY)-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe main.go

clean:
	rm -f $(BINARY)
	rm -rf dist
