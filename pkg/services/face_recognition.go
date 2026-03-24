package services

import (
	"context"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type FaceRecognitionService struct {
	repo repositories.FaceRecognitionRepository
}

func (svc *FaceRecognitionService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.FaceRecognitionRepository{DB: db}
	} else {
		svc.repo = repositories.FaceRecognitionRepository{DB: config.DBConnect}
	}
}

func (svc *FaceRecognitionService) LoginByFaceIdData(ctx context.Context, loginParam dtos.LoginByFaceIdData) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		user dtos.UserDetail
	)

	userSvc := UserService{}
	user, _ = userSvc.GetUserByUsername(ctx, loginParam.Username)
	if dtos.CheckFaceIdData(user.FaceIdDatas, loginParam.FaceIdData) {
		oauthSvc := OAuthService{}
		return oauthSvc.GenerateTokenResponse(ctx, false, config.ApplicationName, user)
	}
	errorData.Err = fmt.Errorf("auth failed")
	return
}

func (svc *FaceRecognitionService) ChangeStatus(ctx context.Context, model dtos.FaceRecognitionStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	//todo 剔除管理员
	errorData = svc.repo.ChangeStatus(ctx, model)
	return
}
func (svc *FaceRecognitionService) GetUserPersonalFaceRecognitionsByUsername(ctx context.Context, username string) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	userSvc := UserService{}
	user, _ := userSvc.GetUserByUsername(ctx, username)
	return svc.repo.GetUserPersonalFaceRecognitions(ctx, user.ID)
}
func (svc *FaceRecognitionService) GetUserPersonalFaceRecognitions(ctx context.Context, userId string) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	return svc.repo.GetUserPersonalFaceRecognitions(ctx, userId)
}
func (svc *FaceRecognitionService) GetFaceRecognitionById(ctx context.Context, id string) (result dtos.FaceRecognitionDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetFaceRecognitionById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getFaceRecognition by id: %s failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *FaceRecognitionService) GetAllFaceRecognitions(ctx context.Context) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	//获取所有的
	current := 1
	pageSize := 1000
	order := "id desc"
	query := ""
	var queryArgs []interface{}
	var records []*dtos.FaceRecognitionDetail
	for {
		r, e := svc.repo.ListFaceRecognition(ctx, current, pageSize, order, query, queryArgs)
		if e.IsNotNil() {
			break
		}
		if len(r.Data) > 0 {
			records = append(records, r.Data...)
		}
		if int64(len(records)) == r.Total {
			break
		}
		current += 1
	}
	results.Data = records
	results.Total = int64(len(results.Data))
	return results, errorData
}
func (svc *FaceRecognitionService) ListFaceRecognition(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.FaceRecognitionDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListFaceRecognition(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listFaceRecognition query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}

	return results, errorData
}

func (svc *FaceRecognitionService) AddFaceRecognition(ctx context.Context, model dtos.FaceRecognitionCreate) (result dtos.FaceRecognitionDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("FaceRecognition: %s create failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.AddFaceRecognition(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("FaceRecognition: %s create failed, err: %s", model.Name, errorData.Err.Error())
	}

	return
}

func (svc *FaceRecognitionService) DeleteFaceRecognition(ctx context.Context, ids []string) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteFaceRecognition(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	return
}
