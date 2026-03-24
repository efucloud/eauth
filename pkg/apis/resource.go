package apis

import (
	"fmt"
	v1 "github.com/efucloud/eauth/pkg/apis/v1"
	"github.com/efucloud/eauth/pkg/embeds"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"path"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

func GetWebServices(container *restful.Container) *restful.WebService {
	ws := new(restful.WebService)
	container.RecoverHandler(func(i interface{}, writer http.ResponseWriter) {
		writer.WriteHeader(http.StatusInternalServerError)
		var body dtos.ResponseError
		body.Detail = fmt.Sprintf("%v", i)
		body.Alert, _ = common.GetLocaleMessage(config.Bundle, nil, "zh", config.MsgCodeApplicationUnExceptPanicError)

	})
	ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	v1.ApplicationResource{Svc: services.ApplicationService{}}.AddWebService(ws)
	v1.DashboardResource{Svc: services.DashboardService{}}.AddWebService(ws)
	v1.FaceRecognitionResource{Svc: services.FaceRecognitionService{}}.AddWebService(ws)
	v1.ProviderLdapResource{Svc: services.ProviderLdapService{}}.AddWebService(ws)
	v1.ProviderOidcResource{Svc: services.ProviderOidcService{}}.AddWebService(ws)
	v1.ProviderSamlResource{Svc: services.ProviderSamlService{}}.AddWebService(ws)
	v1.UserResource{Svc: services.UserService{}}.AddWebService(ws)
	v1.UserAuthProfileResource{Svc: services.UserAuthProfileService{}}.AddWebService(ws)
	v1.UserTokenResource{Svc: services.UserTokenService{}}.AddWebService(ws)
	v1.InfoResource{}.AddWebService(ws)
	v1.MultiFactorAuthResource{Svc: services.MultiFactorAuthService{}}.AddWebService(ws)
	v1.OAuthResource{Svc: services.OAuthService{}}.AddWebService(ws)
	v1.PersonalResource{}.AddWebService(ws)
	v1.SecurityResource{Svc: services.SecurityService{}}.AddWebService(ws)

	return ws
}
func AddResources() {
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	container := restful.DefaultContainer
	container.RecoverHandler(func(i interface{}, writer http.ResponseWriter) {
		writer.WriteHeader(http.StatusInternalServerError)
		var res common.ResponseError
		res.Alert = "Unkown error"
		res.Message = fmt.Sprintf("message: %v", i)
		data, _ := jsoniter.Marshal(res)
		_, _ = writer.Write(data)
	})
	container.Router(restful.CurlyRouter{})
	container.Filter(container.OPTIONSFilter)
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "*"},
		CookiesAllowed: true,
		Container:      container,
	}

	container.Filter(cors.Filter)
	ws := GetWebServices(container)
	container.Add(ws)
	ws.Route(ws.GET("/public/" + config.UserAvatars + "/{subpath:*}").
		Param(ws.PathParameter("subpath", "路径")).To(static(config.UserAvatars)))
	http.Handle("/face-recognition/", http.StripPrefix("/face-recognition/", http.FileServer(embeds.FaceRecognitionModels())))

}
func static(staticDir string) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		actual := path.Join(config.ApplicationConfig.UploadPath, staticDir, req.PathParameter("subpath"))
		http.ServeFile(resp.ResponseWriter, req.Request, actual)
	}
}

func AddSwagger() {
	c := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), //you control what services are visible
		APIPath:                       "/api/v1/swagger.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(c))
}
func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "EAuth",
			Description: "Efu Cloud旗下产品EAuth(统一认证平台)，若采用的是私有化部署，短信发送，邮件发送默认为组织efu-cloud中的配置",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "efucloud",
					Email: "efucloud@aliyun.com",
					URL:   ""},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org"},
			},
			Version: "v1.0.0",
		},
	}
	swo.Tags = []spec.Tag{}
	for _, s := range common.ApiInfos {
		props := spec.TagProps{Name: s.Tag, Description: s.Description}
		swo.Tags = append(swo.Tags, spec.Tag{TagProps: props})
	}
}
