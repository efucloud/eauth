package services

import (
	"context"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type ProviderSamlService struct {
	repo repositories.ProviderSamlRepository
}

func (svc *ProviderSamlService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ProviderSamlRepository{DB: db}
	} else {
		svc.repo = repositories.ProviderSamlRepository{DB: config.DBConnect}
	}
}

func (svc *ProviderSamlService) GetProviderSamlByCategory(ctx context.Context, category string) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetProviderSamlByCategory(ctx, category)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getProviderSaml by category: %s failed, err: %s", category, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) GetProviderSamlById(ctx context.Context, id uint) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetProviderSamlById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getProviderSaml by id: %d failed, err: %s", id, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) ListProviderSaml(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderSamlDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListProviderSaml(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listProviderSaml query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) UpdateProviderSaml(ctx context.Context, model dtos.ProviderSamlUpdate) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateProviderSaml: %s failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateProviderSaml(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateProviderSaml: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) AddProviderSaml(ctx context.Context, model dtos.ProviderSamlCreate) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createProviderSaml: %s failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.AddProviderSaml(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createProviderSaml: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) DeleteProviderSaml(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteProviderSaml(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteProviderSaml by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}

func (svc *ProviderSamlService) ChangeStatus(ctx context.Context, model dtos.ProviderSamlStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData = svc.repo.ChangeStatus(ctx, model)
	return
}
