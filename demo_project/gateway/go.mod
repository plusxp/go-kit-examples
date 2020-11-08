module gateway

go 1.12

require (
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/leigg-go/go-util v0.0.4
	go-util v0.0.0-00010101000000-000000000000
	gokit_foundation v0.0.0-00010101000000-000000000000
	gopkg.in/jose.v1 v1.0.0-20161127122323-a941c3995164
	hello v0.0.0-00010101000000-000000000000
)

replace (
	go-util => ../../go-util
	gokit_foundation => ../../gokit_foundation
	hello => ../hello
)
