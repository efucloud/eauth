package services

import (
	"context"
	"crypto"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ConfigService struct {
	repo repositories.ConfigRepository
}

func (svc *ConfigService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.ConfigRepository{DB: db}
	} else {
		svc.repo = repositories.ConfigRepository{DB: config.DBConnect}
	}
}

func (svc *ConfigService) GetConfigByCode(ctx context.Context, code string) (result dtos.ConfigDetail, errorData common.ErrorData) {
	svc.init(ctx)

	result, errorData = svc.repo.GetConfigByCode(ctx, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getConfig by code: %s failed, err: %s", code, errorData.Err.Error())
	}
	return result, errorData
}

func (svc *ConfigService) GetConfigByID(ctx context.Context, id uint) (result dtos.ConfigDetail, errorData common.ErrorData) {
	svc.init(ctx)

	result, errorData = svc.repo.GetConfigByID(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("getConfig by id: %d failed, err: %s", id, errorData.Err.Error())
	}
	return result, errorData
}
func (svc *ConfigService) ListConfig(ctx context.Context, current, pageSize int, order, query string, queryArgs []interface{}) (results dtos.ConfigDetailList, errorData common.ErrorData) {
	svc.init(ctx)
	results, errorData = svc.repo.ListConfig(ctx, current, pageSize, order, query, queryArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("listConfig  query: [%s] queryArgs: [%+v] failed, err: %s", query, queryArgs, errorData.Err.Error())
	}
	return results, errorData
}
func (svc *ConfigService) UpdateConfig(ctx context.Context, model dtos.ConfigUpdate) (result dtos.ConfigDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("updateConfig: %s failed, err: %s", model.Name, errorData.Err.Error())
		return
	}
	result, errorData = svc.repo.UpdateConfig(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updateConfig: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}
func (svc *ConfigService) AddConfig(ctx context.Context, model dtos.ConfigCreate) (result dtos.ConfigDetail, errorData common.ErrorData) {
	svc.init(ctx)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("createConfig: %s failed, err: %s", model.Name, errorData.Err.Error())

		return
	}
	result, errorData = svc.repo.AddConfig(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("createConfig: %s failed, err: %s", model.Name, errorData.Err.Error())
	}
	return
}

func (svc *ConfigService) DeleteConfig(ctx context.Context, id uint) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteConfig(ctx, id)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteConfig by id: %d failed, err: %s", id, errorData.Err.Error())
	}
	return
}

func (svc *ConfigService) DeleteConfigByCode(ctx context.Context, code string) (errorData common.ErrorData) {
	svc.init(ctx)
	errorData = svc.repo.DeleteConfigByCode(ctx, code)
	if errorData.IsNotNil() {
		config.Logger.Errorf("deleteConfig by code: %s failed, err: %s", code, errorData.Err.Error())
	}
	return
}
func (svc *ConfigService) InitConfig(ctx context.Context) (errorData common.ErrorData) {
	svc.init(ctx)
	var public, private string
	public, private, errorData = svc.GetRsaKeys(ctx)
	if len(public) == 0 || len(private) == 0 {
		public, private = config.GenerateRsaKeys(4096, 50)
		svc.DeleteConfigByCode(ctx, config.JwtPrivateKey)
		svc.DeleteConfigByCode(ctx, config.JwtPublicKey)
		var pubSys, priSys dtos.ConfigCreate
		pubSys.Code = config.JwtPublicKey
		pubSys.ValType = "string"
		pubSys.Name = "OAuth2.0令牌公钥"
		pubSys.Value = public
		pubSys.Description = "OAuth2.0令牌公钥，删除后系统会自动创建，多副本情况下，修改后需重启服务生效"
		svc.AddConfig(ctx, pubSys)
		priSys.Code = config.JwtPrivateKey
		priSys.ValType = "string"
		priSys.Name = "OAuth2.0令牌私钥"
		priSys.Value = private
		priSys.Description = "OAuth2.0令牌私钥，删除后系统会自动创建，多副本情况下，修改后需重启服务生效"
		svc.AddConfig(ctx, priSys)
	}
	// 加载全局证书
	if verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(public)); err == nil {
		keySet := &oidc.StaticKeySet{PublicKeys: []crypto.PublicKey{verifyKey}}
		config.Verifier = oidc.NewVerifier(config.ApplicationConfig.ServerAddress, keySet, &oidc.Config{
			SkipClientIDCheck: true,
			ClientID:          config.ApplicationName,
		})
	}
	return
}

func (svc *ConfigService) GetRsaKeys(ctx context.Context) (public, private string, errorData common.ErrorData) {
	svc.init(ctx)
	var sys dtos.ConfigDetail
	sys, errorData = svc.GetConfigByCode(ctx, config.JwtPrivateKey)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get jwt private key: %s failed, err: %s", config.JwtPrivateKey, errorData.String())
		return
	} else {
		private = sys.Value
	}
	sys, errorData = svc.GetConfigByCode(ctx, config.JwtPublicKey)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get jwt public key: %s failed, err: %s", config.JwtPublicKey, errorData.String())
		return
	} else {
		public = sys.Value
	}
	return
}
