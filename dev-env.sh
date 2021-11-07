#!/bin/sh

docker run --rm -v $(pwd):/app -v $(echo "$HOME/.kube/config"):/root/.kube/config -it -p8080:8080 golang "bash"
