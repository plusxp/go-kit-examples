package main

import (
	"go-kit-examples/template/internal"
	"go-kit-examples/template/pkg"
	"log"
)

func main() {
	internal.Func()
	log.Println("main/some_svc/main.go: -------- running")
	pkg.PublicFunc()

	log.Println("main/some_svc/main.go: -------- exited")
}
