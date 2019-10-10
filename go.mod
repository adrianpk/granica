module gitlab.com/mikrowezel/granica

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/twinj/uuid v1.0.0
	gitlab.com/mikrowezel/backend/config v0.0.0
	gitlab.com/mikrowezel/backend/db v0.0.0-20191008095614-6aaadda5b1e2
	gitlab.com/mikrowezel/backend/db/postgres v0.0.0-20191008095614-6aaadda5b1e2
	gitlab.com/mikrowezel/backend/log v0.0.0
	gitlab.com/mikrowezel/backend/migration v0.0.0-00010101000000-000000000000
	gitlab.com/mikrowezel/backend/service v0.0.0-20191008112211-3ae14b5bbc28
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace gitlab.com/mikrowezel/backend/log => ../log

replace gitlab.com/mikrowezel/backend/config => ../config

replace gitlab.com/mikrowezel/backend/service => ../service

replace gitlab.com/mikrowezel/backend/db => ../db

replace gitlab.com/mikrowezel/backend/migration => ../migration
