# Granica

Authentication and authorization service.

The final intention of this project is to build a complete authentication and authorization system that can be easily deployed as a microservice.

But at the same time it is also a way to find a modular and systematized foundation to build microservices and micromonoliths; therefore, as a pattern emerges some features, if not alli, will serve to update a series of libraries and a [command-line interface](https://gitlab.com/mikrowezel/backend/cli) that will allow to automate as much as possible the generation of Go-based services while trying not to detach too much from the particular idiosyncrasy of this platform.

<img src="docs/img/users_index.png" width="480">

[Screenshots](docs/screenshots.md)

## Dev branch

- [new/wip at GitLab](https://gitlab.com/mikrowezel/backend/granica/tree/new/wip)
- [new/wip at GitHub](https://github.com/adrianpk/granica/tree/new/wip)

## Changelog

* [20191126 - Simplified flash messages handling](/docs/changelog.md#20191126)
* [20191124 - Embedded translations and form data session store](/docs/changelog.md#20191124)
* [20191123 - Internationalization](/docs/changelog.md#20191123)

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
5:31PM INF New handler name=migration-handler
5:31PM INF New handler name=repo-handler
5:31PM INF Cookie store key value=pVuOO7ZPNBnqTb5o13JsBMOPcPAe4pxY
5:31PM INF Reading template path=layout/base.tmpl
5:31PM INF Reading template path=user/_ctxbar.tmpl
5:31PM INF Reading template path=user/_flash.tmpl
5:31PM INF Reading template path=user/_form.tmpl
5:31PM INF Reading template path=user/_header.tmpl
5:31PM INF Reading template path=user/_list.tmpl
5:31PM INF Reading template path=user/create.tmpl
5:31PM INF Reading template path=user/index.tmpl
5:31PM INF Parsed template set path=user/create.tmpl
5:31PM INF Parsed template set path=user/index.tmpl
5:31PM INF Dialing to Postgres host="host=localhost port=5432 user=dbuser password=dbpass dbname=granica sslmode=disable"
5:31PM INF Postgres connection established
5:31PM INF Repo initializated name=repo-handler
2019/11/21 17:31:06 Migration 'enable_postgis' already applied.
2019/11/21 17:31:06 Migration 'create_users_table' already applied.
2019/11/21 17:31:06 Migration 'create_accounts_table' already applied.
5:31PM INF Migrator initializated name=migration-handler
5:31PM INF JSON REST Server initializing port=:8081
5:31PM INF Web server initializing port=:8080
```

## Deployment

[To be completed]

## Packages

**Worker**

[Auth](pkg/auth/readme.md)

## Helpers

[Supervisord and Gulp](docs/draft/helpers.md)
