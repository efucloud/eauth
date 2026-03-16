package system

import (
	"context"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/apis/filters"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/services"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	jsoniter "github.com/json-iterator/go"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type PersonalResource struct {
	Svc services.PersonalService
}

func (r PersonalResource) AddWebService(ws *restful.WebService) {
	apiInfo := common.ApiInfo{}
	apiInfo.Tag = "personal"
	apiInfo.Description = "系统个人信息"
	common.RegisterApiInfo(apiInfo)
	apiExtend := "/personal"

	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/self/avatar").
		Doc("更新个人头像").
		Notes("更新个人头像").
		To(r.selfAvatar).
		Param(ws.FormParameter("file", "上传的头像文件").DataType("File")).
		Consumes(config.RequestFormData).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "updateSelfAvatar"))
	ws.Route(ws.POST(config.V1Prefix+apiExtend+"/set/pwd").
		Doc("设置个人密码").
		Notes("设置个人密码").
		Param(ws.HeaderParameter(config.AuthHeader, "请求Token").Required(true)).
		To(r.setPassword).
		Reads(dtos.SetPassword{}).
		Filter(filters.Log).Filter(filters.I18n).Filter(filters.Auth).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()).
		Metadata(config.FrontApiTag, "setPassword"))
}

func (r PersonalResource) setPassword(req *restful.Request, resp *restful.Response) {
	var (
		errorData common.ErrorData
		model     dtos.SetPassword
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
	errorData = r.Svc.SetPassword(ctx, userId.(uint), model)
	if errorData.IsNotNil() {
		config.Logger.Errorf("set password failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	common.ResponseSuccess(resp, "success")
}

func (r PersonalResource) selfAvatar(req *restful.Request, resp *restful.Response) {
	var (
		errorData  common.ErrorData
		file       multipart.File
		fileHeader *multipart.FileHeader
		data       []byte
		fileSuffix string
	)
	lang := common.GetLanguageFromReq(req, config.RequestLanguage)
	errorData.Lang = lang
	ctx := context.WithValue(context.Background(), config.RequestLanguage, lang)
	if reqCtx := req.Attribute(config.RequestContext); reqCtx != nil {
		ctx = reqCtx.(context.Context)
	}
	ctx = context.WithValue(ctx, config.RequestLanguage, lang)

	file, fileHeader, errorData.Err = req.Request.FormFile("file")
	data, errorData.Err = io.ReadAll(file)
	sp := strings.Split(fileHeader.Filename, ".")
	if len(sp) > 0 {
		fileSuffix = fmt.Sprintf(".%s", sp[len(sp)-1])
	}
	fileName := fmt.Sprintf("%s%s", common.MD5VByte(data), fileSuffix)
	filePath := fmt.Sprintf("%s/%s/%s", config.ApplicationConfig.UploadPath, config.UserAvatars, fileName)
	publicPath := fmt.Sprintf("/public/%s/%s", config.UserAvatars, fileName)
	userId := req.Attribute(config.RequestUserID)
	if userId == nil {
		errorData.MsgCode = config.MsgCodeUserInfoIsEmpty
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	errorData.Err = os.WriteFile(filePath, data, os.ModePerm)
	if errorData.Err != nil {
		errorData.MsgCode = config.MsgCodeGetRecordFailed
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	errorData = r.Svc.UpdateUserAvatar(ctx, userId.(uint), publicPath)
	if errorData.IsNotNil() {
		config.Logger.Errorf("updata user avatar failed, err: %s", errorData.Err.Error())
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, publicPath)

}
