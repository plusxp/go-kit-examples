package db

import (
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_orm/_gorm"
	"github.com/leigg-go/go-util/_redis"
	"go-util/_util"
	"gokit_foundation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hello/config"
)

var logger *gokit_foundation.Logger

type MySQLType struct {
	cli *gorm.DB // 不暴露底层client给其他pkg
}

type RedisType struct {
	cli *redis.Client
}

var GMySQLDao *MySQLType
var GRedisDao *RedisType

// pkg.Init 层可以直接panic if err
func Init(_logger *gokit_foundation.Logger) {
	logger = _logger
	// 注：启动项目前先创建 DB: go_kit_examples
	dsn := "root:123@tcp(127.0.0.1:3306)/go_kit_examples?charset=utf8mb4&parseTime=True&loc=Local"

	// GORM V2 enabled BlockGlobalUpdate mode by default
	gormDB := _gorm.MustInit(mysql.Open(dsn), &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	SingularTable: true,
		//},
	})
	GMySQLDao = &MySQLType{cli: gormDB}
	err := GMySQLDao.Migrate()
	_util.PanicIfErr(err, nil)

	rds := _redis.MustInit(config.GetRedisConf())
	GRedisDao = &RedisType{cli: rds}
}

func Close() {
	db, _ := GMySQLDao.cli.DB()
	err := db.Close()
	if err != nil {
		logger.Log("GMySQLDao.Close", err)
	}
	err = GRedisDao.cli.Close()
	if err != nil {
		logger.Log("GRedisDao.Close", err)
	}
}
