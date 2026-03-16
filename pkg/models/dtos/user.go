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

// UserDetailList  账户列表响应
type UserDetailList struct {
	//当前页数据
	Data []*ShortUser `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ShortUser 简单账户详情
type ShortUser struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//用户名
	Username string `gorm:"type:varchar(255);column:username" json:"username" validate:"alphanum" description:"用户名"`
	//昵称，如中文名
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `gorm:"type:varchar(255);column:job_number" json:"jobNumber" description:"工号"`
	//系统角色
	Role string `gorm:"type:varchar(255);column:role;default:none" json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//邮箱
	Email string `gorm:"type:varchar(255);column:email" json:"email" validate:"email" description:"邮箱"`
	//手机号码
	Phone string `gorm:"type:varchar(255);column:phone" json:"phone" validate:"required" description:"电话"`
	//默认语言
	Language string `gorm:"type:varchar(255);column:language;default:zh" json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
	//头像
	Avatar string `gorm:"type:varchar(1000);column:avatar" json:"avatar" description:"头像"`
	//是否绑定MFA认证
	MFA bool `gorm:"column:mfa" json:"mfa" description:"是否绑定MFA认证"`
}

// UserDetail 账户详情
type UserDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户名
	Username string `gorm:"type:varchar(255);column:username" json:"username" validate:"alphanum" description:"用户名"`
	//昵称，如中文名
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `gorm:"type:varchar(255);column:job_number" json:"jobNumber" description:"工号"`
	//密码
	Password string `gorm:"-" json:"password" example:"admin" description:"密码"`
	//数据库保存的加密密码
	PasswordStore string `gorm:"type:varchar(255);column:password_store" json:"passwordStore"`
	//密码强度
	PasswordStrength string `gorm:"type:varchar(255);column:password_strength;default:weak" json:"passwordStrength" validate:"oneof=strong medium weak" enum:"strong|medium|weak" description:"密码强度"`
	//系统角色
	Role string `gorm:"type:varchar(255);column:role;default:none" json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//邮箱
	Email string `gorm:"type:varchar(255);column:email" json:"email" validate:"email" description:"邮箱"`
	//手机号码
	Phone string `gorm:"type:varchar(255);column:phone" json:"phone" validate:"required" description:"电话"`
	//默认语言
	Language string `gorm:"type:varchar(255);column:language;default:zh" json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
	//头像
	Avatar string `gorm:"type:varchar(1000);column:avatar" json:"avatar" description:"头像"`
	// 人脸识别数据
	FaceIdDatas []ArrayFloat64 `gorm:"-" json:"-" validate:"-" description:"人脸识别数据"`
	//是否绑定MFA认证
	MFA bool `gorm:"column:mfa" json:"mfa" description:"是否绑定MFA认证"`
}

// UserCreate 账户信息创建
// 未来账户信息修改只能从eauth中
type UserCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//用户名
	Username string `gorm:"type:varchar(255);column:username" json:"username" validate:"alphanum" description:"用户名"`
	//昵称，如中文名
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `gorm:"type:varchar(255);column:job_number" json:"jobNumber" description:"工号"`
	//密码
	Password string `gorm:"-" json:"password" example:"admin" description:"密码"`
	//数据库保存的加密密码
	PasswordStore string `gorm:"type:varchar(255);column:password_store" json:"passwordStore"`
	//密码强度
	PasswordStrength string `gorm:"type:varchar(20);column:password_strength;default:weak" json:"passwordStrength" validate:"oneof=strong medium weak" enum:"strong|medium|weak" description:"密码强度"`
	//系统角色
	Role string `gorm:"type:varchar(255);column:role;default:none" json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//邮箱
	Email string `gorm:"type:varchar(255);column:email" json:"email" validate:"email" description:"邮箱"`
	//手机号码
	Phone string `gorm:"type:varchar(255);column:phone" json:"phone" validate:"required" description:"电话"`
	//默认语言
	Language string `gorm:"type:varchar(255);column:language;default:zh" json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
	//头像
	Avatar string `gorm:"type:varchar(1000);column:avatar" json:"-" description:"头像"`
}

func (ins *UserCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	if len(ins.Language) == 0 {
		ins.Language = "zh"
	}
	if len(ins.PasswordStrength) == 0 {
		ins.PasswordStrength = "weak"
	}
	ins.PasswordStrength = utils.CheckPasswordStrength(ins.Password)
}
func (ins *UserCreate) Validate(ctx context.Context) (err error) {
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

// UserStatus 账户信息禁用/启用
// 账户禁用后，用户将不能登陆该系统
type UserStatus struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (ins *UserStatus) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *UserStatus) Validate(ctx context.Context) (err error) {
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

// UserRole 账户系统角色设置
// 设置账户在系统中的角色
type UserRole struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//用户在系统中的角色
	Role string `gorm:"type:varchar(255);default:none;column:role" json:"role" validate:"oneof:admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
}

func (ins *UserRole) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *UserRole) Validate(ctx context.Context) (err error) {
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

// UserUpdate 账户信息更新
// 更新账户信息，未来只能在eauth中更新
type UserUpdate struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//用户名
	Username string `gorm:"type:varchar(255);column:username" json:"username" validate:"alphanum" description:"用户名"`
	//昵称，如中文名
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `gorm:"type:varchar(255);column:job_number" json:"jobNumber" description:"工号"`
	//默认语言
	Language string `gorm:"type:varchar(255);column:language;default:zh" json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
}

func (ins *UserUpdate) Default(ctx context.Context) {
	if len(ins.Language) == 0 {
		ins.Language = "zh"
	}
}
func (ins *UserUpdate) Validate(ctx context.Context) (err error) {
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

type SetPassword struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//密码
	NewPassword string `json:"newPassword" validate:"required" description:"密码"`
	//旧密码
	OldPassword string ` json:"oldPassword" validate:"required" description:"旧密码"`
}

// UserResetPassword 账户修改密码
type UserResetPassword struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//密码
	Password string `json:"password" validate:"required" description:"密码"`
	//密码强度
	PasswordStrength string `gorm:"type:varchar(20);column:password_strength;default:weak" json:"-" validate:"oneof=strong medium weak" enum:"strong|medium|weak" description:"密码强度"`
	//重置密码状态码
	Code string `json:"code" description:"重置密码状态码"`
}

func (ins *UserResetPassword) Default(ctx context.Context) {
	ins.PasswordStrength = utils.CheckPasswordStrength(ins.Password)
}
func (ins *UserResetPassword) Validate(ctx context.Context) (err error) {
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

type UserMfa struct {
	//主键
	Id  uint `gorm:"column:id" json:"id" description:"用户ID"`
	MFA bool `gorm:"column:mfa" json:"mfa" description:"状态"`
}
