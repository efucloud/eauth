package dtos

// DashboardData 看板数据
type DashboardData[T any] struct {
	Name  string `gorm:"name" json:"name" validate:"required" description:"名称"`
	Value T      `gorm:"value" json:"value" validate:"required" description:"值"`
}

// Dashboard 系统统计数据
type Dashboard struct {
	//人脸数据
	FaceRecognition uint `json:"faceRecognition" description:"人脸数据"`
	//认证方式
	AuthProfile []DashboardData[uint] `json:"authProfile" description:"认证方式"`
	//普通应用数据
	Application []DashboardData[uint] `json:"application" description:"普通应用数据"`
	//用户角色
	UserRole []DashboardData[uint] `json:"userRole" description:"用户角色"`
}
type ApplicationAuthTop struct {
	Name          string `json:"name" description:"应用名称"`
	Code          string `json:"code" description:"应用编码"`
	ApplicationId string `gorm:"column:application_id;type:varchar(50)"  json:"applicationId" description:"应用ID"`
	Value         int64  `json:"value" description:"认证数量"`
	Home          string `json:"home" description:"应用主页"`
	Description   string `json:"description" description:"描述"`
	Scope         string `gorm:"type:varchar(50);column:scope;default:tenant" json:"scope" validate:"oneof=system tenant" enum:"system|tenant" description:"应用所属域"`
}
