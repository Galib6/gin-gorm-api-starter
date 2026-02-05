package seeder

import (
	"errors"

	"myapp/core/entity"
	"myapp/support/constant"

	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	var dummyUsers = []entity.User{
		{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Password: "admin1",
			Role:     constant.EnumRoleAdmin,
		},
		{
			Name:     "User",
			Email:    "user@gmail.com",
			Password: "user1",
			Role:     constant.EnumRoleUser,
		},
	}

	hasTable := db.Migrator().HasTable(&entity.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	for _, data := range dummyUsers {
		var user entity.User
		err := db.Unscoped().Where(&entity.User{Email: data.Email}).First(&user).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
