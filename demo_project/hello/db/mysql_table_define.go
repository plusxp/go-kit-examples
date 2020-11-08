package db

import "gorm.io/gorm"

// user 表
type User struct {
	gorm.Model
	Name string `gorm:"index:idx_name_age,unique"`
	Age  uint8  `gorm:"index:idx_name_age"`
}

func (User) TableName() string {
	return "t_user"
}

// See also https://gorm.io/docs/migration.html
func (mysql *MySQLType) Migrate() error {
	db := mysql.cli.Set("gorm:table_options", `ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	// Create these tables If not exist
	err := db.AutoMigrate(&User{})
	return err
}
