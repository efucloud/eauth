package repositories

import (
	"context"
	"fmt"
	"net/http"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/daos"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	DB *gorm.DB
}

func (repo *ApplicationRepository) GetApplicationByClientId(ctx context.Context, clientId string) (result dtos.ApplicationDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("client_id = ?", clientId).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found application by clientId: %s", clientId)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ApplicationRepository) GetApplicationByCode(ctx context.Context, code string) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("code = ?", code).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found application by code: %s", code)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ApplicationRepository) GetApplicationById(ctx context.Context, id string) (result dtos.ApplicationDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found application by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ApplicationRepository) ListApplication(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ApplicationDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.ApplicationTableName)
	if len(query) > 0 {
		db = db.Where(query, queryArgs...)
	}
	errorData.Err = db.Count(&results.Total).Error
	if errorData.IsNil() {

		if len(order) > 0 {
			db = db.Order(order)
		}
		errorData.Err = db.Offset((current - 1) * pageSize).Limit(pageSize).Find(&results.Data).Error
	} else {
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	return
}
func (repo *ApplicationRepository) AddApplication(ctx context.Context, model dtos.ApplicationCreate) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}
func (repo *ApplicationRepository) DeleteApplication(ctx context.Context, ids []string) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return errorData
}
func (repo *ApplicationRepository) UpdateApplication(ctx context.Context, model dtos.ApplicationUpdate) (result dtos.ApplicationDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("id = ?", model.ID).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}

func (repo *ApplicationRepository) ChangeStatus(ctx context.Context, model dtos.ApplicationStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["enable"] = model.Enable
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ApplicationTableName).Where("id IN (?)", model.Ids).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	} else {
		//todo 多条更新
	}
	return
}
