package service

import (
	"context"
	"errors"
	"user/dao"
	"user/generated/user"
)

func (s *server) CheckUser(ctx context.Context, req *user.CheckUserRequest) (*user.CheckUserResponse, error) {
	if req == nil || req.UserId == nil || len(req.UserId) == 0 {
		return nil, errors.New("UserId is empty")
	}
	var matchedCount int64
	result := dao.DB.Table("user").Where("user_name IN ", req.UserId).Count(&matchedCount)
	if result.Error != nil {
		return nil, errors.New("Get matchedCount error when check user")
	}
	if matchedCount != (int64)(len(req.UserId)) {
		return nil, errors.New("len of UserId is not equal to exist user count")
	}
	return &user.CheckUserResponse{ErrorCode: user.UserErrorCode_OK}, nil
}

func CheckUserInDB(userName []string) (bool, error) {
	var matchedCount int64
	result := dao.DB.Table("user").Where("user_name IN ", userName).Count(&matchedCount)
	if result.Error != nil {
		return false, errors.New("Get matchedCount error when check user")
	}
	if matchedCount != (int64)(len(userName)) {
		return false, errors.New("len of UserId is not equal to exist user count")
	}
	return true, nil
}
