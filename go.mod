module gitlab.com/mikrowezel/backend/granica

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.26.5
	github.com/davecgh/go-spew v1.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gorilla/csrf v1.6.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.12.2
	github.com/nicksnyder/go-i18n/v2 v2.0.3
	github.com/satori/go.uuid v1.2.0
	gitlab.com/mikrowezel/backend/config v0.0.0
	gitlab.com/mikrowezel/backend/db v0.0.0-20191014125253-afa2a932cece
	gitlab.com/mikrowezel/backend/db/postgres v0.0.0-20191014125253-afa2a932cece
	gitlab.com/mikrowezel/backend/log v0.0.0
	gitlab.com/mikrowezel/backend/migration v0.0.0-00010101000000-000000000000
	gitlab.com/mikrowezel/backend/model v0.0.0-00010101000000-000000000000
	gitlab.com/mikrowezel/backend/service v0.0.0-20191216010017-5636164a79bf
	gitlab.com/mikrowezel/backend/web v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478 // indirect
	golang.org/x/sys v0.0.0-20190922100055-0a153f010e69 // indirect
	golang.org/x/text v0.3.2
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/stretchr/testify.v1 v1.2.2 // indirect
)

replace gitlab.com/mikrowezel/backend/log => ../log

replace gitlab.com/mikrowezel/backend/config => ../config

replace gitlab.com/mikrowezel/backend/service => ../service

replace gitlab.com/mikrowezel/backend/db => ../db

replace gitlab.com/mikrowezel/backend/db/postgres => ../db/postgres

replace gitlab.com/mikrowezel/backend/migration => ../migration

replace gitlab.com/mikrowezel/backend/model => ../model

replace gitlab.com/mikrowezel/backend/web => ../web
