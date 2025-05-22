package model

import "time"

// RelationType 定义关系类型的枚举
type RelationType int

const (
	UNRELATED RelationType = 0 // 未建立任何关系
	FRIEND    RelationType = 1 // 好友
	BLACK     RelationType = 2 // 黑名单
	DELETED   RelationType = 3 // 删除
)

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

// RelationType 1代表好友 2代表黑名单 3代表删除 0代表未建立任何关系
type UserRelation struct {
	UserName         string    `gorm:"index:user_name_idx;uniqueIndex:user_relation_unique;column:user_name;type:varchar(100);not null"`
	RelationUserName string    `gorm:"uniqueIndex:user_relation_unique;column:relation_user_name;type:varchar(100);not null"`
	RelationType     int       `gorm:"column:relation_type;type:int;not null"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// other fields
}

func (User) TableName() string {
	return "im_users"
}

func (UserRelation) TableName() string {
	return "im_user_relation"
}
