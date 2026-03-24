package daos

import (
	"time"
)

// UserToken 用户生成的Token
// 若使用缓存共享技术，执行token加入黑名单
type UserToken struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//客户端ID
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//过期时间(时间戳)
	Expired int64 `json:"expired" description:" 过期时间(时间戳)"`
	//过期时间
	ExpiredTime time.Time `json:"expiredTime" description:"过期时间"`
	//RefreshToken
	RefreshToken string `gorm:"type:varchar(50);column:refresh_token" json:"refreshToken,omitempty" description:"RefreshToken"`
	//Claims的ID
	ClaimsID string `gorm:"type:varchar(50)" json:"-" validate:"required" description:"Claims的ID"`
	//session key， token MD5
	SessionKey string `gorm:"type:varchar(50)" json:"sessionKey" description:"SessionKey"`
	//Token
	Token string `gorm:"type:longtext" json:"token" description:"Token"`
}

func (t UserToken) TableName() string {
	return UserTokenTableName
}
