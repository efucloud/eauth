package system

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

type ProviderOidcResource struct {
	Svc services.ProviderOidcService
}

func (r ProviderOidcResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "provider-oidc"
	apiInfo.Description = "OIDC认证提供商"
	common.RegisterApiInfo(apiInfo)
	apiExtend := "/provider-oidc"
	ws.Route(ws.POST(config.V1Prefix+apiExtend).
		Doc("创建OIDC认证提供商").
		Notes("创建OIDC认证提供商信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.create).
		Reads(dtos.ProviderOidcCreate{}).
		Returns(http.StatusOK, "成功", dtos.ProviderOidcDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "createProviderOidc"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend).
		Doc("获取OIDC认证提供商列表").
		Notes("获取OIDC认证提供商信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		Param(ws.QueryParameter("current", "页码").DataType("number")).
		Param(ws.QueryParameter("pageSize", "每页大小").DataType("number")).
		Param(ws.QueryParameter("order", "排序")).
		Param(ws.QueryParameter("name", "名称").DataType("string")).
		Param(ws.QueryParameter("code", "编码").DataType("string")).
		To(r.list).
		Returns(http.StatusOK, "成功", dtos.ProviderOidcDetailList{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "listProviderOidc"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/{id}").
		Doc("获取OIDC认证提供商详情").
		Notes("获取OIDC认证提供商信息详情").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.get).
		Param(ws.PathParameter("id", "记录ID").DataType("number")).
		Returns(http.StatusOK, "成功", dtos.ProviderOidcDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getProviderOidc"))
	ws.Route(ws.PUT(config.V1Prefix+apiExtend).
		Doc("更新OIDC认证提供商信息").
		Notes("更新OIDC认证提供商信息").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.update).
		Reads(dtos.ProviderOidcUpdate{}).
		Returns(http.StatusOK, "成功", dtos.ProviderOidcDetail{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "updateProviderOidc"))
	ws.Route(ws.DELETE(config.V1Prefix+apiExtend).
		Doc("删除OIDC认证提供商").
		Notes("删除OIDC认证提供商信息详情").
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
		Metadata(config.FrontApiTag, "deleteProviderOidc"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/status").
		Doc("启用禁用").
		Notes("启用禁用,修改账户状态").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.status).
		Reads(dtos.ProviderOidcStatus{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "changeProviderOidcStatus"))

}

func (r ProviderOidcResource) status(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.ProviderOidcStatus
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

	errorData = r.Svc.ChangeStatus(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("enable oidc provider failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, "success")
}
func (r ProviderOidcResource) get(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	var (
		errorData common.ErrorData
		result    dtos.ProviderOidcDetail
	)
	errorData.Lang = lang
	errorData.Lang = lang
	ctx := context.Background()
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	id := common.StringsToUint(req.PathParameter("id"))
	if id < 1 {
		errorData.MsgCode = config.MsgCodePathIdInvalid
		errorData.ResponseCode = http.StatusBadRequest
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	result, errorData = r.Svc.GetProviderOidcById(ctx, id)
	if !errorData.IsNil() {
		config.Logger.Errorf("get oidc proivder by id: %d failed, err: %s", id, errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
func (r ProviderOidcResource) delete(req *restful.Request, resp *restful.Response) {
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

	errorData = r.Svc.DeleteProviderOidc(ctx, model.Ids)
	if errorData.IsNotNil() {
		config.Logger.Errorf("delete account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, "删除成功")
}
func (r ProviderOidcResource) create(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.ProviderOidcDetail
		model     dtos.ProviderOidcCreate
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
	//判断OIDC认证提供商是否允许创建OIDC认证提供商
	result, errorData = r.Svc.AddProviderOidc(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("add account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r ProviderOidcResource) update(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.ProviderOidcUpdate
		result    dtos.ProviderOidcDetail
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

	result, errorData = r.Svc.UpdateProviderOidc(ctx, model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("update account failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, result)
}

func (r ProviderOidcResource) list(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		result    dtos.ProviderOidcDetailList
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
	result, errorData = r.Svc.ListProviderOidc(ctx, current, pageSize, order, queryParam.WhereQuery, queryParam.WhereArgs)
	if errorData.IsNotNil() {
		config.Logger.Errorf("list account failed, err: %s", errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, result)
}
