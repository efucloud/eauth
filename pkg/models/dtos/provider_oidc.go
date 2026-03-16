package dtos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// ProviderOidcDetailList OIDC提供商列表响应
type ProviderOidcDetailList struct {
	//当前页数据
	Data []*ProviderOidcDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ProviderOidcDetail OIDC提供商详情
type ProviderOidcDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"提供商名称"`
	//提供商类型，根据该字段显示不同的图标
	Category string `gorm:"type:varchar(255)" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"提供商类型"`
	//client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//ClientSecret，在提供商创建应用时生成
	ClientSecret string `gorm:"type:varchar(255)" json:"clientSecret,omitempty" validate:"required" description:"client Secret"`
	//颁发者地址
	Issuer string `gorm:"type:varchar(500);column:issuer" json:"issuer" description:"颁发者地址"`
	//作用域
	Scopes ArrayString `json:"scopes" description:"作用域"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//授权地址
	AuthorizationEndpoint string `gorm:"type:varchar(255)" json:"authorizationEndpoint" validate:"required" description:"授权地址"`
	//令牌获取地址
	TokenEndpoint string `gorm:"type:varchar(255)" json:"tokenEndpoint" validate:"required" description:"令牌获取地址"`
	//用户信息获取地址
	UserinfoEndpoint string `gorm:"type:varchar(255)" json:"userinfoEndpoint" description:"用户信息获取地址"`
}

// ProviderOidcCreate OIDC提供商创建
type ProviderOidcCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"提供商名称"`
	//图标
	Category string `gorm:"type:varchar(255)" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"图标"`
	//client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//ClientSecret，在提供商创建应用时生成
	ClientSecret string `gorm:"type:varchar(255)" json:"clientSecret,omitempty" validate:"required" description:"client Secret"`
	//颁发者地址
	Issuer string `gorm:"type:varchar(500);column:issuer" json:"issuer" description:"颁发者地址"`
	//作用域
	Scopes ArrayString `json:"scopes" description:"作用域"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//授权地址
	AuthorizationEndpoint string `gorm:"type:varchar(255)" json:"authorizationEndpoint" validate:"required" description:"授权地址"`
	//令牌获取地址
	TokenEndpoint string `gorm:"type:varchar(255)" json:"tokenEndpoint" validate:"required" description:"令牌获取地址"`
	//用户信息获取地址
	UserinfoEndpoint string `gorm:"type:varchar(255)" json:"userinfoEndpoint" description:"用户信息获取地址"`
}

