package crons

import (
	"context"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/services"
	"github.com/robfig/cron/v3"
)

func StartCronJob() {
	config.Logger.Info("start cron job")
	c := cron.New()
	go deleteUserTokens()
	_, _ = c.AddFunc("@every 10m", deleteUserTokens)
	c.Start()

}

func deleteUserTokens() {
	svc := services.UserTokenService{}
	svc.DeleteExpireRecord(context.Background())
}
