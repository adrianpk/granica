# Vars
STAGE_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=mwconfig

# Accounts
DOCKERHUB_USER=n/a

# Go
MAKE_CMD=make
# Go
GO_CMD=go

## Docker
DOCKER_CMD=docker

# Gcloud
GCLOUD_CMD=gcloud

# Misc
BINARY_NAME=config
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GO_CMD) build ./...

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_UNIX) -v

test:
# Be sure to set up environment variables that apply for your case.
# PROVIDER_ID_KEY_1, PROVIDER_API_KEY_2, AWS_ACCESS_KEY_ID, AWS_SECRET_KEY
	$(GO_CMD) test -v ./...

clean:
	$(GO_CMD) clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

deps:
	$(GO_CMD) get -u github.com/cenkalti/backoff
	$(GO_CMD) get -u github.com/google/uuid
	$(GO_CMD) get -u github.com/mitchellh/mapstructure
	$(GO_CMD) get -u github.com/streadway/amqp
	$(GO_CMD) get -u gitlab.com/mikrowezel/log

