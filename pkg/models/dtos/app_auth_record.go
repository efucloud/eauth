package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/utils"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// AppAuthRecordDetailList 应用认证记录列表响应
type AppAuthRecordDetailList struct {
	//当前页数据
	Data []*AppAuthRecordDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// AppAuthRecordDetail 应用认证记录详情
type AppAuthRecordDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用ID
	ApplicationId string `json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code" json:"code" validate:"required" description:"响应编码"`
	//PKCE挑战值
	CodeChallenge string `gorm:"type:varchar(255);column:code_challenge" json:"codeChallenge" description:"PKCE挑战值"`
	//PKCE挑战方式
	CodeChallengeMethod string `gorm:"type:varchar(20);column:code_challenge_method" json:"codeChallengeMethod" description:"PKCE挑战方式"`
	//用户ID
	UserId string `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//OIDC nonce
	Nonce string `gorm:"type:varchar(255);column:nonce" json:"nonce" description:"OIDC nonce"`
}

// AppAuthRecordCreate 应用认证记录创建
type AppAuthRecordCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//应用ID
	ApplicationId string `json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code" json:"code" validate:"required" description:"响应编码"`
	//PKCE挑战值
	CodeChallenge string `gorm:"type:varchar(255);column:code_challenge" json:"codeChallenge" description:"PKCE挑战值"`
	//PKCE挑战方式
	CodeChallengeMethod string `gorm:"type:varchar(20);column:code_challenge_method" json:"codeChallengeMethod" description:"PKCE挑战方式"`
	//用户ID
	UserId string `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//OIDC nonce
	Nonce string `gorm:"type:varchar(255);column:nonce" json:"nonce" description:"OIDC nonce"`
}

func (ins *AppAuthRecordCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	ins.ID = utils.GenerateDatabaseId()
}
func (ins *AppAuthRecordCreate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}
