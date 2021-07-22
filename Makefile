all: get_deps build

get_deps:
	rm -rf ./vendor
	go mod download
	go mod vendor

build:
	go build -o build/graphql server.go
