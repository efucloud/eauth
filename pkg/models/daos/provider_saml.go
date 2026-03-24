package daos

import "time"

type ProviderSaml struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255);uniqueIndex" json:"name" validate:"required" description:"提供商名称"`
	//提供商编码
	Category string `gorm:"type:varchar(255);uniqueIndex" json:"category" validate:"required" description:"提供商编码"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//IdP EntityID
	EntityId string `gorm:"type:varchar(500);column:entity_id" json:"entityId" validate:"required" description:"IdP EntityID"`
	//SAML SSO地址
	SsoURL string `gorm:"type:varchar(1000);column:sso_url" json:"ssoUrl" validate:"required" description:"SAML SSO地址"`
	//回调地址(AssertionConsumerService URL)
	AcsURL string `gorm:"type:varchar(1000);column:acs_url" json:"acsUrl" validate:"required" description:"回调地址"`
	//SAML证书
	Certificate string `gorm:"type:longtext;column:certificate" json:"certificate" validate:"required" description:"SAML证书"`
	//元数据地址
	MetadataURL string `gorm:"type:varchar(1000);column:metadata_url" json:"metadataUrl" description:"元数据地址"`
	//元数据内容
	Metadata string `gorm:"type:longtext;column:metadata" json:"metadata" description:"元数据内容"`
	//登录ID字段映射
	LoginIDAttr string `gorm:"type:varchar(255);column:login_id_attr" json:"loginIdAttr" description:"登录ID字段映射"`
	//登录名字段映射
	LoginNameAttr string `gorm:"type:varchar(255);column:login_name_attr" json:"loginNameAttr" description:"登录名字段映射"`
	//邮箱字段映射
	EmailAttr string `gorm:"type:varchar(255);column:email_attr" json:"emailAttr" description:"邮箱字段映射"`
	//手机号字段映射
	PhoneAttr string `gorm:"type:varchar(255);column:phone_attr" json:"phoneAttr" description:"手机号字段映射"`
	//昵称字段映射
	NicknameAttr string `gorm:"type:varchar(255);column:nickname_attr" json:"nicknameAttr" description:"昵称字段映射"`
	//头像字段映射
	AvatarAttr string `gorm:"type:varchar(255);column:avatar_attr" json:"avatarAttr" description:"头像字段映射"`
}

func (pro ProviderSaml) TableName() string {
	return ProviderSamlTableName
}
