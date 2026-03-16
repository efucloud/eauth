package daos

import (
	"time"

	"github.com/efucloud/eauth/pkg/models/dtos"
)

// UserAuthProfileTemp  临时认证信息
type UserAuthProfileTemp struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//请求码
	Code string `gorm:"type:varchar(50);column:code" json:"code" validate:"required,max=50" description:"请求码"`
	//认证类型provider中的code,email,phone
	Provider string `gorm:"type:varchar(255);column:provider" json:"provider" description:"认证类型"`
	//第三方登录用户的id，邮箱，手机号
	LoginID string `gorm:"type:varchar(255);column:login_id" json:"loginId" description:"第三方认证的用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//第三方认证的全部用户信息
	Properties dtos.JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`
}

func (t UserAuthProfileTemp) TableName() string {
	return UserAuthProfileTempTableName
}
