package daos

import (
	"time"
)

// AppAuthRecord 应用认证记录，该记录在1分钟内有效
type AppAuthRecord struct {
	//主键
	ID string `gorm:"primarykey;column:id;type:varchar(50)" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//应用ID
	ApplicationId string `gorm:"column:application_id;type:varchar(50)" json:"applicationId" description:"应用ID"`
	//响应编码 返回给浏览器客户端
	Code string `gorm:"type:varchar(50);column:code;uniqueIndex" json:"code" validate:"required" description:"响应编码"`
	//PKCE挑战值
	CodeChallenge string `gorm:"type:varchar(255);column:code_challenge" json:"codeChallenge" description:"PKCE挑战值"`
	//PKCE挑战方式
	CodeChallengeMethod string `gorm:"type:varchar(20);column:code_challenge_method" json:"codeChallengeMethod" description:"PKCE挑战方式"`
	//用户ID
	UserId string `gorm:"column:user_id;type:varchar(50);index" json:"userId" validate:"required" description:"用户ID"`
	//OIDC nonce
	Nonce string `gorm:"type:varchar(255);column:nonce" json:"nonce" description:"OIDC nonce"`
}
type AppAuthRecordList struct {
	Data  []*AppAuthRecord `json:"data" description:"数据列表"`
	Total int64            `json:"total" description:"记录总数量"`
}

func (app AppAuthRecord) TableName() string {
	return AppAuthRecordTableName
}
