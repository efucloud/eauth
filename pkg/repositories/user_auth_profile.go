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

type AuthProfileRepository struct {
	DB *gorm.DB
}

func (repo *AuthProfileRepository) GetUserAuthProfilesByUserId(ctx context.Context, userId uint) (results dtos.UserAuthProfileDetailList, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Where("user_id = ?", userId).Find(&results.Data).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	return
}
func (repo *AuthProfileRepository) GetUserAuthProfileByProviderAndId(ctx context.Context, provider, providerUserId string) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Debug().Where("provider = ? AND login_id = ?", provider, providerUserId).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	return
}
func (repo *AuthProfileRepository) GetUserAuthProfileByID(ctx context.Context, id uint) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found record by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *AuthProfileRepository) ListUserAuthProfile(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserAuthProfileDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName)

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
func (repo *AuthProfileRepository) AddUserAuthProfile(ctx context.Context, model dtos.UserAuthProfileCreate) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *AuthProfileRepository) DeleteUserAuthProfileByUserIds(ctx context.Context, userIds []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Where("user_id IN (?) ", userIds).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}

func (repo *AuthProfileRepository) DeleteUserAuthProfile(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
func (repo *AuthProfileRepository) UpdateUserAuthProfile(ctx context.Context, model dtos.UserAuthProfileUpdate) (result dtos.UserAuthProfileDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}
func (repo *AuthProfileRepository) ChangeStatus(ctx context.Context, model dtos.UserAuthProfileStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["enable"] = model.Enable
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Where("id IN (?)", model.Ids).
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
