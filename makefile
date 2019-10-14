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
	go test -v -count=1 -timeout=5s  ./internal/repo/user_test.go

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

run:
	./scripts/run.sh

rest-create-user:
	./scripts/rest/create_user.zsh

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

# Tests
test-create-user:
	go test -v -run TestCreateUser -count=1 -timeout=5s  ./internal/repo/user_test.go

test-get-users:
	go test -v -run TestGetAllUsers -count=1 -timeout=5s  ./internal/repo/user_test.go

test-get-user-by-id:
	go test -v -run TestGetUserByID -count=1 -timeout=5s  ./internal/repo/user_test.go

test-get-user-by-slug:
	go test -v -run TestGetUserBySlug -count=1 -timeout=5s  ./internal/repo/user_test.go

test-get-user-by-username:
	go test -v -run TestGetUserByUsername -count=1 -timeout=5s  ./internal/repo/user_test.go

test-update-user:
	go test -v -run TestUpdateUser -count=1 -timeout=5s  ./internal/repo/user_test.go

test-delete-user:
	go test -v -run TestDeleteUser -count=1 -timeout=5s  ./internal/repo/user_test.go

## Misc
custom-build:
	make mod tidy; go mod vendor; go build ./...

current-conn:
	kubectl config current-context

get-deps:
	go get -u "github.com/go-chi/chi"
	go get -u "github.com/jmoiron/sqlx"
	go get -u "github.com/kr/pretty"
	go get -u "github.com/lib/pq"
	go get -u "github.com/mattn/go-sqlite3"
	go get -u "github.com/satori/go.uuid"
	go get -u "gitlab.com/mikrowezel/backend/config"
	go get -u "gitlab.com/mikrowezel/backend/db"
	go get -u "gitlab.com/mikrowezel/backend/db/postgres"
	go get -u "gitlab.com/mikrowezel/backend/log"
	go get -u "gitlab.com/mikrowezel/backend/migration"
	go get -u "gitlab.com/mikrowezel/backend/service"
	go get -u "golang.org/x/crypto"
	go get -u "golang.org/x/net"
	go get -u "gopkg.in/check.v1"
