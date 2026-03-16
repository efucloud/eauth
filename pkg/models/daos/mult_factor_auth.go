package daos

import (
	"time"
)

// MultiFactorAuth 多因子认证
type MultiFactorAuth struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//所属用户
	UserId uint `gorm:"user_id;uniqueIndex" json:"userId" validate:"required" description:"所属用户"`
	//密钥
	Secret string `gorm:"type:longtext;column:secret" json:"secret" validate:"required" description:"密钥"`
	//二维码
	Image string `gorm:"type:longtext;column:image" json:"image" validate:"required" description:"二维码"`
	//状态：是否已绑定
	Status string `gorm:"type:varchar(50);column:status;default:unbound" json:"status" validate:"oneof=bound unbound" enum:"bound|unbound"  description:"状态"`
}

func (t MultiFactorAuth) TableName() string {
	return MultiFactorAuthTableName
}
