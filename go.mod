module gitlab.com/mikrowezel/granica

go 1.12

require (
	github.com/go-sql-driver/mysql v1.4.1
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/db/postgres v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
	gitlab.com/mikrowezel/service v0.0.0
)

replace gitlab.com/mikrowezel/log => ../log

replace gitlab.com/mikrowezel/config => ../config

replace gitlab.com/mikrowezel/service => ../service

replace gitlab.com/mikrowezel/db/postgres => ../db/postgres