func (ins *ProviderOidcCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	switch ins.Category {
	case config.ProviderOidcGitlab:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://gitlab.com/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://gitlab.com/oauth/token"
		}
		if len(ins.Issuer) == 0 {
			ins.Issuer = "https://gitlab.com"
		} else {
			ins.Issuer = strings.TrimSuffix(ins.Issuer, "/")
		}
		if !common.StringInArray("read_user", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "read_user")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
	case config.ProviderOidcGithub:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://github.com/login/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://github.com/login/oauth/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.github.com/user"
		}
		if !common.StringInArray("user:email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "user:email")
		}
		if !common.StringInArray("read:user", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "read:user")
		}
	case config.ProviderOidcGoogle:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://oauth2.googleapis.com/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"
		}
		if !common.StringInArray("openid", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "openid")
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
	case config.ProviderOidcMicrosoft:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://graph.microsoft.com/v1.0/me"
		}
		if !common.StringInArray("openid", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "openid")
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
		if !common.StringInArray("User.Read", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "User.Read")
		}
	case config.ProviderOidcWechat:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.weixin.qq.com/connect/qrconnect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.weixin.qq.com/sns/oauth2/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.weixin.qq.com/sns/userinfo"
		}

		if !common.StringInArray("snsapi_login", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "snsapi_login")
		}
	case config.ProviderOidcAlipay:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://openapi.alipay.com/gateway.do"
		}
	case config.ProviderOidcWeibo:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://api.weibo.com/oauth2/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.weibo.com/oauth2/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.weibo.com/2/users/show.json"
		}

	case config.ProviderOidcWechatWork:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
		}
		if !common.StringInArray("snsapi_login", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "snsapi_login")
		}
	case config.ProviderOidcQq:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://graph.qq.com/oauth2.0/show?which=Login&display=pc"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://graph.qq.com/oauth2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://graph.qq.com/user/get_user_info"
		}

		if !common.StringInArray("get_user_info", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "get_user_info")
		}
	case config.ProviderOidcDingTalk:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://login.dingtalk.com/oauth2/auth"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.dingtalk.com/v1.0/contact/users/me"
		}

	case config.ProviderOidcGitee:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://gitee.com/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://gitee.com/oauth/token"
		}
		if len(ins.Issuer) == 0 {
			ins.Issuer = "https://gitee.com"
		} else {
			ins.Issuer = strings.TrimSuffix(ins.Issuer, "/")
		}
	case config.ProviderOidcFeiShu:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://accounts.feishu.cn/open-apis/authen/v1/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://open.feishu.cn/open-apis/authen/v2/oauth/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://open.feishu.cn/open-apis/authen/v1/user_info"
		}
		if !common.StringInArray("contact:user.base:readonly", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "contact:user.base:readonly")
		}
	case config.ProviderOidcBilibili:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://passport.bilibili.com/register/pc_oauth2.html"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.bilibili.com/x/account-oauth2/v1/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://member.bilibili.com/arcopen/fn/user/account/info"
		}
	case config.ProviderOidcTiktok:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.douyin.com/platform/oauth/connect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://open.douyin.com/oauth/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://open.douyin.com/oauth/userinfo/"
		}
	case config.ProviderOidcBaidu:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://openapi.baidu.com/oauth/2.0/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://openapi.baidu.com/oauth/2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://openapi.baidu.com/rest/2.0/passport/users/getInfo"
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
	case config.ProviderOidcCustom:

	}
}
func (ins *ProviderOidcCreate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}

