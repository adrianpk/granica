#!/bin/sh
# Build
# ./scripts/build.sh

# Free ports
killall -9 granica

# Set environment variables
REV=$(eval git rev-parse HEAD)
# Service
export GRN_SVC_NAME="granica"
export GRN_SVC_REVISION=$REV
export GRN_SVC_PINGPORT=8090
# Servers
export GRN_WEB_SERVER_PORT=8080
export GRN_JSONREST_SERVER_PORT=8081
export GRN_WEB_COOKIESTORE_KEY="pVuOO7ZPNBnqTb5o13JsBMOPcPAe4pxY"
# Postgres
export GRN_PG_SCHEMA="public"
export GRN_PG_DATABASE="granica"
export GRN_PG_HOST="localhost"
export GRN_PG_PORT="5432"
export GRN_PG_USER="granica"
export GRN_PG_PASSWORD="granica"
export GRN_PG_BACKOFF_MAXTRIES="3"
# Amazon SES MAiler
  # These are sample not usable keys
export AWS_ACCESS_KEY_ID=EIIAHI5FF3A2OG3MJEX5
export AWS_SECRET_KEY=8BiWmd5Hdgmk2rR4pwG332bHwvLGiJOoxLLtDy12

# Switches
export GRN_APP_USERNAME_UPDATABLE=false

go build -o ./bin/granica ./cmd/granica.go
./bin/granica
# go run -race main.go
