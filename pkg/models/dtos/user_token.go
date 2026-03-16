package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// UserTokenDetailList 系统用户Token列表响应
type UserTokenDetailList struct {
	//当前页数据
	Data []*UserTokenDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// UserTokenDetail 系统用户Token详情
type UserTokenDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//用户
	User ShortUser `gorm:"-" json:"user" description:"用户"`
	//客户端ID
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//过期时间(时间戳)
	Expired int64 `json:"expired" description:" 过期时间(时间戳)"`
	//过期时间
	ExpiredTime time.Time `json:"expiredTime" description:"过期时间"`
	//RefreshToken
	RefreshToken string `gorm:"type:varchar(50)" json:"refreshToken,omitempty" description:"RefreshToken"`
	//Claims的ID
	ClaimsID string `gorm:"type:varchar(50)" json:"-" validate:"required" description:"Claims的ID"`
	//session key， token MD5
	SessionKey string `gorm:"type:varchar(50)" json:"sessionKey" description:"SessionKey"`
	//Token
	Token string `gorm:"type:longtext" json:"token" description:"Token"`
}

// UserTokenCreate 系统用户Token创建
type UserTokenCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//客户端ID
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//过期时间(时间戳)
	Expired int64 `json:"expired" description:" 过期时间(时间戳)"`
	//过期时间
	ExpiredTime time.Time `json:"expiredTime" description:"过期时间"`
	//RefreshToken
	RefreshToken string `gorm:"type:varchar(50)" json:"refreshToken,omitempty" description:"RefreshToken"`
	//Claims的ID
	ClaimsID string `gorm:"type:varchar(50)" json:"-" validate:"required" description:"Claims的ID"`
	//session key， token MD5
	SessionKey string `gorm:"type:varchar(50)" json:"sessionKey" description:"SessionKey"`
	//Token
	Token string `gorm:"type:longtext" json:"token" description:"Token"`
}

func (ins *UserTokenCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
}
func (ins *UserTokenCreate) Validate(ctx context.Context) (err error) {
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

// UserTokenUpdate 系统用户Token修改
type UserTokenUpdate struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//客户端ID
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//过期时间(时间戳)
	Expired int64 `json:"expired" description:" 过期时间(时间戳)"`
	//过期时间
	ExpiredTime time.Time `json:"expiredTime" description:"过期时间"`
	//RefreshToken
	RefreshToken string `gorm:"type:varchar(50)" json:"refreshToken,omitempty" description:"RefreshToken"`
	//Claims的ID
	ClaimsID string `gorm:"type:varchar(50)" json:"-" validate:"required" description:"Claims的ID"`
	//session key， token MD5
	SessionKey string `gorm:"type:varchar(50)" json:"sessionKey" description:"SessionKey"`
	//Token
	Token string `gorm:"type:longtext" json:"token" description:"Token"`
}

func (ins *UserTokenUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *UserTokenUpdate) Validate(ctx context.Context) (err error) {
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
