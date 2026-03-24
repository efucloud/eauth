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

type AppAuthRecordRepository struct {
	DB *gorm.DB
}

func (repo *AppAuthRecordRepository) GetAppAuthRecordByCode(ctx context.Context, code string) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Where("code = ?", code).Find(&result).Error
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

func (repo *AppAuthRecordRepository) GetAppAuthRecordById(ctx context.Context, id string) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found application by id: %s", id)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *AppAuthRecordRepository) ListAppAuthRecord(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.AppAuthRecordDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName)
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
func (repo *AppAuthRecordRepository) AddAppAuthRecord(ctx context.Context, model dtos.AppAuthRecordCreate) (result dtos.AppAuthRecordDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}
func (repo *AppAuthRecordRepository) DeleteAppAuthRecord(ctx context.Context, ids []string) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return errorData
}
