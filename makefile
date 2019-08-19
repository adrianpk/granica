# Vars
STG_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=mwgranica

# Make
MAKE_CMD=make

# Go
GO_CMD=go

## Docker
DOCKER_CMD=docker

## Kubernetes
KUBECTL_CMD=kubectl

## Helm
HELM_CMD=helm

# Google Cloud
GCLOUD_CMD=gcloud

# Misc
BINARY_NAME=granica
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GO_CMD) build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

build-linux:
	CGOENABLED=0 GOOS=linux GOARCH=amd64; $(GO_CMD) build -o ./bin/$(BINARY_UNIX) ./cmd/$(BINARY_NAME).go

test:
	./scripts/test.sh

clean:
	$(GO_CMD) clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

run:
	./scripts/run.sh

connect-stg:
	$(GCLOUD_CMD) beta container clusters get-credentials ${GC_STG_CLUSTER} --region ${GC_REGION} --project ${GC_STG_PROJECT}

connect-prod:
	$(GCLOUD_CMD) beta container clusters get-credentials ${GC_PROD_CLUSTER} --region ${GC_REGION} --project ${GC_PROD_PROJECT}

build-stg:
	$(MAKE_CMD) build
	$(DOCKER_CMD) login
	$(DOCKER_CMD) build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG) .
	$(DOCKER_CMD) push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG)

build-prod:
	$(MAKE_CMD) build
	$(DOCKER_CMD) login
	$(DOCKER_CMD) build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG) .
	$(DOCKER_CMD) push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG)

template-stg:
	$(HELM_CMD) template --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

template-prod:
	$(HELM_CMD) template --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

install-stg:
	$(MAKE_CMD) connect-stg
	$(HELM_CMD) install --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

install-prod:
	$(MAKE_CMD) connect-prod
	$(HELM_CMD) install --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

delete-stg:
	$(MAKE_CMD) connect-stg
	$(HELM_CMD) del --purge $(IMAGE_NAME)

delete-prod:
	$(MAKE_CMD) connect-prod
	$(HELM_CMD) del --purge $(IMAGE_NAME)

deploy-stg:
	$(MAKE_CMD) build-stg
	$(MAKE_CMD) connect-stg
	$(MAKE_CMD) delete-stg
	$(MAKE_CMD) install-stg

deploy-prod:
	$(MAKE_CMD) build-prod
	$(MAKE_CMD) connect-prod
	$(MAKE_CMD) delete-prod
	$(MAKE_CMD) install-prod

## Misc
custom-build:
	$(MAKE_CMD) mod tidy; go mod vendor; go build ./...

current-conn:
	$(KUBECTL_CMD) config current-context

get-deps:
	$(GO_CMD) get "github.com/cenkalti/backoff"
	$(GO_CMD) get "github.com/go-sql-driver/mysql"
	$(GO_CMD) get "github.com/heptiolabs/healthcheck"
	$(GO_CMD) get "github.com/jmoiron/sqlx"
	$(GO_CMD) get "github.com/lib/pq"
	$(GO_CMD) get "gitlab.com/mikrowezel/config"
	$(GO_CMD) get "gitlab.com/mikrowezel/log"
	$(GO_CMD) get "gitlab.com/mikrowezel/service"
	$(GO_CMD) get "google.golang.org/appengine"
