package daos

import (
	"time"
)

// User 用户表
type User struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户名
	Username string `gorm:"type:varchar(255);column:username;uniqueIndex" json:"username" validate:"alphanum" description:"用户名"`
	//昵称，如中文名
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `gorm:"type:varchar(255);column:job_number" json:"jobNumber" description:"工号"`
	//密码
	Password string `gorm:"-" json:"password" example:"admin" description:"密码"`
	//数据库保存的加密密码
	PasswordStore string `gorm:"type:varchar(255);column:password_store" json:"passwordStore"`
	//密码强度
	PasswordStrength string `gorm:"type:varchar(255);column:password_strength;default:weak" json:"passwordStrength" validate:"oneof=strong medium weak" enum:"strong|medium|weak" description:"密码强度"`
	//系统角色
	Role string `gorm:"type:varchar(255);column:role;default:none" json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//邮箱
	Email string `gorm:"type:varchar(255);column:email;uniqueIndex" json:"email" validate:"email" description:"邮箱"`
	//手机号码
	Phone string `gorm:"type:varchar(255);column:phone;uniqueIndex" json:"phone" validate:"required" description:"电话"`
	//默认语言
	Language string `gorm:"type:varchar(255);column:language" json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
	//头像
	Avatar string `gorm:"type:varchar(1000);column:avatar" json:"avatar" description:"头像"`
	//是否绑定MFA认证
	MFA bool `gorm:"column:mfa" json:"mfa" description:"是否绑定MFA认证"`
}

func (t User) TableName() string {
	return UserTableName
}
