# Granica

 Authentication and authorization service.

<img src="docs/img/users_index.png" width="480">

## Dev branch

* [new/wip at GitLab](https://gitlab.com/mikrowezel/backend/granica/tree/new/wip)
* [new/wip at GitHub](https://github.com/adrianpk/granica/tree/new/wip)

## Installation

[To be completed]

**Test**
```shell
$ make test
```
Use `make grc-test` for colored output.

**Run**
```shell
$ make run
./scripts/run.sh
granica: no process found
5:49PM INF New handler name=migration-handler
5:49PM INF New handler name=repo-handler
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/layout/base.tmpl
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/partial/_flash.tmpl
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/user/_ctxbar.tmpl
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/user/_header.tmpl
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/user/_list.tmpl
5:49PM INF New template file path=gitlab.com/mikrowezel/backend/granica:/assets/web/template/user/index.tmpl
5:49PM INF Dialing to Postgres host="host=localhost port=5432 user=granica password=granica dbname=granica sslmode=disable"
5:49PM INF Postgres connection established
5:49PM INF Repo initializated name=repo-handler
2019/11/14 17:49:22 Migration 'enable_postgis' already applied.
2019/11/14 17:49:22 Migration 'create_users_table' already applied.
2019/11/14 17:49:22 Migration 'create_accounts_table' already applied.
5:49PM INF Migrator initializated name=migration-handler
5:49PM INF JSON REST Server initializing port=:8081
5:49PM INF Web server initializing port=:8080
```

## Deployment

[To be completed]

## Packages

**Worker**

[Auth](pkg/auth/readme.md)

## Helpers

[Supervisord and Gulp](docs/draft/helpers.md)
