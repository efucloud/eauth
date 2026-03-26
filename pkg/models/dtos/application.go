package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/utils"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
)

// ApplicationDetailList 普通应用列表响应
type ApplicationDetailList struct {
	//当前页数据
	Data []*ShortApplication `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ShortApplication 简单应用信息
type ShortApplication struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"应用名称"`
	//应用编码
	Code string `gorm:"type:varchar(255);column:code" json:"code"  validate:"required" description:"应用编码"`
	//应用描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"应用介绍"`
	//应用主页
	Home string `gorm:"type:varchar(1000);column:home" json:"home" description:"应用主页"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//Logo
	Logo string `gorm:"type:varchar(1000);column:logo" json:"logo" description:"Logo"`
	//OIDC client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required_if=Protocol oidc" description:"client ID"`
	//OIDC 客户端密钥
	ClientSecret string `gorm:"column:client_secret" json:"clientSecret,omitempty" validate:"required_if=Protocol oidc" description:"客户端密钥"`
	//OIDC 回调地址
	RedirectUri string `gorm:"type:varchar(1000);column:redirect_uri" json:"redirectUri" validate:"required" description:"回调地址"`
	//OIDC 重定向地址匹配类型
	RedirectUriMatchType string `gorm:"type:varchar(50);default:equal" json:"redirectUriMatchType"  validate:"oneof=regex equal prefix contain" enum:"regex|equal|prefix|contain"  description:"重定向地址匹配类型。regex:正则，all:全路径，prefix:前缀,contain:包含"`
}

// ApplicationDetail 普通应用详情
type ApplicationDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"应用名称"`
	//应用编码
	Code string `gorm:"type:varchar(255);column:code" json:"code"  validate:"required" description:"应用编码"`
	//应用描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"应用介绍"`
	//应用主页
	Home string `gorm:"type:varchar(1000);column:home" json:"home" description:"应用主页"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//Logo
	Logo string `gorm:"type:varchar(1000);column:logo" json:"logo" description:"Logo"`
	//OIDC client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required_if=Protocol oidc" description:"client ID"`
	//OIDC 客户端密钥
	ClientSecret string `gorm:"column:client_secret" json:"clientSecret,omitempty" validate:"required_if=Protocol oidc" description:"客户端密钥"`
	//OIDC 回调地址
	RedirectUri string `gorm:"type:varchar(1000);column:redirect_uri" json:"redirectUri" validate:"required" description:"回调地址"`
	//OIDC 重定向地址匹配类型
	RedirectUriMatchType string `gorm:"type:varchar(50);default:equal" json:"redirectUriMatchType"  validate:"oneof=regex equal prefix contain" enum:"regex|equal|prefix|contain"  description:"重定向地址匹配类型。regex:正则，all:全路径，prefix:前缀,contain:包含"`
}

func (ins *ApplicationDetail) RedirectUriMatch(redirectUri string) bool {
	switch ins.RedirectUriMatchType {
	case "equal":
		if ins.RedirectUri == redirectUri {
			return true
		}
	case "prefix":
		if strings.HasPrefix(redirectUri, ins.RedirectUri) {
			return true
		}
	case "regex":
		if match, _ := regexp.MatchString(ins.RedirectUri, redirectUri); match {
			return true
		}
	case "contain":
		if strings.Contains(redirectUri, ins.RedirectUri) {
			return true
		}
	}
	return false
}

// ApplicationCreate 普通应用创建
type ApplicationCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//应用名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"应用名称"`
	//应用编码
	Code string `gorm:"type:varchar(255);column:code" json:"code"  validate:"required" description:"应用编码"`
	//应用描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"应用介绍"`
	//应用主页
	Home string `gorm:"type:varchar(1000);column:home" json:"home" description:"应用主页"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//Logo
	Logo string `gorm:"type:varchar(1000);column:logo" json:"logo" description:"Logo"`
	//OIDC client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required_if=Protocol oidc" description:"client ID"`
	//OIDC 客户端密钥
	ClientSecret string `gorm:"column:client_secret" json:"clientSecret,omitempty" validate:"required_if=Protocol oidc" description:"客户端密钥"`
	//OIDC 回调地址
	RedirectUri string `gorm:"type:varchar(1000);column:redirect_uri" json:"redirectUri" validate:"required" description:"回调地址"`
	//OIDC 重定向地址匹配类型
	RedirectUriMatchType string `gorm:"type:varchar(50);default:equal" json:"redirectUriMatchType"  validate:"oneof=regex equal prefix contain" enum:"regex|equal|prefix|contain"  description:"重定向地址匹配类型。regex:正则，all:全路径，prefix:前缀,contain:包含"`
}

func (ins *ApplicationCreate) Default(ctx context.Context) {
	ins.ID = utils.GenerateDatabaseId()
	ins.ClientId = common.NewSecureID(16)
	ins.ClientSecret = common.NewSecureID(32)
	if len(ins.RedirectUriMatchType) == 0 {
		ins.RedirectUriMatchType = "equal"
	}

}
func (ins *ApplicationCreate) Validate(ctx context.Context) (err error) {
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

// ApplicationUpdate 普通应用修改
type ApplicationUpdate struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//应用名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"应用名称"`
	//应用编码
	Code string `gorm:"type:varchar(255);column:code" json:"code"  validate:"required" description:"应用编码"`
	//应用描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"应用介绍"`
	//应用主页
	Home string `gorm:"type:varchar(1000);column:home" json:"home" description:"应用主页"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//Logo
	Logo string `gorm:"type:varchar(1000);column:logo" json:"logo" description:"Logo"`
	//OIDC 回调地址
	RedirectUri string `gorm:"type:varchar(1000);column:redirect_uri" json:"redirectUri" validate:"required" description:"回调地址"`
	//OIDC 重定向地址匹配类型
	RedirectUriMatchType string `gorm:"type:varchar(50);default:equal" json:"redirectUriMatchType"  validate:"oneof=regex equal prefix contain" enum:"regex|equal|prefix|contain"  description:"重定向地址匹配类型。regex:正则，all:全路径，prefix:前缀,contain:包含"`
}

func (ins *ApplicationUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
	if len(ins.RedirectUriMatchType) == 0 {
		ins.RedirectUriMatchType = "equal"
	}
}
func (ins *ApplicationUpdate) Validate(ctx context.Context) (err error) {
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

// ApplicationStatus 普通应用状态
// 状态为disable时将不在用户前端显示，同时普通应用中的应用将不能认证
type ApplicationStatus struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (ins *ApplicationStatus) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *ApplicationStatus) Validate(ctx context.Context) (err error) {
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
