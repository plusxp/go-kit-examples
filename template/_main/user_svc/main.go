package main

import (
	"go-kit-examples/template/pkg"
	"log"
)

/*
cmd目录结构(应该仅包含main.go一个文件，从/pkg或/internal目录导入需要的方法)：
	_main/your_app_name/_main.go

- build
$cd /_main/your_app_name && go build .
*/

func main() {
	log.Println("_main/user_svc/_main.go: --------")
	pkg.PublicFunc()

}
