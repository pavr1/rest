module session-service

go 1.25.1

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.46.0
	shared v0.0.0
)

require golang.org/x/sys v0.39.0 // indirect

replace shared => ../shared

replace data-service => ../data-service
