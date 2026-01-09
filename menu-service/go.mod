module menu-service

go 1.24.0

require (
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.3
	shared v0.0.0
)

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace shared => ../shared
