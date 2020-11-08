package db

import (
	"hello/pb/gen-go/pbcommon"
)

/*
此文件只存放操作t_user表的方法
所有mysql的操作都应该定义以MySQLType为receiver的Method
*/

// Update user info from mysql
func (mysql *MySQLType) UpdateUserInfo(uid uint) pbcommon.R {
	err := mysql.cli.Model(&User{}).Where("uid=?", uid).Update("name", "晓明").Error
	if err != nil {
		return pbcommon.R_MYSQL_ERR
	}
	return pbcommon.R_OK
}
