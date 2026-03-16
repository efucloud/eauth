package services

import (
	"context"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/repositories"
	"gorm.io/gorm"
)

type DashboardService struct {
	repo repositories.DashboardRepository
}

func (svc *DashboardService) init(ctx context.Context) {
	db, ok := ctx.Value(config.ContextDBTx).(*gorm.DB)
	if ok {
		svc.repo = repositories.DashboardRepository{DB: db}
	} else {
		svc.repo = repositories.DashboardRepository{DB: config.DBConnect}
	}
}
func (svc *DashboardService) Dashboard(ctx context.Context) (result dtos.Dashboard) {
	svc.init(ctx)
	return svc.repo.GetDashboard(ctx)
}
func (svc *DashboardService) ApplicationAuthTop10(ctx context.Context) (result []dtos.ApplicationAuthTop) {
	svc.init(ctx)
	result = svc.repo.ApplicationAuthTop10(ctx)
	appSvc := ApplicationService{}
	for i, _ := range result {
		app, _ := appSvc.GetApplicationById(ctx, result[i].ApplicationId)
		result[i].Name = app.Name
		result[i].Code = app.Code
		result[i].Home = app.Home
		result[i].Description = app.Description
	}
	return
}

func (svc *DashboardService) ApplicationAuth30Days(ctx context.Context) (dashboard []dtos.DashboardData[int64]) {
	svc.init(ctx)
	return svc.repo.ApplicationAuth30Days(ctx)
}
