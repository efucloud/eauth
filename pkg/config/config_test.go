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
	"os"
	"testing"

	"github.com/ghodss/yaml"
)

func TestConfig(t *testing.T) {
	config := &Config{
		LoginConfig: LoginConfig{
			MFA:             false,
			FaceRecognition: true,
		},
		LogConfig: &LogConfig{
			Filename:   "",
			MaxSize:    0,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  false,
			Compress:   false,
			Production: true,
		},
		TokenPeriod:   16,
		ServerAddress: "http://localhost:8000",
		Mysql: &MysqlConfig{
			Host:                      "localhost:3306",
			User:                      "root",
			Password:                  "EfuCloud",
			Dbname:                    "eauth",
			Charset:                   "utf8",
			Loc:                       "",
			DefaultStringSize:         0,
			DisableDatetimePrecision:  false,
			DontSupportRenameIndex:    false,
			DontSupportRenameColumn:   false,
			SkipInitializeWithVersion: false,
		},
		Email: &EmailConfig{
			Username:   "272685110@qq.com",
			Password:   "ltapxououmwlbggd",
			SmtpServer: "smtp.qq.com",
			SmtpPort:   465,
		},
	}
	data, _ := yaml.Marshal(config)
	_ = os.WriteFile("../../config/config.yaml", data, os.ModePerm)
}
