package system

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

type DashboardResource struct {
	Svc services.DashboardService
}

func (r DashboardResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "dashboard"
	apiInfo.Description = "仪表盘"
	common.RegisterApiInfo(apiInfo)
	apiExtend := "/dashboard"
	ws.Route(ws.GET(config.V1Prefix+apiExtend).
		Doc("系统面板").
		Notes("系统面板").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.system).
		Returns(http.StatusOK, "成功", dtos.Dashboard{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "systemDashboard"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/application-auth-top10").
		Doc("系统应用认证TOP10").
		Notes("系统应用认证TOP10").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.applicationAuthTop10).
		Returns(http.StatusOK, "成功", []dtos.ApplicationAuthTop{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "applicationAuthTop10"))
	ws.Route(ws.GET(config.V1Prefix+apiExtend+"/application-auth-30days").
		Doc("30天内应用认证统计").
		Notes("30天内应用认证统计").
		Param(ws.HeaderParameter(config.AuthHeader, "系统用户Token")).
		To(r.applicationAuth30Days).
		Returns(http.StatusOK, "成功", []dtos.DashboardData[int64]{}).
		Returns(http.StatusBadRequest, "请求数据无法处理", dtos.ResponseError{}).
		Returns(http.StatusForbidden, "用户没有权限", dtos.ResponseError{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Filter(filters.Permission([]string{config.RoleAdmin, config.RoleEdit, config.RoleView})).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "applicationAuth30Days"))
}
func (r DashboardResource) system(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	common.ResponseSuccess(resp, r.Svc.Dashboard(ctx))
}
func (r DashboardResource) applicationAuthTop10(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	common.ResponseSuccess(resp, r.Svc.ApplicationAuthTop10(ctx))
}

func (r DashboardResource) applicationAuth30Days(req *restful.Request, resp *restful.Response) {
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	common.ResponseSuccess(resp, r.Svc.ApplicationAuth30Days(ctx))
}
