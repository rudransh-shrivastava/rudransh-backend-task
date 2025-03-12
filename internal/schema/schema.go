package schema

import "time"

type Role string

const (
	Student  Role = "STUDENT"
	Educator Role = "EDUCATOR"
	Admin    Role = "ADMIN"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	UID       string    `gorm:"uniqueIndex"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"not null" json:"name"`
	Role      Role      `gorm:"not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Course struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
