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

// ConfigDetailList 配置列表响应
type ConfigDetailList struct {
	//当前页数据
	Data []*ConfigDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ConfigDetail 配置详情
type ConfigDetail struct {
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

// ConfigCreate 配置创建
type ConfigCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
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

func (ins *ConfigCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	ins.ID = utils.GenerateDatabaseId()
}
func (ins *ConfigCreate) Validate(ctx context.Context) (err error) {
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

// ConfigUpdate 配置修改
type ConfigUpdate struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
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

func (ins *ConfigUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()

}
func (ins *ConfigUpdate) Validate(ctx context.Context) (err error) {
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
