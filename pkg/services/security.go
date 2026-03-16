package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/embeds"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"gopkg.in/gomail.v2"
	"html/template"
	"path"
	"time"
)

type SecurityService struct {
}

func (svc *SecurityService) UserInfo(ctx context.Context, code string) (user dtos.ShortUser, errorData common.ErrorData) {
	var validateCode dtos.ValidateCodeDetail
	validateSvc := ValidateCodeService{}
	validateCode, errorData = validateSvc.GetValidateCodeByCode(ctx, "email", code)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return
	}
	userSvc := UserService{}
	user, errorData = userSvc.GetShortUserByID(ctx, validateCode.UserId)
	return
}
func (svc *SecurityService) UserResetPassword(ctx context.Context, model dtos.UserResetPassword) (errorData common.ErrorData) {
	var validateCode dtos.ValidateCodeDetail
	validateSvc := ValidateCodeService{}
	validateCode, errorData = validateSvc.GetValidateCodeByCode(ctx, "email", model.Code)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return errorData
	}
	userSvc := UserService{}
	errorData = userSvc.ResetPassword(ctx, dtos.UserResetPassword{
		ID:       validateCode.UserId,
		Password: model.Password,
	})
	return
}
func (svc *SecurityService) SendEmail(ctx context.Context, lang, to string) (errorData common.ErrorData) {
	validateSvc := ValidateCodeService{}
	var (
		validate dtos.ValidateCodeCreate
		data     []byte
	)
	validate.Code = common.NewSecureID(30)
	validate.Category = "email"
	validate.Action = "forgetpwd"
	validate.Expired = time.Now().Add(10 * time.Minute)
	userSvc := UserService{}
	user, _ := userSvc.GetUserByEmail(ctx, to)
	validate.UserId = user.ID
	_, errorData = validateSvc.AddValidateCode(ctx, validate)
	//发送邮件
	params := make(map[string]string)
	params["email"] = to
	params["time"] = time.Now().Format("2006-01-02 15:04:05")
	params["expire"] = validate.Expired.Format("2006-01-02 15:04:05")
	params["server"] = config.ApplicationConfig.ServerAddress
	params["link"] = fmt.Sprintf("%s/reset/password/%s", config.ApplicationConfig.ServerAddress, validate.Code)
	data, errorData.Err = embeds.Templates.ReadFile(path.Join("templates", fmt.Sprintf("password.%s.tpl", lang)))
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return
	}
	t, _ := template.New("email").Delims("_{{_", "_}}_").Parse(string(data))
	b := new(bytes.Buffer)
	errorData.Err = t.Execute(b, params)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return
	}
	client := config.NewMailClient()
	if client == nil {
		errorData.Err = fmt.Errorf("mail client init failed")
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", config.ApplicationConfig.Email.Username)
	m.SetHeader("To", to)
	if lang == "en" {
		m.SetHeader("Subject", "EAuth Reset Password")
	} else {
		m.SetHeader("Subject", "易认证重置密码")
	}
	m.SetBody("text/html", b.String())
	errorData.Err = client.DialAndSend(m)
	return errorData

}
