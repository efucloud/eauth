package daos

import (
	"time"
)

// Application 全局应用
// 使用全局认证的应用，如公司的OA系统
type Application struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用名称
	Name string `gorm:"type:varchar(255);unique" json:"name" validate:"required" description:"应用名称"`
	//应用编码
	Code string `gorm:"type:varchar(255);column:code;unique" json:"code"  validate:"required" description:"应用编码"`
	//应用描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"应用介绍"`
	//应用主页
	Home string `gorm:"type:varchar(1000);column:home" json:"home" description:"应用主页"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true;unique" json:"enable" description:"是否有效"`
	//Logo
	Logo string `gorm:"type:varchar(1000);column:logo" json:"logo" description:"Logo"`
	//OIDC client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255);uniqueIndex" json:"clientId" validate:"required_if=Protocol oidc" description:"client ID"`
	//OIDC 客户端密钥
	ClientSecret string `gorm:"column:client_secret;uniqueIndex" json:"clientSecret,omitempty" validate:"required_if=Protocol oidc" description:"客户端密钥"`
	//OIDC 回调地址
	RedirectUri string `gorm:"type:varchar(1000);column:redirect_uri" json:"redirectUri" validate:"required" description:"回调地址"`
}

func (app Application) TableName() string {
	return ApplicationTableName
}
