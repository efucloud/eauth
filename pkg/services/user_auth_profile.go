package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type UserAuthProfileService struct {
	repo repositories.AuthProfileRepository
}

func (svc *UserAuthProfileService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.AuthProfileRepository{DB: db}
	} else {
		svc.repo = repositories.AuthProfileRepository{DB: config.DBConnect}
	}
}
func (svc *UserAuthProfileService) ChangeStatus(ctx context.Context, model dtos.UserAuthProfileStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	//todo 剔除管理员
	errorData = svc.repo.ChangeStatus(ctx, model)
	return
}
func (svc *UserAuthProfileService) GetUserAuthProfilesByUserId(ctx context.Context, userId uint) (results dtos.UserAuthProfileDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.GetUserAuthProfilesByUserId(ctx, userId)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUserAuthProfile by user id: %d failed, err: %s", userId, errorData.Err.Error())
	}
	return results, errorData
}

func (svc *UserAuthProfileService) GetUserAuthProfileByProviderAndId(ctx context.Context, provider, providerUserId string) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserAuthProfileByProviderAndId(ctx, provider, providerUserId)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
	}
	return result, errorData
}

func (svc *UserAuthProfileService) GetUserAuthProfileByID(ctx context.Context, id uint) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserAuthProfileByID(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUserAuthProfile by id: %d failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *UserAuthProfileService) ListUserAuthProfile(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserAuthProfileDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListUserAuthProfile(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listUserAuthProfile  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *UserAuthProfileService) UpdateUserAuthProfile(ctx context.Context, model dtos.UserAuthProfileUpdate) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateUserAuthProfile: %s failed, err: %s", model.LoginName, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateUserAuthProfile(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateUserAuthProfile: %s failed, err: %s", model.LoginName, errorData.Err.Error())
	}
	return
}
func (svc *UserAuthProfileService) AddUserAuthProfile(ctx context.Context, model dtos.UserAuthProfileCreate) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createUserAuthProfile: %s failed, err: %s", model.LoginName, errorData.Err.Error())

		return
	}
	result, errorData = svc.repo.AddUserAuthProfile(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createUserAuthProfile: %s failed, err: %s", model.LoginName, errorData.Err.Error())
	}
	return
}

func (svc *UserAuthProfileService) DeleteUserAuthProfile(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteUserAuthProfile(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteUserAuthProfile by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}

func (svc *UserAuthProfileService) DeleteUserAuthProfileByUserIds(ctx context.Context, userIds []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteUserAuthProfileByUserIds(ctx, userIds)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteUserAuthProfile by user ids: %v failed, err: %s", userIds, errorData.Err.Error())
	}
	return
}
