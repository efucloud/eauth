package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type UserTokenService struct {
	repo repositories.UserTokenRepository
}

func (svc *UserTokenService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.UserTokenRepository{DB: db}
	} else {
		svc.repo = repositories.UserTokenRepository{DB: config.DBConnect}
	}
}

func (svc *UserTokenService) GetUserTokensByreRefreshToken(ctx context.Context, refreshToken string) (results dtos.UserTokenDetail, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.GetUserTokensByreRefreshToken(ctx, refreshToken)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get OrgAccountAuthProfile by  refreshToken: %s failed, err: %s", refreshToken, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *UserTokenService) GetUserTokensByUserId(ctx context.Context, userId string) (results dtos.UserTokenDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	userSvc := UserService{}
	user, _ := userSvc.GetShortUserByID(ctx, userId)
	results, errorData = svc.repo.GetUserTokensByUserId(ctx, userId)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUserAuthProfile by user id: %s failed, err: %s", userId, errorData.Err.Error())
	}
	for i, _ := range results.Data {
		results.Data[i].User = user
	}
	return results, errorData
}

func (svc *UserTokenService) GetUserTokenDetailById(ctx context.Context, id string) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserTokenDetailById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUserTokenDetail by id: %s failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *UserTokenService) ListUserToken(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserTokenDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListUserToken(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listUserToken  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *UserTokenService) UpdateUserToken(ctx context.Context, model dtos.UserTokenUpdate) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateUserToken: %s failed, err: %s", model.UserId, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateUserToken(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateUserToken: %s failed, err: %s", model.UserId, errorData.Err.Error())

	}
	return
}
func (svc *UserTokenService) AddUserToken(ctx context.Context, model dtos.UserTokenCreate) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createUserToken: %s failed, err: %s", model.UserId, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.AddUserToken(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createUserToken: %s failed, err: %s", model.UserId, errorData.Err.Error())
	}
	return
}

func (svc *UserTokenService) DeleteUserToken(ctx context.Context, ids []string) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteUserToken(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteUserToken by id: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}

func (svc *UserTokenService) DeleteExpireRecord(ctx context.Context) {
	svc.init(ctx)
	svc.repo.DeleteExpireRecord(ctx)
}
