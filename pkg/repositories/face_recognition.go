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

type FaceRecognitionRepository struct {
	DB *gorm.DB
}

func (repo *FaceRecognitionRepository) GetFaceRecognitionById(ctx context.Context, id uint) (result dtos.FaceRecognitionDetail, errorData common.ErrorData) {

	errorData.Err = repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Where("id = ?", id).Find(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeGetRecordFailed
	} else if result.ID == 0 {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.ResponseCode = http.StatusNotFound
		errorData.Err = fmt.Errorf("not foundFaceRecognition by id: %d", id)
		config.Logger.Error(errorData.Err)
	}
	return
}

func (repo *FaceRecognitionRepository) ListFaceRecognition(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	db := repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName)
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
func (repo *FaceRecognitionRepository) AddFaceRecognition(ctx context.Context, model dtos.FaceRecognitionCreate) (result dtos.FaceRecognitionDetail, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Create(&model).Scan(&result).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeCreateRecordFailed
	}
	return result, errorData
}
func (repo *FaceRecognitionRepository) DeleteFaceRecognition(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Where("id IN (?) ", ids).Delete(nil).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeDeleteRecordFailed
	}
	return errorData
}

func (repo *FaceRecognitionRepository) ChangeStatus(ctx context.Context, model dtos.FaceRecognitionStatus) (errorData common.ErrorData) {
	columns := make(map[string]interface{})
	columns["enable"] = model.Enable
	columns["updated_at"] = model.UpdatedAt
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Where("id IN (?)", model.Ids).
		UpdateColumns(columns).Error
	if errorData.IsNotNil() {
		errorData = daos.ParserDatabaseError(ctx, errorData)
		config.Logger.Error(errorData.Err)
		errorData.ResponseCode = http.StatusInternalServerError
		errorData.MsgCode = config.MsgCodeUpdateRecordFailed
	}

	return
}

func (repo *FaceRecognitionRepository) GetUserPersonalFaceRecognitions(ctx context.Context, userId uint) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	errorData.Err = repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Where("user_id = ?", userId).Find(&results.Data).Error
	results.Total = int64(len(results.Data))
	return
}
