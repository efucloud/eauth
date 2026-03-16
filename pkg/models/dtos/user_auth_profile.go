package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// UserAuthProfileDetailList 系统用户认证方式列表响应
type UserAuthProfileDetailList struct {
	//当前页数据
	Data []*UserAuthProfileDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// UserAuthProfileDetail 系统用户认证方式详情
type UserAuthProfileDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"本系统用户ID"`
	//认证类型provider中的code,email,phone
	Provider string `gorm:"type:varchar(255);column:provider" json:"provider" description:"认证类型"`
	//第三方登录用户的id，邮箱，手机号
	LoginID string `gorm:"type:varchar(255);column:login_id" json:"loginId" description:"第三方认证的用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//第三方认证的全部用户信息
	Properties JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`

	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}

// UserAuthProfileCreate 系统用户认证方式创建
type UserAuthProfileCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"本系统用户ID"`
	//认证类型provider中的code,email,phone
	Provider string `gorm:"type:varchar(255);column:provider" json:"provider" validate:"required" description:"认证类型"`
	//第三方登录用户的id，邮箱，手机号
	LoginID string `gorm:"type:varchar(255);column:login_id" json:"loginId" validate:"required" description:"第三方认证的用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//第三方认证的全部用户信息
	Properties JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`

	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}

func (ins *UserAuthProfileCreate) Default(ctx context.Context) {
	ins.LatestUsedTime = time.Now().Format(time.DateTime)
	ins.CreatedAt = time.Now()

}
func (ins *UserAuthProfileCreate) Validate(ctx context.Context) (err error) {
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

// UserAuthProfileUpdate 系统用户认证方式修改
type UserAuthProfileUpdate struct {
	//记录ID
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"本系统用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//第三方认证的全部用户信息
	Properties JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`

	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}

func (ins *UserAuthProfileUpdate) Default(ctx context.Context) {
	ins.LatestUsedTime = time.Now().Format(time.DateTime)
	ins.UpdatedAt = time.Now()
}
func (ins *UserAuthProfileUpdate) Validate(ctx context.Context) (err error) {
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

// UserAuthProfileStatus 禁用后，用户将不能使用该认证方式登陆系统
type UserAuthProfileStatus struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (ins *UserAuthProfileStatus) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *UserAuthProfileStatus) Validate(ctx context.Context) (err error) {
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