// ProviderOidcUpdate OIDC提供商修改
type ProviderOidcUpdate struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name"  validate:"required" description:"提供商名称"`
	//图标
	Category string `gorm:"type:varchar(255)" json:"category" validate:"oneof=gitlab github google microsoft wechat alipay taobao weibo wechatWork qq dingTalk gitee feiShu bilibili tiktok baidu custom" description:"图标"`
	//client ID，在提供商创建应用时生成
	ClientId string `gorm:"type:varchar(255)" json:"clientId" validate:"required" description:"client ID"`
	//ClientSecret，在提供商创建应用时生成
	ClientSecret string `gorm:"type:varchar(255)" json:"clientSecret,omitempty" validate:"required" description:"client Secret"`
	//颁发者地址
	Issuer string `gorm:"type:varchar(500);column:issuer" json:"issuer" description:"颁发者地址"`
	//作用域
	Scopes ArrayString `json:"scopes" description:"作用域"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//授权地址
	AuthorizationEndpoint string `gorm:"type:varchar(255)" json:"authorizationEndpoint" validate:"required" description:"授权地址"`
	//令牌获取地址
	TokenEndpoint string `gorm:"type:varchar(255)" json:"tokenEndpoint" validate:"required" description:"令牌获取地址"`
	//用户信息获取地址
	UserinfoEndpoint string `gorm:"type:varchar(255)" json:"userinfoEndpoint" description:"用户信息获取地址"`
}

func (ins *ProviderOidcUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
	switch ins.Category {
	case config.ProviderOidcGitlab:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://gitlab.com/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://gitlab.com/oauth/token"
		}
		if len(ins.Issuer) == 0 {
			ins.Issuer = "https://gitlab.com"
		} else {
			ins.Issuer = strings.TrimSuffix(ins.Issuer, "/")
		}
		if !common.StringInArray("read_user", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "read_user")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
	case config.ProviderOidcGithub:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://github.com/login/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://github.com/login/oauth/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.github.com/user"
		}
		if !common.StringInArray("user:email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "user:email")
		}
		if !common.StringInArray("read:user", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "read:user")
		}
	case config.ProviderOidcGoogle:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://oauth2.googleapis.com/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"
		}
		if !common.StringInArray("openid", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "openid")
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
	case config.ProviderOidcMicrosoft:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://graph.microsoft.com/v1.0/me"
		}
		if !common.StringInArray("openid", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "openid")
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
		if !common.StringInArray("profile", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "profile")
		}
		if !common.StringInArray("User.Read", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "User.Read")
		}
	case config.ProviderOidcWechat:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.weixin.qq.com/connect/qrconnect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.weixin.qq.com/sns/oauth2/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.weixin.qq.com/sns/userinfo"
		}

		if !common.StringInArray("snsapi_login", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "snsapi_login")
		}
	case config.ProviderOidcAlipay:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://openapi.alipay.com/gateway.do"
		}
	case config.ProviderOidcWeibo:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://api.weibo.com/oauth2/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.weibo.com/oauth2/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.weibo.com/2/users/show.json"
		}

	case config.ProviderOidcWechatWork:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
		}
		if !common.StringInArray("snsapi_login", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "snsapi_login")
		}
	case config.ProviderOidcQq:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://graph.qq.com/oauth2.0/show?which=Login&display=pc"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://graph.qq.com/oauth2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://graph.qq.com/user/get_user_info"
		}

		if !common.StringInArray("get_user_info", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "get_user_info")
		}
	case config.ProviderOidcDingTalk:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://login.dingtalk.com/oauth2/auth"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://api.dingtalk.com/v1.0/contact/users/me"
		}

	case config.ProviderOidcGitee:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://gitee.com/oauth/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://gitee.com/oauth/token"
		}
		if len(ins.Issuer) == 0 {
			ins.Issuer = "https://gitee.com"
		} else {
			ins.Issuer = strings.TrimSuffix(ins.Issuer, "/")
		}
	case config.ProviderOidcFeiShu:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://accounts.feishu.cn/open-apis/authen/v1/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://open.feishu.cn/open-apis/authen/v2/oauth/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://open.feishu.cn/open-apis/authen/v1/user_info"
		}
		if !common.StringInArray("contact:user.base:readonly", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "contact:user.base:readonly")
		}
	case config.ProviderOidcBilibili:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://passport.bilibili.com/register/pc_oauth2.html"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://api.bilibili.com/x/account-oauth2/v1/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://member.bilibili.com/arcopen/fn/user/account/info"
		}
	case config.ProviderOidcTiktok:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://open.douyin.com/platform/oauth/connect"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://open.douyin.com/oauth/access_token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://open.douyin.com/oauth/userinfo/"
		}
	case config.ProviderOidcBaidu:
		if len(ins.AuthorizationEndpoint) == 0 {
			ins.AuthorizationEndpoint = "https://openapi.baidu.com/oauth/2.0/authorize"
		}
		if len(ins.TokenEndpoint) == 0 {
			ins.TokenEndpoint = "https://openapi.baidu.com/oauth/2.0/token"
		}
		if len(ins.UserinfoEndpoint) == 0 {
			ins.UserinfoEndpoint = "https://openapi.baidu.com/rest/2.0/passport/users/getInfo"
		}
		if !common.StringInArray("email", ins.Scopes) {
			ins.Scopes = append(ins.Scopes, "email")
		}
	case config.ProviderOidcCustom:

	}
}
func (ins *ProviderOidcUpdate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}

// ProviderOidcStatus 认证提供商状态
// 状态为disable时将不在用户前端显示
type ProviderOidcStatus struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (t *ProviderOidcStatus) Default(ctx context.Context) {
	t.UpdatedAt = time.Now()
}
func (t *ProviderOidcStatus) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(t)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}
