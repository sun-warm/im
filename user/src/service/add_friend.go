package service

import (
	"context"
	"errors"
	"fmt"
	"user/client"
	"user/dao"
	"user/dao/model"
	"user/generated/push"
	usr "user/generated/user"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// A申请添加B为好友的逻辑
// 1、A直接写入数据库B对A的好友关系
// 2、查找A对B的好友关系（可能是A把B删了，又申请加好友）
// 3、如果B对A的好友关系是好友，则推送消息给双方，表明又成为了双向好友
// 4、如果B对A的好友关系是未添加，则推送消息给B，表明A申请添加B为好友
func (s *server) AddFriend(ctx context.Context, req *usr.AddFriendRequest) (*usr.AddFriendResponse, error) {
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

	//直接写在relationDB中B对A的好友关系，然后推给B，如果B收到消息可以同时请求该接口，建立A对B的好友关系
	//最后做一个校验，如果双方已经确认好友关系，推一条消息给双方

	var reverseRelation model.UserRelation
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 插入或更新 A 对 B 的好友关系
		relation := model.UserRelation{
			UserName:         req.UserName,
			RelationUserName: req.AddedUserName,
			RelationType:     int(model.FRIEND),
		}
		if err := tx.Table("im_user_relation").Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_name"}, {Name: "relation_user_name"}}, // 唯一键字段
				DoUpdates: clause.AssignmentColumns([]string{"relation_type", "updated_at"}),  // 冲突时更新的字段
			},
		).Create(&relation).Error; err != nil {
			return fmt.Errorf("添加好友时写入数据库失败: %v", err)
		}

		// 查询 B 对 A 的好友关系
		if err := tx.Table("im_user_relation").
			Where("user_name = ? AND relation_user_name = ?", req.AddedUserName, req.UserName).
			First(&reverseRelation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果记录不存在，可以继续处理
				fmt.Println("好友关系记录不存在")
			} else {
				return err
			}
		}
		// 如果需要，可以在事务中继续处理其他逻辑
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 双方都添加了好友关系，推送消息给双方
	if reverseRelation.RelationType == int(model.FRIEND) {
		//TODO:push消息给前端
		request := push.PushMessageRequest{
			Content:  "",
			UserName: req.AddedUserName,
			Receiver: req.UserName,
		}

		_, err = client.PushServiceClient.Client.PushMessage(ctx, &request)
		if err != nil {
			//FIXME：其实应该丢一个消息队列里，支持延迟消费
		}

	} else {
		//TODO:push消息给前端
		request := push.PushMessageRequest{
			Content:  "",
			UserName: req.UserName,
			Receiver: req.AddedUserName,
		}

		_, err = client.PushServiceClient.Client.PushMessage(ctx, &request)
		if err != nil {
			//FIXME：其实应该丢一个消息队列里，支持延迟消费
			return nil, errors.New("添加好友失败，请稍后重试")
		}
	}
	return &usr.AddFriendResponse{ErrorCode: usr.UserErrorCode_OK}, nil
}
