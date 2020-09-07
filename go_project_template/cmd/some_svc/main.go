package main

import (
	"go_pj_template/internal"
	"go_pj_template/pkg"
	"log"
)

func main() {
	internal.Func()
	log.Println("main/some_svc/main.go: -------- running")
	pkg.PublicFunc()

	log.Println("main/some_svc/main.go: -------- exited")
}
