package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/daos"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
	DB *gorm.DB
}

func (repo *UserTokenRepository) GetUserTokensByreRefreshToken(ctx context.Context, refreshToken string) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("refresh_token = ?", refreshToken).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	return
}
func (repo *UserTokenRepository) GetUserTokenDetailById(ctx context.Context, id uint) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found application by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *UserTokenRepository) GetUserTokensByUserId(ctx context.Context, userId uint) (results dtos.UserTokenDetailList, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("user_id = ?", userId).Find(&results.Data).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	results.Total = int64(len(results.Data))
	return
}
func (repo *UserTokenRepository) GetUserTokenByID(ctx context.Context, id uint) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("id = ?", id).Find(&result).Error
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
func (repo *UserTokenRepository) ListUserToken(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserTokenDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.UserTokenTableName)

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
func (repo *UserTokenRepository) AddUserToken(ctx context.Context, model dtos.UserTokenCreate) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}

func (repo *UserTokenRepository) DeleteUserToken(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("id IN (?)", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
func (repo *UserTokenRepository) UpdateUserToken(ctx context.Context, model dtos.UserTokenUpdate) (result dtos.UserTokenDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Save(&model).Find(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}
func (repo *UserTokenRepository) DeleteExpireRecord(ctx context.Context) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTokenTableName).Where("expired <= ? ", time.Now().Unix()).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
