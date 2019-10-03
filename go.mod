module gitlab.com/mikrowezel/granica

go 1.12

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0 // indirect
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/db v0.0.0
	gitlab.com/mikrowezel/db/postgres v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
	gitlab.com/mikrowezel/migration v0.0.0
	gitlab.com/mikrowezel/service v0.0.0
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace gitlab.com/mikrowezel/log => ../log

replace gitlab.com/mikrowezel/config => ../config

replace gitlab.com/mikrowezel/service => ../service

replace gitlab.com/mikrowezel/db => ../db

replace gitlab.com/mikrowezel/db/postgres => ../db/postgres

replace gitlab.com/mikrowezel/migration => ../migration
