package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"net/http"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/common/signals"
	"github.com/efucloud/eauth/cmd/server/options"
	"github.com/efucloud/eauth/pkg/apis"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/crons"
	"github.com/efucloud/eauth/pkg/embeds"
	"github.com/efucloud/eauth/pkg/migrations"
	"github.com/efucloud/eauth/pkg/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func NewRunnerServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:          "server",
		Long:         `eauth server`,
		Short:        "eauth server",
		Example:      `eauth server`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return run(s, signals.SetupSignalHandler())
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&s.Config, "config", "c", "./config/config.yaml", "config file path")
	return cmd
}
func run(o *options.ServerRunOptions, stopCh <-chan struct{}) error {
	common.LoadConfig(o.Config, config.ApplicationConfig)
	config.ApplicationConfig.Init()

	config.Bundle, _ = common.I18nInit(embeds.I18nFiles, config.Logger)
	ctx := context.Background()
	//数据库表创建
	migrations.DatabaseMigrate()
	sys := services.ConfigService{}
	sys.InitConfig(ctx)

	userSvc := services.UserService{}
	user, _ := userSvc.GetUserByUsername(ctx, config.AdminUsername)
	config.SupperAdminID = user.ID
	if config.SupperAdminID == 0 {
		var userCreate dtos.UserCreate
		userCreate.Username = config.AdminUsername
		userCreate.Nickname = "系统管理员"
		userCreate.Email = "admin@efucloud.com"
		userCreate.Phone = "13988888888"
		userCreate.Enable = true
		userCreate.Password = config.AdminPassword
		userCreate.Role = config.RoleAdmin
		var errData common.ErrorData
		user, errData = userSvc.AddUser(ctx, userCreate)
		if errData.IsNotNil() {
			config.Logger.Fatalf("can't get or create admin use, err: %s", errData.Err.Error())
		} else {
			config.SupperAdminID = user.ID
		}
	}
	apis.AddResources()
	apis.AddSwagger()
	go crons.StartCronJob()
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.ServerPort),
	}
	errCh := make(chan error, 1)
	config.Logger.Infof("ready to start server with http and port is: %d", config.ServerPort)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-stopCh:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			config.Logger.Errorf("server shutdown failed, err: %s", err.Error())
			return err
		}
		config.Logger.Info("server shutdown completed")
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		if err != nil {
			config.Logger.Errorf("ready to start server failed, err: %s", err.Error())
		}
		return err
	}
}
