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

const (
	V1Prefix        = "/api"
	ServerPort      = 9001
	RequestForm     = "application/x-www-form-urlencoded"
	RequestFormData = "multipart/form-data"
)
const (
	ProviderOidcGitlab     = "gitlab"
	ProviderOidcGithub     = "github"
	ProviderOidcGoogle     = "google"
	ProviderOidcMicrosoft  = "microsoft"
	ProviderOidcWechat     = "wechat"
	ProviderOidcAlipay     = "alipay"
	ProviderOidcWeibo      = "weibo"
	ProviderOidcWechatWork = "wechatWork"
	ProviderOidcQq         = "qq"
	ProviderOidcDingTalk   = "dingTalk"
	ProviderOidcGitee      = "gitee"
	ProviderOidcFeiShu     = "feiShu"
	ProviderOidcBilibili   = "bilibili"
	ProviderOidcTiktok     = "tiktok"
	ProviderOidcBaidu      = "baidu"
	ProviderOidcCustom     = "custom"
)

const (
	AdminUsername = "admin"
	AdminPassword = "EfuCloud"
	RoleAdmin     = "admin"
	RoleView      = "view"
	RoleEdit      = "edit"
	RoleNone      = "none"
	FrontApiTag   = "FrontApiTag"
	UserAvatars   = "user-avatars"
)
const (
	PasswordSalt = "Cdq3rsd12&"

	ApplicationName = "eauth"
)
const (
	RequestUserID   = "RequestUserID"
	RequestLanguage = "RequestLanguage"
	RequestContext  = "RequestContext"
	RequestNonce    = "RequestNonce"
	AuthHeader      = "Authorization"
)

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeToken             = "token"
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypePassword          = "password"
	GrantTypeImplicit          = "implicit"
	GrantTypeDeviceCode        = "urn:ietf:params:oauth:grant-type:device_code"
)

// 系统配置的key
const (
	JwtPublicKey  = "JwtPublicKey"
	JwtPrivateKey = "JwtPrivateKey"
)

type ContextDatabaseTx string

var ContextDBTx ContextDatabaseTx
