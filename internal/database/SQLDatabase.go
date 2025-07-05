package database

import (
	"main/internal/database/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SQLDatabase struct {
	DB *gorm.DB
}

func GetDB(Db_url string) Database {
	db, err := gorm.Open(mysql.Open(Db_url), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SQLDatabase{DB: db}
}

func (db *SQLDatabase) AuthUser(email string) *models.Student {
	student := &models.Student{}
	db.DB.Where("email = ?", email).First(student)
	return student
}
