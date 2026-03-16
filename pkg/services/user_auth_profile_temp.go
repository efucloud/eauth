package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type UserAuthProfileTempService struct {
	repo repositories.AuthProfileTempRepository
}

func (svc *UserAuthProfileTempService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.AuthProfileTempRepository{DB: db}
	} else {
		svc.repo = repositories.AuthProfileTempRepository{DB: config.DBConnect}
	}
}

func (svc *UserAuthProfileTempService) GetUserAuthProfileTempByCode(ctx context.Context, code string) (result dtos.UserAuthProfileTempDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserAuthProfileTempByCode(ctx, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUserAuthProfileTemp by code: %s failed, err: %s", code, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *UserAuthProfileTempService) ListUserAuthProfileTemp(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserAuthProfileTempDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListUserAuthProfileTemp(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listUserAuthProfileTemp  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *UserAuthProfileTempService) AddUserAuthProfileTemp(ctx context.Context, model dtos.UserAuthProfileTempCreate) (result dtos.UserAuthProfileTempDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createUserAuthProfileTemp: %s failed, err: %s", model.LoginName, errorData.Err.Error())

		return
	}
	result, errorData = svc.repo.AddUserAuthProfileTemp(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createUserAuthProfileTemp: %s failed, err: %s", model.LoginName, errorData.Err.Error())
	}
	return
}

func (svc *UserAuthProfileTempService) DeleteUserAuthProfileTemp(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteUserAuthProfileTemp(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteUserAuthProfileTemp by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}
