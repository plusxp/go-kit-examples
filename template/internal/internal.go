package internal

import "log"

/*
internal目录: 存放这个app私有的代码，不能被其他app引入(编译器也会限制)
-	可以包含子目录
*/

func InternalFunc() {

	log.Println("internal/internal.go:this is a InternalFunc")
}
