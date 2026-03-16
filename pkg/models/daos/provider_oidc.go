package daos

import (
	"github.com/efucloud/eauth/pkg/models/dtos"
	"time"
)

type ProviderOidc struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255);uniqueIndex" json:"name"  validate:"required" description:"提供商名称"`
	//图标
	Category string `gorm:"type:varchar(255);uniqueIndex" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"图标"`
	//client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//ClientSecret，在提供商创建应用时生成
	ClientSecret string `gorm:"type:varchar(255)" json:"clientSecret,omitempty" validate:"required" description:"client Secret"`
	//颁发者地址
	Issuer string `gorm:"type:varchar(500);column:issuer" json:"issuer" description:"颁发者地址"`
	//作用域
	Scopes dtos.ArrayString `json:"scopes" validate:"required" description:"作用域"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//授权地址
	AuthorizationEndpoint string `gorm:"type:varchar(255)" json:"authorizationEndpoint" validate:"required" description:"授权地址"`
	//令牌获取地址
	TokenEndpoint string `gorm:"type:varchar(255)" json:"tokenEndpoint" validate:"required" description:"令牌获取地址"`
	//用户信息获取地址
	UserinfoEndpoint string `gorm:"type:varchar(255)" json:"userinfoEndpoint" description:"用户信息获取地址"`
}

func (pro ProviderOidc) TableName() string {
	return ProviderOidcTableName
}
