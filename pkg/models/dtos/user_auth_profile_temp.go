package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/utils"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// UserAuthProfileTempDetailList 系统用户认证方式列表响应
type UserAuthProfileTempDetailList struct {
	//当前页数据
	Data []*UserAuthProfileTempDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// UserAuthProfileTempDetail 系统用户认证方式详情
type UserAuthProfileTempDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
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
	Properties JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`
}

// UserAuthProfileTempCreate 系统用户认证方式创建
type UserAuthProfileTempCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
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
	Properties JsonMap `gorm:"column:properties" json:"properties" description:"第三方认证的全部用户信息"`
}

func (ins *UserAuthProfileTempCreate) Default(ctx context.Context) {
	ins.ID = utils.GenerateDatabaseId()
	ins.CreatedAt = time.Now()
}
func (ins *UserAuthProfileTempCreate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}
