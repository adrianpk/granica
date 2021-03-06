# Vars
STG_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=mw_granica

# Misc
BINARY_NAME=granica
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	go  build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

build-linux:
	CGOENABLED=0 GOOS=linux GOARCH=amd64; go build -o ./bin/$(BINARY_UNIX) ./cmd/$(BINARY_NAME).go

test:
	make -f makefile.test test-selected

grc-test:
	grc make -f makefile.test test-selected

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

run:
	./scripts/run.sh

package-resources:
	pkger -include /assets/web/embed -o pkg/auth/web

list-package-resources:
	pkger list -include /assets/web --json

# Generators
gen-resource:
	mw generate all assets/gen/resource.yaml

# Cloud
connect-stg:
	gcloud beta container clusters get-credentials ${GC_STG_CLUSTER} --region ${GC_REGION} --project ${GC_STG_PROJECT}

connect-prod:
	gcloud  beta container clusters get-credentials ${GC_PROD_CLUSTER} --region ${GC_REGION} --project ${GC_PROD_PROJECT}

build-stg:
	make build
	docker login
	docker build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG) .
	docker push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(STG_TAG)

build-prod:
	make build
	docker login
	docker build -t ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG) .
	docker push ${DOCKER_ACCOUNT}/$(IMAGE_NAME):$(PROD_TAG)

template-stg:
	helm template --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

template-prod:
	helm template --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

install-stg:
	make connect-stg
	helm install --name $(IMAGE_NAME) -f ./deployments/helm/values-stg.yaml ./deployments/helm

install-prod:
	make connect-prod
	helm install --name $(IMAGE_NAME) -f ./deployments/helm/values-prod.yaml ./deployments/helm

delete-stg:
	make connect-stg
	helm del --purge $(IMAGE_NAME)

delete-prod:
	make connect-prod
	helm del --purge $(IMAGE_NAME)

deploy-stg:
	make build-stg
	make connect-stg
	make delete-stg
	make install-stg

deploy-prod:
	make build-prod
	make connect-prod
	make delete-prod
	make install-prod

## Misc
custom-build:
	make mod tidy; go mod vendor; go build ./...

clean-and-run:
	clear
	make package-resources
	make run

gen-sample:
	mw generate Sample --all --force

current-conn:
	kubectl config current-context

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
