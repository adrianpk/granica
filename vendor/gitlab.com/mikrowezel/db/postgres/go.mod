module gitlab.com/mikrowezel/db/postgres

go 1.12

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.2.0
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
	gitlab.com/mikrowezel/service v0.0.0
	google.golang.org/appengine v1.6.1 // indirect
)

replace gitlab.com/mikrowezel/log => ../../log

replace gitlab.com/mikrowezel/config => ../../config

replace gitlab.com/mikrowezel/service => ../../service
