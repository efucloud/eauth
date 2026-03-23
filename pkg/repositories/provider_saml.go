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

type ProviderSamlRepository struct {
	DB *gorm.DB
}

func (repo *ProviderSamlRepository) GetProviderSamlByCategory(ctx context.Context, category string) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Where("category = ?", category).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found provider saml by category: %s", category)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ProviderSamlRepository) GetProviderSamlById(ctx context.Context, id uint) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found provider saml by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ProviderSamlRepository) ListProviderSaml(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderSamlDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName)

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

func (repo *ProviderSamlRepository) AddProviderSaml(ctx context.Context, model dtos.ProviderSamlCreate) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *ProviderSamlRepository) DeleteProviderSaml(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Where("id IN (?)", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}

func (repo *ProviderSamlRepository) UpdateProviderSaml(ctx context.Context, model dtos.ProviderSamlUpdate) (result dtos.ProviderSamlDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}

func (repo *ProviderSamlRepository) ChangeStatus(ctx context.Context, model dtos.ProviderSamlStatus) (errorData common.ErrorData) {
	columns := map[string]interface{}{
		"enable":     model.Enable,
		"updated_at": model.UpdatedAt,
	}
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderSamlTableName).Where("id IN (?)", model.Ids).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return
}
