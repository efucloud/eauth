package daos

import (
	"github.com/efucloud/eauth/pkg/models/dtos"
	"time"
)

type ProviderLdap struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//名称
	Name string `gorm:"type:varchar(255);uniqueIndex" json:"name"  validate:"required" description:"名称"`
	//编码
	Code string `gorm:"type:varchar(255);uniqueIndex" json:"code" validate:"required" description:"编码"` //图标
	//提供商分类
	Category string `gorm:"type:varchar(255);uniqueIndex" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"提供商分类"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//认证提供商
	Provider string `gorm:"type:varchar(255)" json:"provider" validate:"oneof=ad openldap"  enum:"ad|openldap" description:"认证提供商"`
	//LDAP配置
	LdapConfig *dtos.LdapConfig `json:"ldapConfig" description:"LDAP配置"`
}

func (ld ProviderLdap) TableName() string {
	return ProviderLdapTableName
}
