package v1

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"net/http"

	"github.com/efucloud/eauth/pkg/models/dtos"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/services"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

type UserTokenResource struct {
	Svc services.UserTokenService
}

func (r UserTokenResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "user-token"
	apiInfo.Description = "系统用户令牌"
	common.RegisterApiInfo(apiInfo)
	apiExtend := "/user-token"
	ws.Route(ws.POST(config.V1Prefix+apiExtend).
		Doc("创建系统用户令牌").
		Notes("创建系统用户令牌信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.create).
		Reads(dtos.UserTokenCreate{}).
		Returns(http.StatusOK, "成功", dtos.UserTokenDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "createUserToken"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend).
		Doc("获取系统用户令牌列表").
		Notes("获取系统用户令牌信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		Param(ws.QueryParameter("current", "页码").DataType("number")).
		Param(ws.QueryParameter("pageSize", "每页大小").DataType("number")).
		Param(ws.QueryParameter("order", "排序")).
		To(r.list).
		Returns(http.StatusOK, "成功", dtos.UserTokenDetailList{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "listUserToken"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/{id}").
		Doc("获取系统用户令牌详情").
		Notes("获取系统用户令牌信息详情").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.get).
		Param(ws.PathParameter("id", "记录ID").DataType("string")).
		Returns(http.StatusOK, "成功", dtos.UserTokenDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getUserToken"))
	ws.Route(ws.PUT(config.V1Prefix+apiExtend).
		Doc("更新系统用户令牌信息").
		Notes("更新系统用户令牌信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.update).
		Reads(dtos.UserTokenUpdate{}).
		Returns(http.StatusOK, "成功", dtos.UserTokenDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "updateUserToken"))
	ws.Route(ws.DELETE(config.V1Prefix+apiExtend).
		Doc("删除系统用户令牌").
		Notes("删除系统用户令牌信息详情").
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
		Metadata(config.FrontApiTag, "deleteUserToken"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/personal/authed").
		Doc("获取个人的所有访问令牌").
		Notes("获取个人的所有访问令牌").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.personal).
		Param(ws.QueryParameter("current", "页码").DataType("number")).
		Param(ws.QueryParameter("pageSize", "每页大小").DataType("number")).
		Param(ws.QueryParameter("order", "排序")).
		Returns(http.StatusOK, "成功", dtos.UserAuthProfileDetailList{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getUserPersonalUserToken"))

}

func (r UserTokenResource) personal(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.UserTokenDetailList
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
	result, errorData = r.Svc.GetUserTokensByUserId(ctx, userId.(string))
	if errorData.IsNotNil() {
		config.Logger.Errorf("get personal auth profile failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
func (r UserTokenResource) get(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.UserTokenDetail
	)
	errorData.Lang = lang
	errorData.Lang = lang
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	result, errorData = r.Svc.GetUserTokenDetailById(ctx, req.PathParameter("id"))
	if !errorData.IsNil() {
		config.Logger.Errorf("get oidc proivder by id: %s failed, err: %s", req.PathParameter("id"), errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
func (r UserTokenResource) delete(req *restful.Request, resp *restful.Response) {
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

	errorData = r.Svc.DeleteUserToken(ctx, model.Ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, "删除成功")
}
func (r UserTokenResource) create(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.UserTokenDetail
		model     dtos.UserTokenCreate
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
	//判断系统用户令牌是否允许创建系统用户令牌
	result, errorData = r.Svc.AddUserToken(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("add account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r UserTokenResource) update(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.UserTokenUpdate
		result    dtos.UserTokenDetail
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

	result, errorData = r.Svc.UpdateUserToken(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("update account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r UserTokenResource) list(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.UserTokenDetailList
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
	result, errorData = r.Svc.ListUserToken(ctx, current, pageSize, order, queryParam.WhereQuery, queryParam.WhereArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("list account failed, err: %s", errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
