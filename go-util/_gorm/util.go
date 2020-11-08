package _gorm

import "gorm.io/gorm"

func IsDBErr(err error) bool {
	return err != nil && err != gorm.ErrRecordNotFound
}
