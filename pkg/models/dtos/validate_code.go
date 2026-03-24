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

// ValidateCodeDetailList 系统操作校验码列表响应
type ValidateCodeDetailList struct {
	//当前页数据
	Data []*ValidateCodeDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ValidateCodeDetail 系统操作校验码详情
type ValidateCodeDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId string `gorm:"column:user_id" json:"userId" description:"用户ID"`
	//验证类型
	Category string `gorm:"type:varchar(50);column:category" json:"category" validate:"oneof=phone email" enum:"phone|email" description:"验证类型"`
	//验证码
	Code string `gorm:"type:varchar(50);column:code" json:"code" description:"验证码"`
	//动作
	Action string `gorm:"type:varchar(50);column:action" json:"action" validate:"oneof=registry login changepwd forgetpwd"  enum:"registry|login|changepwd|forgetpwd" description:"动作"`
	//过期时间
	Expired time.Time `gorm:"column:expired" json:"expired" validate:"required" description:"过期时间"`
}

// ValidateCodeCreate 系统操作校验码创建
type ValidateCodeCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//用户ID
	UserId string `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//验证类型
	Category string `gorm:"type:varchar(50);column:category" json:"category" validate:"oneof=phone email" enum:"phone|email" description:"验证类型"`
	//验证码
	Code string `gorm:"type:varchar(50);column:code" json:"code" description:"验证码"`
	//动作
	Action string `gorm:"type:varchar(50);column:action" json:"action" validate:"oneof=registry login changepwd forgetpwd"  enum:"registry|login|changepwd|forgetpwd" description:"动作"`
	//过期时间
	Expired time.Time `gorm:"column:expired" json:"expired" validate:"required" description:"过期时间"`
}

func (ins *ValidateCodeCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	ins.Expired = time.Now().Add(10 * time.Minute)
	ins.ID = utils.GenerateDatabaseId()
}
func (ins *ValidateCodeCreate) Validate(ctx context.Context) (err error) {
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
