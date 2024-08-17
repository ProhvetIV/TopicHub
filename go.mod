module social-network

go 1.22.1

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.23.0
)

replace github.com/golang-migrate/migrate/v4 => github.com/golang-migrate/migrate/v4 v4.15.2

require golang.org/x/net v0.25.0 // indirect
