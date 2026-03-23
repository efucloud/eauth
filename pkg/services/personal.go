package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
)

type PersonalService struct {
}

func (svc *PersonalService) UpdateUserAvatar(ctx context.Context, userId uint, avatarAddress string) (errorData common.ErrorData) {
	userSvc := UserService{}
	return userSvc.UpdateUserAvatar(ctx, userId, avatarAddress)
}
func (svc *PersonalService) SetPassword(ctx context.Context, userId uint, model dtos.SetPassword) (errorData common.ErrorData) {
	userSvc := UserService{}
	if len(model.OldPassword) > 0 {
		user, _ := userSvc.GetUserByID(ctx, userId)
		errorData.Err = common.ComparePassword(user.PasswordStore, model.OldPassword, config.PasswordSalt)
		if errorData.IsNotNil() {
			return errorData
		}
	}
	var pwd dtos.UserResetPassword
	pwd.Password = model.NewPassword
	pwd.ID = userId
	return userSvc.ResetPassword(ctx, pwd)
}
