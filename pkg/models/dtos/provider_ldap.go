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

// ProviderLdapDetailList LDAP提供商列表响应
type ProviderLdapDetailList struct {
	//当前页数据
	Data []*ProviderLdapDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ProviderLdapDetail LDAP提供商详情
type ProviderLdapDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"名称"`
	//图标
	Icon string `gorm:"type:varchar(255)" json:"icon" validate:"required" description:"图标"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//认证提供商
	Provider string `gorm:"type:varchar(255)" json:"provider" validate:"oneof=ad openldap"  enum:"ad|openldap" description:"认证提供商"`

	//LDAP配置
	LdapConfig *LdapConfig `json:"ldapConfig" description:"LDAP配置"`
}

// ProviderLdapCreate LDAP提供商创建
type ProviderLdapCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"名称"`
	//编码
	Code string `gorm:"type:varchar(255)" json:"code" validate:"required" description:"编码"` //编码
	//提供商类型
	Category string `gorm:"type:varchar(255)" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"提供商类型"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//认证提供商
	Provider string `gorm:"type:varchar(255)" json:"provider" validate:"oneof=ad openldap"  enum:"ad|openldap" description:"认证提供商"`

	//LDAP配置
	LdapConfig *LdapConfig `json:"ldapConfig" description:"LDAP配置"`
}

func (ins *ProviderLdapCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	ins.ID = utils.GenerateDatabaseId()

}
func (ins *ProviderLdapCreate) Validate(ctx context.Context) (err error) {
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

// ProviderLdapUpdate LDAP提供商修改
type ProviderLdapUpdate struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"名称"`
	//编码
	Code     string `gorm:"type:varchar(255)" json:"code" validate:"required" description:"编码"` //图标
	Category string `gorm:"type:varchar(255)" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"图标"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//认证提供商
	Provider string `gorm:"type:varchar(255)" json:"provider" validate:"oneof=ad openldap"  enum:"ad|openldap" description:"认证提供商"`
	//LDAP配置
	LdapConfig *LdapConfig `json:"ldapConfig" description:"LDAP配置"`
}

func (ins *ProviderLdapUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()

}
func (ins *ProviderLdapUpdate) Validate(ctx context.Context) (err error) {
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
