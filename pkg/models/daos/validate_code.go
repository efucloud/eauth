package daos

import (
	"time"
)

type ValidateCode struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" description:"用户ID"`
	//验证类型
	Category string `gorm:"type:varchar(50);column:category" json:"category" validate:"oneof=phone email" enum:"phone|email" description:"验证类型"`
	//验证码
	Code string `gorm:"type:varchar(50);column:code" json:"code" description:"验证码"`
	//动作
	Action string `gorm:"type:varchar(50);column:action" json:"action" validate:"oneof=registry login changepwd forgetpwd"  enum:"registry|login|changepwd|forgetpwd" description:"动作"`
	//过期时间
	Expired time.Time `gorm:"column:expired" json:"expired" validate:"required" description:"过期时间"`
}

func (pro ValidateCode) TableName() string {
	return ValidateCodeTableName
}
