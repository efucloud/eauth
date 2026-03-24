package system

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

type MultiFactorAuthResource struct {
	Svc services.MultiFactorAuthService
}

func (r MultiFactorAuthResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "mfa"
	apiInfo.Description = "MFA"
	apiExtend := "/mfa"
	common.RegisterApiInfo(apiInfo)
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/personal").
		Doc("获取个人MFA信息列表").
		Notes("获取个人MFA信息列表").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.personalMfaInfo).
		Returns(http.StatusOK, "", dtos.MultiFactorAuthDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getMultiFactorAuthPersonal"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/reset").
		Doc("重置个人MFA信息").
		Notes("重置个人MFA信息，会重新生成新的").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.reset).
		Reads(dtos.PersonalBoundMFA{}).
		Returns(http.StatusOK, "", dtos.MultiFactorAuthDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "resetMultiFactorAuthPersonal"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/bound").
		Doc("个人重新绑定MFA").
		Notes("个人重新绑定MFA， 在个人中心页面操作 ").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.bound).
		Reads(dtos.PersonalBoundMFA{}).
		Returns(http.StatusOK, "", dtos.MultiFactorAuthDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "boundMultiFactorAuthPersonal"))
	ws.Route(ws.DELETE(config.V1Prefix+apiExtend).
		Doc("删除MFA数据").
		Notes("删除MFA数据信息详情").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.delete).
		Reads(dtos.BatchOperationIds{}).
		Returns(http.StatusOK, "成功", "成功").
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "deleteUserMultiFactorAuth"))

}
func (r MultiFactorAuthResource) bound(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.PersonalBoundMFA
		result    dtos.MultiFactorAuthDetail
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
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
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result = r.Svc.BoundUserMultiFactorAuth(ctx, userId.(string), model.Code)
	common.ResponseSuccess(resp, result)
}
func (r MultiFactorAuthResource) list(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.MultiFactorAuthDetailList
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	current, pageSize, order := common.GetRequestPaginationInformation(req)
	queryParam := &common.QueryParam{}
	common.RequestQuery("userId", common.ParamTypeNumber, common.QueryTypeEqual, req, queryParam)
	result, errorData = r.Svc.ListMultiFactorAuth(ctx, current, pageSize, order, queryParam.WhereQuery, queryParam.WhereArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("list account failed, err: %s", errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}

func (r MultiFactorAuthResource) reset(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.PersonalBoundMFA
		result    dtos.MultiFactorAuthDetail
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
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
	errorData.Err = jsoniter.NewDecoder(req.Request.Body).Decode(&model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("decode json format data failed, err: %s", errorData.Err.Error())
		errorData.MsgCode = config.MsgCodeJsonDecodeFailed
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.ResetUserMultiFactorAuth(ctx, userId.(string), model.Code)
	if errorData.IsNotNil() {
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
func (r MultiFactorAuthResource) delete(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.BatchOperationIds
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
	errorData = r.Svc.DeleteMultiFactorAuth(ctx, model.Ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, "删除成功")
}
func (r MultiFactorAuthResource) personalMfaInfo(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.MultiFactorAuthDetail
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
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
	result = r.Svc.GetUserMultiFactorAuthByUserId(ctx, userId.(string))
	if result.Status == "bound" {
		result.Secret = ""
		result.Image = ""
	}
	common.ResponseSuccess(resp, result)

}
