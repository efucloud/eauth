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

type MultiFactorAuthRepository struct {
	DB *gorm.DB
}

func (repo *MultiFactorAuthRepository) GetMultiFactorAuthById(ctx context.Context, id uint) (result dtos.MultiFactorAuthDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not foundMultiFactorAuth by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *MultiFactorAuthRepository) ListMultiFactorAuth(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.MultiFactorAuthDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName)
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
func (repo *MultiFactorAuthRepository) AddMultiFactorAuth(ctx context.Context, model dtos.MultiFactorAuthCreate) (result dtos.MultiFactorAuthDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *MultiFactorAuthRepository) DeleteMultiFactorAuthByUserIds(ctx context.Context, userIds []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Where("user_id IN (?) ", userIds).Delete(nil).Error

	return errorData
}

func (repo *MultiFactorAuthRepository) DeleteMultiFactorAuth(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return errorData
}

func (repo *MultiFactorAuthRepository) ChangeStatus(ctx context.Context, userId uint, model dtos.MultiFactorAuthStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["status"] = model.Status
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Where("id = ? AND user_id = ?", model.Id, userId).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return
}

func (repo *MultiFactorAuthRepository) GetUserMultiFactorAuth(ctx context.Context, userId uint) (result dtos.MultiFactorAuthDetail) {
	repo.DB.WithContext(ctx).Table(daos.MultiFactorAuthTableName).Where("user_id = ?", userId).Find(&result)
	return
}
