/*
Copyright 2022 The itcloudy.com Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var (
	GoVersion string
	Commit    string
	BuildDate string
)
var (
	ApplicationConfig *Config
	configOnce        sync.Once
	DBConnect         *gorm.DB
	Bundle            *i18n.Bundle
	SupperAdminID     string
	Logger            *zap.SugaredLogger
	Verifier          *oidc.IDTokenVerifier
)

func init() {
	if ApplicationConfig == nil {
		configOnce.Do(func() {
			ApplicationConfig = new(Config)
		})
	}

	ContextDBTx = "ContextDBTx"
}
