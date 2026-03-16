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

type AuthProfileTempRepository struct {
	DB *gorm.DB
}

func (repo *AuthProfileTempRepository) GetUserAuthProfileTempByCode(ctx context.Context, code string) (result dtos.UserAuthProfileTempDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTempTableName).Where("code = ?", code).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found record by code: %s", code)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *AuthProfileTempRepository) ListUserAuthProfileTemp(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserAuthProfileTempDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTempTableName)

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
func (repo *AuthProfileTempRepository) AddUserAuthProfileTemp(ctx context.Context, model dtos.UserAuthProfileTempCreate) (result dtos.UserAuthProfileTempDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTempTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *AuthProfileTempRepository) DeleteUserAuthProfileTemp(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTempTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
