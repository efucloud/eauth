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

type ValidateCodeRepository struct {
	DB *gorm.DB
}

func (repo *ValidateCodeRepository) GetValidateCodeByCode(ctx context.Context, category, code string) (result dtos.ValidateCodeDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ValidateCodeTableName).Where("category = ? AND code = ?", category, code).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not foundValidateCode by code: %s", code)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *ValidateCodeRepository) AddValidateCode(ctx context.Context, model dtos.ValidateCodeCreate) (result dtos.ValidateCodeDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ValidateCodeTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}
func (repo *ValidateCodeRepository) DeleteValidateCode(ctx context.Context) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.ValidateCodeTableName).Where("expired <= ? ", time.Now()).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return errorData
}
