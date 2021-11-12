run:
	go run main.go -out-cluster

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -o bin/kubevis-darwin-arm64 main.go

build-linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -o bin/kubevis-linux-arm64 main.go


build-docker:
	docker build -t kubevis .

build: build-darwin-arm64 build-linux-arm64

kind-load-image: build-docker
	kind load docker-image kubevis

k8s-deploy:
	kubectl apply -f k8s

k8s-undeploy:
	kubectl delete -f k8s
