package model

/*
此文件只存放操作user表的方法
*/

func GetUserInfo(uid uint) *User {
	/*
	 Query from storage.
	*/
	return &User{
		Id:   0,
		Name: "",
		Age:  0,
	}
}
