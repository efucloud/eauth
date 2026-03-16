package migrations

import (
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/daos"
)

func DatabaseMigrate() {
	var (
		migs []interface{}
	)
	migs = []interface{}{
		&daos.Application{},
		&daos.AppAuthRecord{},
		&daos.UserAuthProfile{},
		&daos.Config{},
		&daos.FaceRecognition{},
		&daos.MultiFactorAuth{},
		&daos.ProviderLdap{},
		&daos.ProviderOidc{},
		&daos.ValidateCode{},
		&daos.User{},
		&daos.UserAuthProfileTemp{},
		&daos.UserToken{},
	}
	err := config.DBConnect.Migrator().AutoMigrate(migs...)
	if err != nil {
		config.Logger.Fatalf("migrate database failed, err: %s", err.Error())
	} else {
		config.Logger.Debugf("migrate database success")
	}

}
