# Granica

Authentication and authorization service.

<img src="docs/img/users_index.png" width="480">

## Dev branch

- [new/wip at GitLab](https://gitlab.com/mikrowezel/backend/granica/tree/new/wip)
- [new/wip at GitHub](https://github.com/adrianpk/granica/tree/new/wip)

Repository mirroring (GitLab -> GitHub) seems to work erratically from time to time and for that reason the latter and because of this the latter may not be showing the current state of development.

## Installation

[TODO: Create database steps]

```shell
$ git clone https://gitlab.com/mikrowezel/backend/granica
$ make run
```

[TODO: additional steps]

**Test**

```shell
$ make test
```

Use `make grc-test` for colored output.

**Run**

```shell
$ make run
./scripts/run.sh
2:27AM INF New handler name=migration-handler
2:27AM INF New handler name=repo-handler
2:27AM INF Template file path=layout/base.tmpl
2:27AM INF Template file path=user/_ctxbar.tmpl
2:27AM INF Template file path=user/_flash.tmpl
2:27AM INF Template file path=user/_header.tmpl
2:27AM INF Template file path=user/_list.tmpl
2:27AM INF Template file path=user/index.tmpl
2:27AM INF Template processed template=./assets/web/embed/template/user/index.tmpl
2:27AM INF Dialing to Postgres host="host=localhost port=5432 user=granica password=granica dbname=granica sslmode=disable"
2:27AM INF Postgres connection established
2:27AM INF Repo initializated name=repo-handler
2019/11/15 02:27:10 Migration 'enable_postgis' already applied.
2019/11/15 02:27:10 Migration 'create_users_table' already applied.
2019/11/15 02:27:10 Migration 'create_accounts_table' already applied.
2:27AM INF Migrator initializated name=migration-handler
2:27AM INF JSON REST Server initializing port=:8081
2:27AM INF Web server initializing port=:8080
```

## Deployment

[To be completed]

## Packages

**Worker**

[Auth](pkg/auth/readme.md)

## Helpers

[Supervisord and Gulp](docs/draft/helpers.md)
