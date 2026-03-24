package services

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type AppAuthRecordService struct {
	repo repositories.AppAuthRecordRepository
}

func (svc *AppAuthRecordService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.AppAuthRecordRepository{DB: db}
	} else {
		svc.repo = repositories.AppAuthRecordRepository{DB: config.DBConnect}
	}
}
func (svc *AppAuthRecordService) GetAppAuthRecordById(ctx context.Context, id string) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetAppAuthRecordById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get AppAuthRecord by id: %s failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}
func (svc *AppAuthRecordService) GetAppAuthRecordByCode(ctx context.Context, code string) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetAppAuthRecordByCode(ctx, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get AppAuthRecord by code: %s failed, err: %s", code, errorData.Err.Error())
	}
	return result, errorData
}
func (svc *AppAuthRecordService) GetAllAppAuthRecords(ctx context.Context) (results dtos.AppAuthRecordDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	//获取所有的
	current := 1
	pageSize := 1000
	order := "id desc"
	query := ""
	var queryArgs []interface{}
	var records []*dtos.AppAuthRecordDetail
	for {
		r, e := svc.repo.ListAppAuthRecord(ctx, current, pageSize, order, query, queryArgs)
		if e.IsNotNil() {
			break
		}
		if len(r.Data) > 0 {
			records = append(records, r.Data...)
		}
		if int64(len(records)) == r.Total {
			break
		}
		current += 1
	}
	results.Data = records
	results.Total = int64(len(results.Data))
	return results, errorData
}
func (svc *AppAuthRecordService) ListAppAuthRecord(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.AppAuthRecordDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListAppAuthRecord(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("list AppAuthRecord query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}

	return results, errorData
}

func (svc *AppAuthRecordService) AddAppAuthRecord(ctx context.Context, model dtos.AppAuthRecordCreate) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("application: %s user: %s create failed, err: %s", model.ApplicationId, model.UserId, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.AddAppAuthRecord(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("application: %s user: %s create failed, err: %s", model.ApplicationId, model.UserId, errorData.Err.Error())
	}

	return
}

func (svc *AppAuthRecordService) DeleteAppAuthRecord(ctx context.Context, ids []string) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteAppAuthRecord(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}
