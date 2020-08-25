package config

var conf *Conf

/*
存放全局配置
*/

type Conf struct {
	MysqlHost string
	MysqlPort int
}
