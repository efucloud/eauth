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

// FaceRecognitionDetailList 用户人脸识别信息
type FaceRecognitionDetailList struct {
	//当前页数据
	Data []*FaceRecognitionDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// FaceRecognitionDetail 用户人脸识别信息详情
type FaceRecognitionDetail struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//用户ID
	UserId string `gorm:"column:user_id;type:varchar(50);index" json:"userId" validate:"required" description:"本系统用户ID"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//录入名称
	Name string `gorm:"type:varchar(50);column:name" json:"name" validate:"required,max=50" description:"录入名称"`
	//人脸数据
	FaceIdData ArrayFloat64 `gorm:"column:face_id_data" json:"-" description:"人脸数据"`

	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}

// FaceRecognitionCreate 用户人脸识别信息创建
type FaceRecognitionCreate struct {
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"-" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//用户ID
	UserId string `gorm:"column:user_id;type:varchar(50);index" json:"userId" validate:"required" description:"本系统用户ID"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//录入名称
	Name string `gorm:"type:varchar(50);column:name" json:"name" validate:"required,max=50" description:"录入名称"`
	//人脸数据
	FaceIdData ArrayFloat64 `gorm:"column:face_id_data" json:"faceIdData" description:"人脸数据"`
}

func (ins *FaceRecognitionCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	ins.ID = utils.GenerateDatabaseId()
}
func (ins *FaceRecognitionCreate) Validate(ctx context.Context) (err error) {
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

// FaceRecognitionStatus 禁用后，用户将不能使用该认证方式登陆系统
type FaceRecognitionStatus struct {
	//主键
	Ids []string `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (ins *FaceRecognitionStatus) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
}
func (ins *FaceRecognitionStatus) Validate(ctx context.Context) (err error) {
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
