package daos

import (
	"time"

	"github.com/efucloud/eauth/pkg/models/dtos"
)

// UserAuthProfile 认证信息，会定时清掉OwnerID为空或ApprovalExpire过期的记录，
// 允许用户自己删除 若删除后，再次认证需要重新走所有的流程
type UserAuthProfile struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId string `gorm:"column:user_id;type:varchar(50);index" json:"userId" validate:"required" description:"本系统用户ID"`
	//认证类型provider中的code,email,phone
	Provider string `gorm:"type:varchar(255);column:provider" json:"provider" description:"认证类型"`
	//第三方登录用户的id，邮箱，手机号
	LoginID string `gorm:"type:varchar(255);column:login_id" json:"loginId" description:"第三方认证的用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//第三方认证的全部用户信息
	Properties dtos.JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`
	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}

func (t UserAuthProfile) TableName() string {
	return UserAuthProfileTableName
}
