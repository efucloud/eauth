package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// AppAuthRecordDetailList 应用认证记录列表响应
type AppAuthRecordDetailList struct {
	//当前页数据
	Data []*AppAuthRecordDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// AppAuthRecordDetail 应用认证记录详情
type AppAuthRecordDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用ID
	ApplicationId uint `json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code" json:"code" validate:"required" description:"响应编码"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
}

// AppAuthRecordCreate 应用认证记录创建
type AppAuthRecordCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//应用ID
	ApplicationId uint `json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code" json:"code" validate:"required" description:"响应编码"`
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
}

func (ins *AppAuthRecordCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
}
func (ins *AppAuthRecordCreate) Validate(ctx context.Context) (err error) {
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
