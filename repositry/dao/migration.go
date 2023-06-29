package dao

import "easy-drive/repositry/model"

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.File{},
		)
	if err != nil {
		panic(err)
	}
}
