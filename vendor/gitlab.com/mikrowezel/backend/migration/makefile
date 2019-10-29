# Vars
STG_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=mwgranica

# Misc
BINARY_NAME=granica
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	go  build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

build-linux:
	CGOENABLED=0 GOOS=linux GOARCH=amd64; go build -o ./bin/$(BINARY_UNIX) ./cmd/$(BINARY_NAME).go

test:
	go test -v -count=1 ./...

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

## Misc
custom-build:
	make mod tidy; go mod vendor; go build ./...

get-deps:
	go get "github.com/jmoiron/sqlx v1.2.0"
	go get "github.com/lib/pq v1.2.0"
	go get "gitlab.com/mikrowezel/config"
	go get "google.golang.org/appengine"
