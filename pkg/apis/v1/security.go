package v1

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type SecurityResource struct {
	Svc services.SecurityService
}

func (r SecurityResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "security"
	apiInfo.Description = "安全"
	common.RegisterApiInfo(apiInfo)
	apiExtend := ""
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/email").
		Doc("判断邮箱在系统中是否存在").
		Notes("判断邮箱在系统中是否存在").
		Param(ws.QueryParameter("email", "邮箱")).
		To(r.emailExist).
		Returns(http.StatusOK, "成功", dtos.ExistResponse{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemEmailExist"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/phone").
		Doc("判断手机号码在系统中是否存在").
		Notes("判断手机号码在系统中是否存在").
		Param(ws.QueryParameter("phone", "手机号码")).
		To(r.phoneExist).
		Returns(http.StatusOK, "成功", dtos.ExistResponse{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemPhoneExist"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/email").
		Doc("发送重置邮件").
		Notes("发送重置邮件").
		To(r.sendResetPwdEmail).
		Reads(dtos.ExistResponse{}).
		Returns(http.StatusOK, "成功", "").
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "sendResetPwdEmail"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/user/{code}").
		Doc("获取重置密码用户信息").
		Notes("获取重置密码用户信息").
		Param(ws.PathParameter("code", "重置校验码")).
		To(r.systemUserInfo).
		Returns(http.StatusOK, "成功", dtos.ShortUser{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemValidateUserInfo"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/resetpwd").
		Doc("重置密码").
		Notes("重置密码").
		To(r.systemUserResetPassword).
		Reads(dtos.UserResetPassword{}).
		Returns(http.StatusOK, "成功", "").
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemUserResetPassword"))

}

func (r SecurityResource) systemUserInfo(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.ShortUser
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	config.Logger.Info(lang)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	result, errorData = r.Svc.UserInfo(ctx, req.PathParameter("code"))
	if errorData.IsNotNil() {
		errorData.Lang = lang
		errorData.MsgCode = config.MsgCodeRecordNotExist
		errorData.ResponseCode = http.StatusInternalServerError
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}

func (r SecurityResource) systemUserResetPassword(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.UserResetPassword
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	config.Logger.Info(lang)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	errorData = r.Svc.UserResetPassword(ctx, model)
	if errorData.IsNotNil() {
		errorData.Lang = lang
		errorData.MsgCode = config.MsgCodeResetPasswordFailed
		errorData.ResponseCode = http.StatusInternalServerError
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, "success")

}
func (r SecurityResource) sendResetPwdEmail(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.ExistResponse
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	config.Logger.Info(lang)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	errorData = r.Svc.SendEmail(ctx, lang, model.Name)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		errorData.MsgCode = config.MsgCodeSendEmailFailed
		errorData.ResponseCode = http.StatusInternalServerError
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, "success")

}
func (r SecurityResource) emailExist(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		user      dtos.UserDetail
		res       dtos.ExistResponse
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	userSvc := services.UserService{}
	user, errorData = userSvc.GetUserByEmail(ctx, req.QueryParameter("email"))
	res.Name = user.Email
	res.Exist = len(user.Email) > 0
	common.ResponseSuccess(resp, res)
}

func (r SecurityResource) phoneExist(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		user      dtos.UserDetail
		res       dtos.ExistResponse
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	userSvc := services.UserService{}
	user, errorData = userSvc.GetUserByPhone(ctx, req.QueryParameter("phone"))
	res.Name = user.Phone
	res.Exist = len(user.Phone) > 0
	common.ResponseSuccess(resp, res)
}
