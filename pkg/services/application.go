package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
	"time"
)

type ApplicationService struct {
	repo repositories.ApplicationRepository
}

func (svc *ApplicationService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ApplicationRepository{DB: db}
	} else {
		svc.repo = repositories.ApplicationRepository{DB: config.DBConnect}
	}
}
func (svc *ApplicationService) ChangeStatus(ctx context.Context, model dtos.ApplicationStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	model.UpdatedAt = time.Now()
	//todo 剔除管理员
	errorData = svc.repo.ChangeStatus(ctx, model)
	return
}
func (svc *ApplicationService) GetApplicationById(ctx context.Context, id uint) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetApplicationById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get Application by id: %d failed, err: %s", id, errorData.Err.Error())
	}

	return result, errorData
}
func (svc *ApplicationService) GetApplicationByClientId(ctx context.Context, clientId string) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetApplicationByClientId(ctx, clientId)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get Application by clientId: %s failed, err: %s", clientId, errorData.Err.Error())
	}

	return result, errorData
}

func (svc *ApplicationService) GetApplicationByCode(ctx context.Context, code string) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetApplicationByCode(ctx, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get Application by code: %s failed, err: %s", code, errorData.Err.Error())
	}

	return result, errorData
}
func (svc *ApplicationService) ListApplication(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ApplicationDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListApplication(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("list Application query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}

	return results, errorData
}

func (svc *ApplicationService) UpdateApplication(ctx context.Context, model dtos.ApplicationUpdate) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("Application: %s update failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateApplication(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("Application: %s update failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}
func (svc *ApplicationService) AddApplication(ctx context.Context, model dtos.ApplicationCreate) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("Application: %s create failed, err: %s", model.Name, errorData.Err.Error())
		return
	}

	result, errorData = svc.repo.AddApplication(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("Application: %s create failed, err: %s", model.Name, errorData.Err.Error())
	}

	return
}

func (svc *ApplicationService) DeleteApplication(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteApplication(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}
