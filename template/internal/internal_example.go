package internal

import "log"

/*
internal目录: 存放这个app私有的代码，不能被其他app引入(IDE/编译器都会限制，提示)
-	不能import项目下其他sub-pkg，util除外；一般都是被import
-	可以包含子目录
*/

func Func() {
	log.Println("internal/internal_example.go:this is a internal Func")
}
