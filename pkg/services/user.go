package services

import (
	"context"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
	"net/http"

	"github.com/efucloud/common"
)

type UserService struct {
	repo repositories.UserRepository
}

func (svc *UserService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.UserRepository{DB: db}
	} else {
		svc.repo = repositories.UserRepository{DB: config.DBConnect}
	}
}
func (svc *UserService) UpdateUserMFa(ctx context.Context, userIds []uint, mfa bool) (errorData common.ErrorData) {
	svc.init(ctx)
	return svc.repo.UpdateUserMFa(ctx, userIds, mfa)
}
func (svc *UserService) UpdateUserAvatar(ctx context.Context, userId uint, avatarAddress string) (errorData common.ErrorData) {
	svc.init(ctx)
	return svc.repo.UpdateUserAvatar(ctx, userId, avatarAddress)
}
func (svc *UserService) SetRole(ctx context.Context, model dtos.UserRole) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	result, errorData = svc.repo.SetRole(ctx, model)
	return
}

func (svc *UserService) ResetPassword(ctx context.Context, model dtos.UserResetPassword) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	password := ""
	password, errorData.Err = common.GeneratePassword(model.Password, config.PasswordSalt)
	if errorData.IsNotNil() {
		return
	}
	errorData = svc.repo.ResetUserPassword(ctx, model, password)
	return
}
func (svc *UserService) ChangeStatusUser(ctx context.Context, model dtos.UserStatus) (errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	//todo 剔除管理员
	errorData = svc.repo.ChangeStatusUser(ctx, model)
	return
}
func (svc *UserService) GetUsersByIds(ctx context.Context, ids []uint) (results dtos.UserDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.GetUsersByIds(ctx, ids)
	return results, errorData
}
func (svc *UserService) GetShortUserByID(ctx context.Context, id uint) (result dtos.ShortUser, errorData common.ErrorData) {
	svc.init(ctx)
	var user dtos.UserDetail
	user, errorData = svc.repo.GetUserByID(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUser by id: %d failed, err: %s", id, errorData.Err.Error())
		return
	}
	copyByJSON(user, &result)
	return result, errorData
}

func (svc *UserService) GetUserByPhone(ctx context.Context, phone string) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserByPhone(ctx, phone)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUser by phone: %s failed, err: %s", phone, errorData.Err.Error())
		return
	}
	return result, errorData
}

func (svc *UserService) GetUserByEmail(ctx context.Context, email string) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserByEmail(ctx, email)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUser by email: %s failed, err: %s", email, errorData.Err.Error())
		return
	}

	return result, errorData
}
func (svc *UserService) GetUserByID(ctx context.Context, id uint) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserByID(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUser by id: %d failed, err: %s", id, errorData.Err.Error())
		return
	}
	faceSvc := FaceRecognitionService{}
	faceDatas, _ := faceSvc.GetUserPersonalFaceRecognitions(ctx, result.ID)
	for _, item := range faceDatas.Data {
		result.FaceIdDatas = append(result.FaceIdDatas, item.FaceIdData)
	}
	return result, errorData
}

func (svc *UserService) GetUserByUsername(ctx context.Context, name string) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	result, errorData = svc.repo.GetUserByUsername(ctx, name)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getUser by username: %s failed, err: %s", name, errorData.Err.Error())
	}

	faceSvc := FaceRecognitionService{}
	faceDatas, _ := faceSvc.GetUserPersonalFaceRecognitions(ctx, result.ID)
	for _, item := range faceDatas.Data {
		result.FaceIdDatas = append(result.FaceIdDatas, item.FaceIdData)
	}
	return result, errorData
}
func (svc *UserService) ListUser(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.UserDetailList, err common.ErrorData) {
	svc.init(ctx)
	results, err = svc.repo.ListUser(ctx, current, pageSize, order, query, queryArgs)
	if !err.IsNil() {
		config.Logger.Errorf("listUser  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, err.Err.Error())
	}
	return
}
func (svc *UserService) UpdateUser(ctx context.Context, model dtos.UserUpdate) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if !errorData.IsNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("User: %s create failed, errorData: %s", model.Username, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateUser(ctx, model)
	if !errorData.IsNil() {
		config.Logger.Errorf("updateUser: %s failed, errorData: %s", model.Username, errorData.Err.Error())
		errorData.ResponseCode = http.StatusInternalServerError
	}
	return
}

func (svc *UserService) AddUser(ctx context.Context, model dtos.UserCreate) (result dtos.UserDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if !errorData.IsNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("User: %s create failed, errorData: %s", model.Username, errorData.Err.Error())
		return
	}

	if len(model.Password) > 0 {
		model.PasswordStore, errorData.Err = common.GeneratePassword(model.Password, config.PasswordSalt)
		if errorData.IsNotNil() {
			config.Logger.Errorf("User: %s create failed, err: %s", model.Username, errorData.Err.Error())
			return
		}
	}
	result, errorData = svc.repo.AddUser(ctx, model)

	return
}
func (svc *UserService) DeleteUser(ctx context.Context, ids []uint) (errorData common.ErrorData) {
	svc.init(ctx)
	tx := config.DBConnect.Begin()
	defer func() {
		if errorData.IsNotNil() {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	ctx = context.WithValue(ctx, config.ContextDBTx, tx)

	//todo 内建的管理员名不可删除
	errorData = svc.repo.DeleteUser(ctx, ids)
	if !errorData.IsNil() {
		config.Logger.Errorf("deleteUser by ids: %v failed, errorData: %s", ids, errorData.Err.Error())
	}
	authProfileSvc := UserAuthProfileService{}
	authProfileSvc.DeleteUserAuthProfileByUserIds(ctx, ids)
	return
}
