package system

import (
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	restful "github.com/emicklei/go-restful/v3"
	"net/http"

	"github.com/efucloud/common"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
)

type InfoResource struct {
}

func (r InfoResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "info"
	apiInfo.Description = "应用信息"
	common.RegisterApiInfo(apiInfo)
	ws.Route(ws.GET(config.V1Prefix+"/health").
		Doc("健康检查").
		Notes("健康检查").
		To(r.health).
		Returns(http.StatusOK, "成功", "ok").
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "health"))
	ws.Route(ws.GET(config.V1Prefix+"/info").
		Doc("查看应用信息").
		Notes("查看应用的编译信息").
		To(r.info).
		Returns(http.StatusOK, "成功", common.ApplicationInfo{}).
		Returns(http.StatusInternalServerError, "内部处理逻辑错误", dtos.ResponseError{}).
		Filter(filters.Log).Filter(filters.I18n).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "getAppInformation"))

}
func (r InfoResource) info(req *restful.Request, resp *restful.Response) {
	var app common.ApplicationInfo
	app.Application = config.ApplicationName
	app.BuildDate = config.BuildDate
	app.GoVersion = config.GoVersion
	app.Commit = config.Commit
	app.Data = common.GetLanguageFromReq(req, config.RequestLanguage)
	common.ResponseSuccess(resp, app)
}

func (r InfoResource) health(req *restful.Request, resp *restful.Response) {
	common.ResponseSuccess(resp, "ok")
}
