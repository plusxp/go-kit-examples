package pkg

import (
	"go-kit-examples/template/internal"
	"log"
)

/*
pkg目录：存放可以被重用的代码
-	不过请注意, 区别于util目录，这个目录存放是具有一定业务相关性的代码
-	可以包含子目录
*/

func PublicFunc() {
	internal.InternalFunc()
	log.Println("pkg/share.go: this is a publicFunc")
}
