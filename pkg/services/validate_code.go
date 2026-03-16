package services

import (
	"context"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ValidateCodeService struct {
	repo repositories.ValidateCodeRepository
}

func (svc *ValidateCodeService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ValidateCodeRepository{DB: db}
	} else {
		svc.repo = repositories.ValidateCodeRepository{DB: config.DBConnect}
	}
}
func (svc *ValidateCodeService) GetValidateCodeByCode(ctx context.Context, category, code string) (result dtos.ValidateCodeDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetValidateCodeByCode(ctx, category, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getValidateCode by code: %s failed, err: %s", code, errorData.Err.Error())
	}
	if !result.Expired.After(time.Now()) {
		errorData.Err = fmt.Errorf("action expired")
		errorData.MsgCode = config.MsgCodeActionExpired
		errorData.ResponseCode = http.StatusBadRequest
	}
	return result, errorData
}
func (svc *ValidateCodeService) AddValidateCode(ctx context.Context, model dtos.ValidateCodeCreate) (result dtos.ValidateCodeDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("ValidateCode: %s create failed, err: %s", model.Code, errorData.Err.Error())
		return
	}
	return svc.repo.AddValidateCode(ctx, model)
}

func (svc *ValidateCodeService) DeleteValidateCode(ctx context.Context) (errorData common.ErrorData) {
	svc.init(ctx)
	_ = svc.repo.DeleteValidateCode(ctx)
	return
}
