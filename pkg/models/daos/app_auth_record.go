package daos

import (
	"time"
)

// AppAuthRecord 应用认证记录，该记录在1分钟内有效
type AppAuthRecord struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用ID
	ApplicationId uint `json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code;uniqueIndex" json:"code" validate:"required" description:"响应编码"`
	//用户ID
	UserId uint `gorm:"column:user_id;index" json:"userId" validate:"required" description:"用户ID"`
}
type AppAuthRecordList struct {
	Data  []*AppAuthRecord `json:"data" description:"数据列表"`
	Total int64            `json:"total" description:"记录总数量"`
}

func (app AppAuthRecord) TableName() string {
	return AppAuthRecordTableName
}
