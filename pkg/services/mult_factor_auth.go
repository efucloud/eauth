package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"image"
	"image/png"
	"time"
)

type MultiFactorAuthService struct {
	repo repositories.MultiFactorAuthRepository
}

func (svc *MultiFactorAuthService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.MultiFactorAuthRepository{DB: db}
	} else {
		svc.repo = repositories.MultiFactorAuthRepository{DB: config.DBConnect}
	}
}

func (svc *MultiFactorAuthService) ChangeStatus(ctx context.Context, userId uint, model dtos.MultiFactorAuthStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("MultiFactorAuth: %d create failed, err: %s", model.Id, errorData.Err.Error())
		return
	}
	errorData = svc.repo.ChangeStatus(ctx, userId, model)
	return
}

func (svc *MultiFactorAuthService) GetUserMultiFactorAuthByUserId(ctx context.Context, userId uint) (results dtos.MultiFactorAuthDetail) {
	svc.init(ctx)
	return svc.repo.GetUserMultiFactorAuth(ctx, userId)
}
func (svc *MultiFactorAuthService) GetMultiFactorAuthById(ctx context.Context, id uint) (result dtos.MultiFactorAuthDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetMultiFactorAuthById(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getMultiFactorAuth by id: %d failed, err: %s", id, errorData.Err.Error())
	}

	return result, errorData
}

func (svc *MultiFactorAuthService) ListMultiFactorAuth(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.MultiFactorAuthDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListMultiFactorAuth(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listMultiFactorAuth query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	for i, _ := range results.Data {
		results.Data[i].Secret = ""
		results.Data[i].Image = ""
	}
	return results, errorData
}
func (svc *MultiFactorAuthService) BoundUserMultiFactorAuth(ctx context.Context, userId uint, code string) (model dtos.MultiFactorAuthDetail) {
	model = svc.GetUserMultiFactorAuthByUserId(ctx, userId)
	if model.ID == 0 {
		return
	}
	if totp.Validate(code, model.Secret) {
		if model.Status == "unbound" {
			var status dtos.MultiFactorAuthStatus
			status.Id = model.ID
			status.UpdatedAt = time.Now()
			status.Status = "bound"
			svc.ChangeStatus(ctx, model.UserId, status)
			userSvc := UserService{}
			userSvc.UpdateUserMFa(ctx, []uint{model.UserId}, true)
		}
	}
	return svc.GetUserMultiFactorAuthByUserId(ctx, userId)

}
func (svc *MultiFactorAuthService) ResetUserMultiFactorAuth(ctx context.Context, userId uint, code string) (result dtos.MultiFactorAuthDetail, errorData common.ErrorData) {
	old := svc.GetUserMultiFactorAuthByUserId(ctx, userId)
	if totp.Validate(code, old.Secret) {
		svc.DeleteMultiFactorAuth(ctx, []uint{old.ID})
		userSvc := UserService{}
		userSvc.UpdateUserMFa(ctx, []uint{userId}, false)
		return svc.AddMultiFactorAuth(ctx, userId)
	}
	errorData.Err = fmt.Errorf(" MFA code is invalid")
	return
}
func (svc *MultiFactorAuthService) AddMultiFactorAuth(ctx context.Context, userId uint) (result dtos.MultiFactorAuthDetail, errorData common.ErrorData) {
	svc.init(ctx)
	var (
		key       *otp.Key
		img       image.Image
		buf       bytes.Buffer
		mfaCreate dtos.MultiFactorAuthCreate
	)
	key, _ = totp.Generate(totp.GenerateOpts{
		Issuer:      config.ApplicationName,
		AccountName: fmt.Sprintf("%d", userId),
	})
	mfaCreate.Secret = key.Secret()
	mfaCreate.UserId = userId
	img, _ = key.Image(200, 200)
	_ = png.Encode(&buf, img)
	mfaCreate.Image = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	result, errorData = svc.repo.AddMultiFactorAuth(ctx, mfaCreate)
	if errorData.IsNotNil() {
		config.Logger.Errorf("MultiFactorAuth for user: %d create failed, err: %s", userId, errorData.Err.Error())
	}

	return
}
func (svc *MultiFactorAuthService) DeleteMultiFactorAuthByUserIds(ctx context.Context, userIds []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	return svc.repo.DeleteMultiFactorAuthByUserIds(ctx, userIds)
}
func (svc *MultiFactorAuthService) DeleteMultiFactorAuth(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	tx := config.DBConnect.Begin()
	defer func() {
		if errorData.IsNotNil() {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	svc.init(ctx)
	records, _ := svc.ListMultiFactorAuth(ctx, 1, 1000, "", "id IN (?)", []interface{}{ids})
	var userIds []uint
	for _, item := range records.Data {
		userIds = append(userIds, item.UserId)
	}
	errorData = svc.repo.DeleteMultiFactorAuth(ctx, ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete by ids: %v failed, err: %s", ids, errorData.Err.Error())
	}
	userSvc := UserService{}
	userSvc.UpdateUserMFa(ctx, userIds, false)
	return
}
