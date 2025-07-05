package database

import "main/internal/database/models"

type Database interface {
	AuthUser(email string) *models.Student
}
