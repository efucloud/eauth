package daos

import (
	"time"
)

// Config 全局配置
type Config struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//配置名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"配置名称"`
	//配置编码
	Code string `gorm:"type:varchar(255);uniqueIndex" json:"code" validate:"alpha,ne=eauth" description:"配置编码"`
	//描述
	Description string `gorm:"type:varchar(255)" json:"description" description:"描述"`
	//值
	Value string `gorm:"type:longtext" json:"value" validate:"required" description:"值"`
	//值类型
	ValType string `gorm:"type:varchar(20)" json:"valType"  validate:"oneof=string number" enum:"string|number" description:"值类型"`
}

func (sys Config) TableName() string {
	return ConfigTableName
}
