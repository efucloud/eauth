package daos

import (
	"time"

	"github.com/efucloud/eauth/pkg/models/dtos"
)

// FaceRecognition 用户人脸识别信息
type FaceRecognition struct {
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
	FaceIdData dtos.ArrayFloat64 `gorm:"column:face_id_data" json:"faceIdData" description:"人脸数据"`
}

func (t FaceRecognition) TableName() string {
	return FaceRecognitionTableName
}
