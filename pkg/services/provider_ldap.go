package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type ProviderLdapService struct {
	repo repositories.ProviderLdapRepository
}

func (svc *ProviderLdapService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ProviderLdapRepository{DB: db}
	} else {
		svc.repo = repositories.ProviderLdapRepository{DB: config.DBConnect}
	}
}

func (svc *ProviderLdapService) GetProviderLdapById(ctx context.Context, id string) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetProviderLdapById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getProviderLdap by id: %s failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *ProviderLdapService) ListProviderLdap(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderLdapDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListProviderLdap(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listProviderLdap  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *ProviderLdapService) UpdateProviderLdap(ctx context.Context, model dtos.ProviderLdapUpdate) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateProviderLdap: %s failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateProviderLdap(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateProviderLdap: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}
func (svc *ProviderLdapService) AddProviderLdap(ctx context.Context, model dtos.ProviderLdapCreate) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createProviderLdap: %s failed, err: %s", model.Name, errorData.Err.Error())

		return
	}
	result, errorData = svc.repo.AddProviderLdap(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createProviderLdap: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}

func (svc *ProviderLdapService) DeleteProviderLdap(ctx context.Context, ids []string) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteProviderLdap(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteProviderLdap by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}
