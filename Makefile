build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -o bin/kubevis-darwin-arm64 main.go

build-linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -o bin/kubevis-linux-arm64 main.go

build: build-darwin-arm64 build-linux-arm64
