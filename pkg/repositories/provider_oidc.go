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

type ProviderOidcRepository struct {
	DB *gorm.DB
}

func (repo *ProviderOidcRepository) GetProviderOidcByCategory(ctx context.Context, category string) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Where("category = ?", category).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found  provider oidc by category: %s", category)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *ProviderOidcRepository) GetProviderOidcById(ctx context.Context, id string) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found  provider oidc by id: %s", id)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *ProviderOidcRepository) ListProviderOidc(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderOidcDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName)

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
func (repo *ProviderOidcRepository) AddProviderOidc(ctx context.Context, model dtos.ProviderOidcCreate) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Save(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *ProviderOidcRepository) DeleteProviderOidc(ctx context.Context, ids []string) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Where("id IN (?)", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
func (repo *ProviderOidcRepository) UpdateProviderOidc(ctx context.Context, model dtos.ProviderOidcUpdate) (result dtos.ProviderOidcDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}

func (repo *ProviderOidcRepository) ChangeStatus(ctx context.Context, model dtos.ProviderOidcStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["enable"] = model.Enable
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderOidcTableName).Where("id IN (?)", model.Ids).
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
