export tag=v2.0
root:
	export ROOT=github.com/cncamp/golang

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t csmiao/httpserver:${tag} .

push: release
	echo "pushing csmiao/httpserver"
	docker push csmiao/httpserver:${tag}
