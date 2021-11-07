#!/bin/sh

docker run --rm -v $(pwd):/app -it -p8080:8080 golang "bash"
