package system

import (
	"context"
	"errors"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-jose/go-jose/v4"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

type OAuthResource struct {
	Svc services.OAuthService
}

func (r OAuthResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "oauth"
	apiInfo.Description = "OAuth"
	apiExtend := ""
	common.RegisterApiInfo(apiInfo)
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/.well-known/openid-configuration").
		Doc("openid configuration").
		Notes("openid configuration").
		To(r.oidcCfg).
		Returns(http.StatusOK, "成功", dtos.OpenIDConfiguration{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemOidcConfig"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/discovery/jwks").
		Doc("openid jwks").
		Notes("openid jwks").
		To(r.publicKeys).
		Returns(http.StatusOK, "成功", jose.JSONWebKeySet{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemPublicKeys"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/login/username").
		Doc("用户名/手机号码/邮箱 + 密码登录").
		Notes("用户名/手机号码/邮箱 + 密码登录").
		To(r.loginUsername).
		Reads(dtos.LoginByUsername{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemLoginByUsername"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/login/oidc").
		Doc("OIDC方式登录").
		Notes("OIDC回调后前端给到后端的Code接口，用于换取第三方的token并获取用户信息，若用户在系统不存在，"+
			"则根据组织是否允许自动注册来决定是否自动创建用户信息，若第一次是通过第三方登录，需要先设置密码，"+
			"若组织设置了MFA则需要再次输入验证码，若用户没有绑定过验证器，则返回验证器的二维码和密钥").
		To(r.loginByOIDC).
		Reads(dtos.LoginByOIDC{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemLoginByOidc"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/login/saml").
		Doc("SAML方式登录").
		Notes("SAML回调后前端给到后端的SAMLResponse接口，用于解析SAML断言并登录系统").
		To(r.loginBySAML).
		Reads(dtos.LoginBySAML{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemLoginBySaml"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/register/oidc").
		Doc("OIDC认证新用户注册接口").
		Notes("OIDC认证新用户注册接口，若系统用户注册且当前认证方式需在系统中创建新用户时").
		To(r.registerByOIDC).
		Reads(dtos.RegisterByOIDC{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "registerByOIDC"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/mfa/validate").
		Doc("MFA认证").
		Notes("MFA认证，使用MFA验证码获取token").
		To(r.mfaValidate).
		Reads(dtos.MfaCode{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "mfaValidate"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/third/auth/methods").
		Doc("列出支持的第三方登录方式").
		Notes("列出支持的第三方登录方式包括OIDC，LDAP，未来包括cas等").
		To(r.thirdAuthMethods).
		Returns(http.StatusOK, "", dtos.ThirdAuthMethod{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "thirdAuthMethods"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/login/ldap").
		Doc("LDAP方式登录").
		Notes("LDAP登录，认证成功节后，若用户在系统不存在，"+
			"则根据组织是否允许自动注册来决定是否自动创建用户信息，若第一次是通过第三方登录，需要先设置密码，"+
			"若组织设置了MFA则需要再次输入验证码，若用户没有绑定过验证器，则返回验证器的二维码和密钥").
		To(r.loginByLDAP).
		Reads(dtos.LoginByLDAP{}).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemLoginByLdap"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/token").
		Doc("获取Token").
		Notes("获取Token，支持四种token获取方式").
		Param(ws.QueryParameter("client_id", "client ID")).
		Param(ws.QueryParameter("client_secret", "client TotpSecret，使用场景:grant_type 为"+
			config.GrantTypePassword+";"+config.GrantTypeClientCredentials)).
		Param(ws.QueryParameter("grant_type", "Grant Type").
			PossibleValues([]string{
				config.GrantTypePassword,
				config.GrantTypeAuthorizationCode,
				config.GrantTypeToken,
				config.GrantTypeRefreshToken,
				config.GrantTypeClientCredentials})).
		Param(ws.QueryParameter("code", "StateCode")).
		Param(ws.QueryParameter("redirect_uri", "Redirect URI")).
		Param(ws.QueryParameter("code_verifier", "PKCE code verifier")).
		Param(ws.QueryParameter("username", "Username")).
		Param(ws.QueryParameter("password", "Password")).
		Consumes(restful.MIME_JSON, config.RequestForm).
		To(r.token).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemFetchTokenWithGet"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/token").
		Doc("获取Token").
		Notes("获取Token，支持四种token获取方式").
		Param(ws.QueryParameter("client_id", "client ID")).
		Param(ws.QueryParameter("client_secret", "client TotpSecret，使用场景:grant_type 为"+
			config.GrantTypePassword+";"+config.GrantTypeClientCredentials)).
		Param(ws.QueryParameter("grant_type", "Grant Type").
			PossibleValues([]string{
				config.GrantTypePassword,
				config.GrantTypeAuthorizationCode,
				config.GrantTypeToken,
				config.GrantTypeRefreshToken,
				config.GrantTypeClientCredentials})).
		Param(ws.QueryParameter("code", "StateCode")).
		Param(ws.QueryParameter("redirect_uri", "Redirect URI")).
		Param(ws.QueryParameter("code_verifier", "PKCE code verifier")).
		Param(ws.QueryParameter("username", "Username")).
		Param(ws.QueryParameter("password", "Password")).
		Consumes(restful.MIME_JSON, config.RequestForm).
		To(r.token).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemFetchTokenWithPost"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/refresh_token").
		Doc("刷新Token").
		Notes("刷新Token").
		Param(ws.QueryParameter("refresh_token", "Refresh Token")).
		Consumes(restful.MIME_JSON).
		To(r.refreshToken).
		Returns(http.StatusOK, "成功", dtos.AccessTokenResponse{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemRefreshTokenWithGet"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/oidc/code").
		Doc("获取code").
		Notes("客户端获取code后，给后端获取用户信息").
		Consumes(restful.MIME_JSON).
		Reads(dtos.OidcCodeRequest{}).
		To(r.getOidcCode).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Returns(http.StatusOK, "成功", dtos.OidcCodeResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getOidcCode"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/refresh_token").
		Doc("刷新Token").
		Notes("刷新Token").
		Param(ws.QueryParameter("refresh_token", "Refresh Token")).
		Consumes(restful.MIME_JSON).
		To(r.refreshToken).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemRefreshTokenWithPost"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/introspect").
		Doc("令牌内省").
		Notes("校验Token是否合法并查询信息，比如，access_token是否还有效，谁颁发的，颁发给谁的，scope又哪些等等的元数据信息").
		Param(ws.HeaderParameter(config.AuthHeader, "Basic Token，clientId:clientSecret base64").Required(true)).
		Param(ws.FormParameter("token", "可以是access_token或者refresh_token").Required(true)).
		Param(ws.FormParameter("token_type_hint", "表示token的类型").PossibleValues([]string{"access_token", "refresh_token"})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		To(r.introspect). //https://blog.51cto.com/demo007x/6193913
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemIntrospect"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/revoke").
		Doc("Token撤销").
		Notes("校验Token是否合法并查询信息，比如，access_token是否还有效，谁颁发的，颁发给谁的，scope又哪些等等的元数据信息").
		Param(ws.HeaderParameter(config.AuthHeader, "Basic Token，clientId:clientSecret base64").Required(true)).
		Param(ws.FormParameter("token", "可以是access_token或者refresh_token").Required(true)).
		Param(ws.FormParameter("token_type_hint", "表示token的类型").PossibleValues([]string{"access_token", "refresh_token"})).
		Param(ws.FormParameter("client_id", "client ID")).
		Param(ws.FormParameter("client_secret", "client TotpSecret，使用场景:grant_type 为"+
			config.GrantTypePassword+";"+config.GrantTypeClientCredentials)).
		Consumes(config.RequestForm).
		To(r.revoke).
		Consumes(config.RequestForm).
		Returns(http.StatusOK, "成功", nil).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemRevoke"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/userinfo").
		Doc("获取用户信息").
		Notes("获取用户信息").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.userinfo).
		Returns(http.StatusOK, "请求成功", dtos.AuthedUserInfo{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemGetUserinfo"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/userinfo-no-error").
		Doc("获取用户信息").
		Notes("获取用户信息，在没有认证的情况下返回空内容").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.userinfoNoError).
		Returns(http.StatusOK, "请求成功", dtos.AuthedUserInfo{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemGetUserinfoNoError"))

	ws.Route(ws.PUT(config.V1Prefix+apiExtend+"/self/info").
		Doc("更新个人信息").
		Notes("用户更新个人信息").
		To(r.selfInfo).
		Reads(dtos.UserUpdate{}).
		Consumes(restful.MIME_JSON).
		Returns(http.StatusOK, "Success", dtos.AuthedUserInfo{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "updateSelfInfo"))

}

func (r OAuthResource) userinfoNoError(req *restful.Request, resp *restful.Response) {
	var userinfo dtos.AuthedUserInfo
	userId := filters.GetUserInfo(req)
	if userId > 0 {
		userinfo, _ = r.Svc.Userinfo(context.Background(), userId)
	}
	common.ResponseSuccess(resp, userinfo)
}

func (r OAuthResource) mfaValidate(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.AccessTokenResponse
		model     dtos.MfaCode
	)

	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.MfaValidate(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}
func (r OAuthResource) registerByOIDC(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.AccessTokenResponse
		model     dtos.RegisterByOIDC
	)

	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.RegisterByOIDC(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r OAuthResource) getOidcCode(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.OidcCodeResponse
		model     dtos.OidcCodeRequest
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	contentType := req.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, config.RequestForm) {
		errorData.Err = req.Request.ParseForm()
		if errorData.IsNil() {
			fillOidcCodeRequestFromForm(&model, req.Request.PostForm)
		}
	} else {
		errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	}
	if len(model.GetCodeChallenge()) == 0 {
		model.CodeChallenge = firstNotEmpty(req.QueryParameter("codeChallenge"), req.QueryParameter("code_challenge"))
	}
	if len(model.GetCodeChallengeMethod()) == 0 {
		model.CodeChallengeMethod = firstNotEmpty(req.QueryParameter("codeChallengeMethod"), req.QueryParameter("code_challenge_method"))
	}
	if len(model.GetCodeChallenge()) == 0 || len(model.GetCodeChallengeMethod()) == 0 {
		fillOidcCodeRequestFromReferer(&model, req.Request.Header.Get("Referer"))
	}
	if len(model.Nonce) == 0 {
		model.Nonce = req.QueryParameter("nonce")
	}
	if len(model.Nonce) == 0 {
		fillOidcCodeRequestFromReferer(&model, req.Request.Header.Get("Referer"))
	}
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	userId := req.Attribute(config.RequestUserID)
	if userId == nil {
		errorData.MsgCode = config.MsgCodeUserInfoIsEmpty
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.ApplicationAuthCode(ctx, userId.(uint), model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("add account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}

func fillOidcCodeRequestFromForm(model *dtos.OidcCodeRequest, form map[string][]string) {
	model.ClientId = formFirstValue(form, "clientId", "client_id")
	model.RedirectUri = formFirstValue(form, "redirectUri", "redirect_uri")
	model.State = formFirstValue(form, "state")
	model.ResponseType = formFirstValue(form, "responseType", "response_type")
	model.CodeChallenge = formFirstValue(form, "codeChallenge")
	model.CodeChallengeMethod = formFirstValue(form, "codeChallengeMethod")
	model.CodeChallengeStd = formFirstValue(form, "code_challenge")
	model.CodeChallengeMethodStd = formFirstValue(form, "code_challenge_method")
	model.Nonce = formFirstValue(form, "nonce")
}

func fillOidcCodeRequestFromReferer(model *dtos.OidcCodeRequest, referer string) {
	if len(referer) == 0 {
		return
	}
	parsed, err := url.Parse(referer)
	if err != nil {
		return
	}
	query := parsed.Query()
	if len(model.ClientId) == 0 {
		model.ClientId = firstNotEmpty(query.Get("clientId"), query.Get("client_id"))
	}
	if len(model.RedirectUri) == 0 {
		model.RedirectUri = firstNotEmpty(query.Get("redirectUri"), query.Get("redirect_uri"))
	}
	if len(model.State) == 0 {
		model.State = query.Get("state")
	}
	if len(model.ResponseType) == 0 {
		model.ResponseType = firstNotEmpty(query.Get("responseType"), query.Get("response_type"))
	}
	if len(model.GetCodeChallenge()) == 0 {
		model.CodeChallenge = firstNotEmpty(query.Get("codeChallenge"), query.Get("code_challenge"))
	}
	if len(model.GetCodeChallengeMethod()) == 0 {
		model.CodeChallengeMethod = firstNotEmpty(query.Get("codeChallengeMethod"), query.Get("code_challenge_method"))
	}
	if len(model.Nonce) == 0 {
		model.Nonce = query.Get("nonce")
	}
}

func formFirstValue(form map[string][]string, keys ...string) string {
	for _, key := range keys {
		values, ok := form[key]
		if ok && len(values) > 0 && len(values[0]) > 0 {
			return values[0]
		}
	}
	return ""
}

func firstNotEmpty(values ...string) string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return ""
}

func (r OAuthResource) selfInfo(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.AuthedUserInfo
		model     dtos.UserUpdate
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	userId := req.Attribute(config.RequestUserID)
	if userId == nil {
		errorData.MsgCode = config.MsgCodeUserInfoIsEmpty
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	if userId.(uint) != model.ID {
		errorData.MsgCode = config.MsgCodeStatusForbidden
		errorData.Lang = lang
		errorData.Err = fmt.Errorf("only can change your self information")
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.ChangeSelfInformation(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("add account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}
func (r OAuthResource) thirdAuthMethods(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	common.ResponseSuccess(resp, r.Svc.ThirdAuthMethods(ctx))
}

func (r OAuthResource) revoke(req *restful.Request, resp *restful.Response)     {}
func (r OAuthResource) introspect(req *restful.Request, resp *restful.Response) {}

func (r OAuthResource) loginByOIDC(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData  common.ErrorData
		result     dtos.AccessTokenResponse
		loginParam dtos.LoginByOIDC
	)

	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	// 获取登录信息，如果组织开启了MFA认证，则返回若clientid或者redirect uri为空，则认为是登录eauth的组织org，
	result, errorData = r.Svc.LoginByOIDC(ctx, loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}

func (r OAuthResource) loginBySAML(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData  common.ErrorData
		result     dtos.AccessTokenResponse
		loginParam dtos.LoginBySAML
	)

	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	contentType := req.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, config.RequestForm) {
		errorData.Err = req.Request.ParseForm()
		if errorData.IsNil() {
			loginParam.Provider = firstNotEmpty(req.Request.PostForm.Get("provider"), req.Request.Form.Get("provider"), req.QueryParameter("provider"))
			loginParam.SamlResponse = firstNotEmpty(req.Request.PostForm.Get("samlResponse"), req.Request.PostForm.Get("SAMLResponse"), req.Request.Form.Get("samlResponse"), req.Request.Form.Get("SAMLResponse"))
			loginParam.RelayState = firstNotEmpty(req.Request.PostForm.Get("relayState"), req.Request.PostForm.Get("RelayState"), req.Request.Form.Get("relayState"), req.Request.Form.Get("RelayState"))
		}
	} else {
		errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&loginParam)
	}
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	result, errorData = r.Svc.LoginBySAML(ctx, loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}

func (r OAuthResource) loginByLDAP(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.AccessTokenResponse
	)

	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	var loginParam dtos.LoginByLDAP
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	// 获取登录信息，如果组织开启了MFA认证，则返回若clientid或者redirect uri为空，则认为是登录eauth的组织org，
	result, errorData = r.Svc.LoginByLDAP(ctx, loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r OAuthResource) loginUsername(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.AccessTokenResponse
	)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	var loginParam dtos.LoginByUsername
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	result, errorData = r.Svc.LoginByUsername(ctx, loginParam)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Err = errors.New("username or password not right")
		errorData.MsgCode = config.MsgCodeUsernameOrPasswordNotRight
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}
func (r OAuthResource) publicKeys(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    jose.JSONWebKeySet
	)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	result, errorData = r.Svc.GetPublicKeys(ctx)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get org: %d public keys failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
func (r OAuthResource) oidcCfg(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	result := r.Svc.GetDiscoveryInfo(ctx)
	common.ResponseSuccess(resp, result)
}

func (r OAuthResource) refreshToken(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
	)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	tokenResponse, errorData := r.Svc.GetTokenByRefreshToken(ctx, req.QueryParameter("refresh_token"))
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, tokenResponse)
}
func (r OAuthResource) token(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
	)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	grantType := req.QueryParameter("grant_type")
	if len(grantType) == 0 {
		_ = req.Request.ParseForm()
		grantType = firstNotEmpty(req.Request.Form.Get("grant_type"), req.Request.Form.Get("grant_type "))
	}
	switch grantType {
	case config.GrantTypeAuthorizationCode:
		// 授权码
		r.Svc.GetTokenByAuthorizationCode(ctx, req, resp, lang)
	case config.GrantTypePassword:
		// 密码
		r.Svc.GetTokenByPassword(ctx, req, resp, lang)
	case config.GrantTypeToken:
		// 隐藏式
		r.Svc.GetTokenByToken(ctx, req, resp, lang)
	case config.GrantTypeClientCredentials:
		// 凭证式
		r.Svc.GetTokenByClientCredentials(ctx, req, resp, lang)
	case config.GrantTypeRefreshToken:
		r.Svc.RefreshToken(ctx, req, resp, lang)
	default:
		config.Logger.Errorf("unsupport grant_type: [%s]", grantType)
		errorData.MsgCode = config.MsgCodeUnsupportedGrantType
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Params = map[string]interface{}{"Type": grantType}
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

}

func (r OAuthResource) userinfo(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		userinfo  dtos.AuthedUserInfo
	)
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)
	userId := req.Attribute(config.RequestUserID)
	if userId == nil {
		errorData.MsgCode = config.MsgCodeUserInfoIsEmpty
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	userinfo, errorData = r.Svc.Userinfo(ctx, userId.(uint))
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, userinfo)
}
