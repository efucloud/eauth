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

type ProviderLdapRepository struct {
	DB *gorm.DB
}

func (repo *ProviderLdapRepository) GetProviderLdapById(ctx context.Context, id string) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderLdapTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if len(result.ID) == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found system provider ldap by id: %s", id)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *ProviderLdapRepository) ListProviderLdap(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ProviderLdapDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.ProviderLdapTableName)

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
func (repo *ProviderLdapRepository) AddProviderLdap(ctx context.Context, model dtos.ProviderLdapCreate) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderLdapTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *ProviderLdapRepository) DeleteProviderLdap(ctx context.Context, ids []string) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderLdapTableName).Where("id IN (?)", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
func (repo *ProviderLdapRepository) UpdateProviderLdap(ctx context.Context, model dtos.ProviderLdapUpdate) (result dtos.ProviderLdapDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ProviderLdapTableName).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}
