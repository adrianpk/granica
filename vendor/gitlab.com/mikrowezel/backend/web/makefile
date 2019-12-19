# Vars
STG_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=mw_web

# Misc
BINARY_NAME=web
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	go  build -o ./bin/$(BINARY_NAME) $(BINARY_NAME).go form.go

build-linux:
	CGOENABLED=0 GOOS=linux GOARCH=amd64; go build -o ./bin/$(BINARY_UNIX) $(BINARY_NAME).go form.go

test:
	make -f makefile.test test-selected

grc-test:
	grc make -f makefile.test test-selected

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

## Misc
custom-build:
	make mod tidy; go mod vendor; go build ./...

grc-install:
	sudo apt-get install grc
	make grc-configure

spacer:
	@echo "\n"

get-deps:
	go get -u "github.com/davecgh/go-spew"
	go get -u "github.com/go-chi/chi"
	go get -u "github.com/gorilla/sessions"
	go get -u "github.com/jmoiron/sqlx"
	go get -u "github.com/lib/pq"
	go get -u "github.com/markbates/pkger"
	go get -u "github.com/mattn/go-sqlite3"
	go get -u "gitlab.com/mikrowezel/backend/config"
	go get -u "gitlab.com/mikrowezel/backend/db"
	go get -u "gitlab.com/mikrowezel/backend/db/postgres"
	go get -u "gitlab.com/mikrowezel/backend/log"
	go get -u "gitlab.com/mikrowezel/backend/migration"
	go get -u "gitlab.com/mikrowezel/backend/model"
	go get -u "gitlab.com/mikrowezel/backend/service"
	go get -u "golang.org/x/crypto"
	go get -u "golang.org/x/net"
	go get -u "gopkg.in/check.v1"
