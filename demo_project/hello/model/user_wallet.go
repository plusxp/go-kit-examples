package model

/*
此文件只存放操作user_wallet表的方法
*/

func GetUserWallet(uid uint) *UserWallet {
	/*
	 Query from storage.
	*/
	return &UserWallet{
		Uid:     0,
		Balance: 0,
		Status:  0,
	}
}
