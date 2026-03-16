package repositories

import (
	"context"
	"github.com/efucloud/eauth/pkg/models/daos"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"gorm.io/gorm"
	"time"
)

type DashboardRepository struct {
	DB *gorm.DB
}

func (repo *DashboardRepository) GetDashboard(ctx context.Context) (dashboard dtos.Dashboard) {
	repo.DB.WithContext(ctx).Table(daos.FaceRecognitionTableName).Select(" COUNT(*) AS value").Scan(&dashboard.FaceRecognition)
	repo.DB.WithContext(ctx).Table(daos.UserAuthProfileTableName).Select("provider as name, COUNT(*) AS value").Where("provider != ?", "").Group("provider").Scan(&dashboard.AuthProfile)
	repo.DB.WithContext(ctx).Table(daos.UserTableName).Select("role as name, COUNT(*) AS value").Group("role").Scan(&dashboard.UserRole)
	return
}
func (repo *DashboardRepository) ApplicationAuthTop10(ctx context.Context) (dashboard []dtos.ApplicationAuthTop) {
	repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Select("application_id, COUNT(*) AS value").Group("application_id").Limit(10).Scan(&dashboard)
	return
}
func (repo *DashboardRepository) ApplicationAuth30Days(ctx context.Context) (dashboard []dtos.DashboardData[int64]) {
	repo.DB.WithContext(ctx).Table(daos.AppAuthRecordTableName).Select("DATE_FORMAT(created_at, '%Y-%m-%d') name, count(*) value").Where("created_at > ?", time.Now().Add(-30*24*time.Hour)).Group("name").Scan(&dashboard)
	return
}
