module gitlab.com/mikrowezel/service

go 1.12

require (
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/prometheus/client_golang v1.1.0 // indirect
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
)

replace gitlab.com/mikrowezel/log => ../log

replace gitlab.com/mikrowezel/config => ../config

replace gitlab.com/mikrowezel/service => ../service
