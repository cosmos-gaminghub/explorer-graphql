all: get_deps build

get_deps:
	rm -rf ./vendor
	rm -rf ./build
	go mod download
	go mod vendor

build:
	go build -o build/graphql server.go
