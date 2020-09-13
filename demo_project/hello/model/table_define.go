package model

// user 表
type User struct {
	Id   uint
	Name string
	Age  uint8
}

// user wallet表
type UserWallet struct {
	Uid     uint
	Balance uint
	Status  int8
}

func Migrate() {
	// Create these tables If not exist

	// Add index for these tables
}
