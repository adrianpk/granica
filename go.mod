module gitlab.com/mikrowezel/granica

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.2.0
	github.com/satori/go.uuid v1.2.0
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/db v0.0.0
	gitlab.com/mikrowezel/db/postgres v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
	gitlab.com/mikrowezel/service v0.0.0
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
)

replace gitlab.com/mikrowezel/log => ../log

replace gitlab.com/mikrowezel/config => ../config

replace gitlab.com/mikrowezel/service => ../service

replace gitlab.com/mikrowezel/db => ../db

replace gitlab.com/mikrowezel/db/postgres => ../db/postgres
