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

type UserRepository struct {
	DB *gorm.DB
}

func (repo *UserRepository) GetUserByUsername(ctx context.Context, username string) (result dtos.UserDetail, errorData common.ErrorData) {

	//用户名，手机号码，邮箱，工号登录
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("username = ? OR phone = ? OR email = ? OR job_number = ?",
		username, username, username, username).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found account by username: %s", username)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *UserRepository) SetRole(ctx context.Context, model dtos.UserRole) (result dtos.UserDetail, errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["role"] = model.Role
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id IN (?)", model.Ids).
		UpdateColumns(columns).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return
}
func (repo *UserRepository) ChangeStatusUser(ctx context.Context, model dtos.UserStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["enable"] = model.Enable
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id IN (?)", model.Ids).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return
}

func (repo *UserRepository) ResetUserPassword(ctx context.Context, model dtos.UserResetPassword, password string) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["password_store"] = password
	columns["updated_at"] = time.Now()
	columns["password_strength"] = model.PasswordStrength
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id = ?", model.ID).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return
}
func (repo *UserRepository) GetUsersByIds(ctx context.Context, ids []uint) (results dtos.UserDetailList, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id IN (?)", ids).Find(&results.Data).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	}
	return
}

func (repo *UserRepository) UpdateUserMFa(ctx context.Context, userIds []uint, mfa bool) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id IN (?)", userIds).Update("mfa", mfa).Error
	return
}
func (repo *UserRepository) UpdateUserAvatar(ctx context.Context, userId uint, avatarAddress string) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id = ?", userId).Update("avatar", avatarAddress).Error
	return
}

func (repo *UserRepository) GetUserByPhone(ctx context.Context, phone string) (result dtos.UserDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("phone = ?", phone).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found account by phone: %s", phone)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (result dtos.UserDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("email = ?", email).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found account by email: %s", email)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *UserRepository) GetUserByID(ctx context.Context, id uint) (result dtos.UserDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not found  user by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}
func (repo *UserRepository) ListUser(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.UserTableName)

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
func (repo *UserRepository) AddUser(ctx context.Context, model dtos.UserCreate) (result dtos.UserDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		if len(errorData.MsgCode) == 0 {
			errorData.MsgCode = config.MsgCodeCreateRecordFailed
		}
	}
	return result, errorData
}

func (repo *UserRepository) DeleteUser(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id IN (?)", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return
}
func (repo *UserRepository) UpdateUser(ctx context.Context, model dtos.UserUpdate) (result dtos.UserDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.UserTableName).Where("id = ?", model.ID).Updates(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		config.Logger.Errorf("err: %s", errorData.Err.Error())
		errorData = daos.ParserDatabaseError(ctx, errorData)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}
	return result, errorData
}
