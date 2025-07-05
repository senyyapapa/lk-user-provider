package models

type Student struct {
	ID           int    `gorm:"primaryKey"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	Role         string `gorm:"column:role"`
	Version      int32  `gorm:"column:version"`
}
