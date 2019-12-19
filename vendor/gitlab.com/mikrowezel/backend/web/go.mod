module gitlab.com/mikrowezel/backend/web

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gorilla/csrf v1.6.1
	github.com/gorilla/schema v1.1.0
	github.com/gorilla/sessions v1.2.0
	github.com/markbates/pkger v0.12.2
	github.com/nicksnyder/go-i18n/v2 v2.0.3
	gitlab.com/mikrowezel/backend/config v0.0.0
	gitlab.com/mikrowezel/backend/log v0.0.0
	golang.org/x/text v0.3.2
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace gitlab.com/mikrowezel/backend/log => ../log

replace gitlab.com/mikrowezel/backend/config => ../config

replace gitlab.com/mikrowezel/backend/service => ../service

replace gitlab.com/mikrowezel/backend/db => ../db

replace gitlab.com/mikrowezel/backend/db/postgres => ../db/postgres

replace gitlab.com/mikrowezel/backend/migration => ../migration

replace gitlab.com/mikrowezel/backend/model => ../model
