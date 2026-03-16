package filters

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
)

func GetUserInfo(req *restful.Request) (userId uint) {
	var (
		errorData common.ErrorData
	)
	ctx := context.Background()
	if req.Attribute(config.RequestContext) != nil {
		ctx = req.Attribute(config.RequestContext).(context.Context)
	}

	token := GetRequestToken(config.AuthHeader, req)
	if len(token) > 0 {
		var (
			err     error
			idToken *oidc.IDToken
			claims  dtos.UserClaims
		)
		idToken, errorData.Err = config.Verifier.Verify(ctx, token)
		if errorData.IsNotNil() {
			return
		}
		err = idToken.Claims(&claims)
		if err != nil {
			return
		}
		userId = claims.Id
	}
	return
}
func Auth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	var (
		errorData common.ErrorData
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.Background()
	if req.Attribute(config.RequestContext) != nil {
		ctx = req.Attribute(config.RequestContext).(context.Context)
	}
	userId := GetUserInfo(req)
	if userId == 0 {
		errorData.Lang = lang
		errorData.Err = fmt.Errorf("未登录或者认证信息无效")
		errorData.ResponseCode = http.StatusUnauthorized
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	req.SetAttribute(config.RequestUserID, userId)
	chain.ProcessFilter(req, resp)
}

func Permission(roles []string) func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		lang := common.GetLanguageFromReq(req, config.RequestLanguage)
		ctx := context.Background()
		var errorData common.ErrorData

		userId := req.Attribute(config.RequestUserID)
		if userId == nil {
			errorData.Lang = lang
			errorData.ResponseCode = http.StatusForbidden
			errorData.MsgCode = config.MsgCodeCurrentActionIsForbidden
			common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
			return
		}
		userSvc := services.UserService{}
		user, _ := userSvc.GetUserByID(ctx, userId.(uint))
		if !common.StringKeyInArray(user.Role, roles) {
			errorData.Lang = lang
			errorData.ResponseCode = http.StatusForbidden
			errorData.MsgCode = config.MsgCodeCurrentActionIsForbidden
			common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
			return
		}
		chain.ProcessFilter(req, resp)
	}
}
