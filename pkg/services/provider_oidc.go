package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type ProviderOidcService struct {
	repo repositories.ProviderOidcRepository
}

func (svc *ProviderOidcService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ProviderOidcRepository{DB: db}
	} else {
		svc.repo = repositories.ProviderOidcRepository{DB: config.DBConnect}
	}
}

func (svc *ProviderOidcService) GetProviderOidcByCategory(ctx context.Context, category string) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetProviderOidcByCategory(ctx, category)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getProviderOidc by category: %s failed, err: %s", category, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *ProviderOidcService) GetProviderOidcById(ctx context.Context, id uint) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetProviderOidcById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getProviderOidc by id: %d failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *ProviderOidcService) ListProviderOidc(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderOidcDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListProviderOidc(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listProviderOidc  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *ProviderOidcService) UpdateProviderOidc(ctx context.Context, model dtos.ProviderOidcUpdate) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateProviderOidc: %s failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateProviderOidc(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateProviderOidc: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}
func (svc *ProviderOidcService) AddProviderOidc(ctx context.Context, model dtos.ProviderOidcCreate) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createProviderOidc: %s failed, err: %s", model.Name, errorData.Err.Error())

		return
	}
	result, errorData = svc.repo.AddProviderOidc(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createProviderOidc: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}

func (svc *ProviderOidcService) DeleteProviderOidc(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteProviderOidc(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteProviderOidc by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}

func (svc *ProviderOidcService) ChangeStatus(ctx context.Context, model dtos.ProviderOidcStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData = svc.repo.ChangeStatus(ctx, model)
	return
}
