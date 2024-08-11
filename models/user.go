package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type User struct {
// 	gorm.Model `json:"_"`
// 	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
// 	Username   string    `json:"username"`
// 	Email      string    `json:"email"`
// 	Password   string    `json:"password"`
// }

type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
