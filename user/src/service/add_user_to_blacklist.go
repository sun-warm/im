package service

import (
	"context"
	"errors"
	"user/dao"
	"user/dao/model"
	usr "user/generated/user"

	"gorm.io/gorm/clause"
)

// A申请添加B为好友的逻辑
// 1、A直接写入数据库B对A的好友关系
// 2、查找A对B的好友关系（可能是A把B删了，又申请加好友）
// 3、如果B对A的好友关系是好友，则推送消息给双方，表明又成为了双向好友
// 4、如果B对A的好友关系是未添加，则推送消息给B，表明A申请添加B为好友
func (s *server) AddUserToBlackList(ctx context.Context, req *usr.AddUserToBlackListRequest) (*usr.AddUserToBlackListResponse, error) {
	if req == nil {
		return nil, errors.New("req is empty")
	}
	if req.UserName == "" || req.AddedUserName == "" {
		return nil, errors.New("username or added username is invalid")
	}
	//FIXME：需要校验吗？
	_, err := CheckUserInDB([]string{req.UserName, req.AddedUserName})
	if err != nil {
		return nil, err
	}

	relation := model.UserRelation{
		UserName:         req.UserName,
		RelationUserName: req.AddedUserName,
		RelationType:     int(model.BLACK),
	}
	if err := dao.DB.Table("im_user_relation").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_name"}, {Name: "relation_user_name"}}, // 唯一键字段
			DoUpdates: clause.AssignmentColumns([]string{"relation_type", "updated_at"}),  // 冲突时更新的字段
		},
	).Create(&relation).Error; err != nil {
		return nil, err
	}

	return &usr.AddUserToBlackListResponse{ErrorCode: usr.UserErrorCode_OK}, nil
}
