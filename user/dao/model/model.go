package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"column:user_id;type:varchar(100);not null"`
	UserName  string    `gorm:"uniqueIndex;column:user_name;type:varchar(100);not null"`
	PassWord  string    `gorm:"column:pass_word;type:varchar(100);not null"`
	Avatar    string    `gorm:"column:avatar;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// other fields
}

func (User) TableName() string {
	return "im_users"
}
