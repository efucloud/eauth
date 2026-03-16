package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// MultiFactorAuthDetailList 用户MFA信息
type MultiFactorAuthDetailList struct {
	//当前页数据
	Data []*MultiFactorAuthDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// MultiFactorAuthDetail 用户MFA信息详情
type MultiFactorAuthDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//所属用户
	UserId uint `gorm:"user_id" json:"userId" validate:"required" description:"所属用户"`
	//密钥
	Secret string `gorm:"type:varchar(50);column:secret" json:"secret" validate:"required" description:"密钥"`
	//二维码
	Image string `gorm:"type:longtext;column:image" json:"image" validate:"required" description:"二维码"`
	//状态：是否已绑定
	Status string `gorm:"type:varchar(50);column:status;default:unbound" json:"status" validate:"oneof=bound unbound" enum:"bound|unbound"  description:"状态"`
}

// MultiFactorAuthCreate 用户MFA信息创建
type MultiFactorAuthCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//所属用户
	UserId uint `gorm:"user_id" json:"userId" validate:"required" description:"所属用户"`
	//密钥
	Secret string `gorm:"type:varchar(50);column:secret" json:"secret" validate:"required" description:"密钥"`
	//二维码
	Image string `gorm:"type:longtext;column:image" json:"image" validate:"required" description:"二维码"`
	//状态：是否已绑定
	Status string `gorm:"type:varchar(50);column:status;default:unbound" json:"status" validate:"oneof=bound unbound" enum:"bound|unbound"  description:"状态"`
}

func (ins *MultiFactorAuthCreate) Default(ctx context.Context) {
	ins.Status = "unbound"
}
func (ins *MultiFactorAuthCreate) Validate(ctx context.Context) (err error) {
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

// MultiFactorAuthStatus 禁用后，用户将不能使用该认证方式登陆系统
type MultiFactorAuthStatus struct {
	//主键
	Id uint ` json:"id" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//状态：是否已绑定
	Status string `gorm:"type:varchar(50);column:status;default:unbound" json:"status" validate:"oneof=bound unbound" enum:"bound|unbound"  description:"状态"`
}

func (ins *MultiFactorAuthStatus) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *MultiFactorAuthStatus) Validate(ctx context.Context) (err error) {
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

// PersonalBoundMFA 个人绑定MFA
type PersonalBoundMFA struct {
	Code   string `json:"code" description:"MFA验证码"`
	Client string `json:"client" description:"客户端信息"`
}
