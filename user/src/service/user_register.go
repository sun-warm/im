package service

import (
	"context"
	"errors"
	"fmt"
	"user/dao"
	"user/dao/model"
	usr "user/generated/user"

	"github.com/go-sql-driver/mysql"
	// MySQL 驱动
)

func (s *server) UserRegister(ctx context.Context, req *usr.UserRegisterRequest) (*usr.UserRegisterResponse, error) {
	if req == nil || req.User == nil {
		return nil, errors.New("User is empty")
	}
	user := req.User
	if user.UserName == "" || user.UserId == "" || user.PassWord == "" {
		return nil, errors.New("UserName, UserId or PassWord is empty")
	}
	insertUser := model.User{
		UserName: user.UserName,
		UserID:   user.UserId,
		PassWord: user.PassWord,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
	result := dao.DB.Table("im_users").Create(&insertUser)
	var mysqlErr *mysql.MySQLError
	if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
		return &usr.UserRegisterResponse{
			ErrorCode: usr.UserErrorCode_DUPLICATE_USERNAME,
		}, errors.New("用户名重复")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("failed to register user: %v", result.Error)
	}
	return &usr.UserRegisterResponse{ErrorCode: usr.UserErrorCode_OK}, nil
}
