package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type BaiduProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type baiduUserInfo struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Portrait string `json:"portrait"`
	Email    string `json:"email"`
}

var _ IdProvider = &BaiduProvider{}

func NewBaiduProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *BaiduProvider {
	pro := &BaiduProvider{
		ctx:      ctx,
		provider: &provider,
		scopes:   provider.Scopes,
	}
	pro.oauthCfg = &oauth2.Config{
		ClientID:     provider.ClientId,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthorizationEndpoint,
			TokenURL: provider.TokenEndpoint,
		},
		RedirectURL: redirectURL,
		Scopes:      provider.Scopes,
	}
	return pro
}

func (pro *BaiduProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}

func (pro *BaiduProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://openapi.baidu.com/rest/2.0/passport/users/getInfo"
	}
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	requestURL := userInfoEndpoint + "?" + params.Encode()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(requestURL)
	if err != nil {
		errorData.Err = err
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorData.Err = err
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		errorData.Err = fmt.Errorf("request baidu userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfo baiduUserInfo
	if err = jsoniter.Unmarshal(data, &userInfo); err != nil {
		errorData.Err = err
		return
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcBaidu
	profile.LoginID = userInfo.UserID
	profile.Username = userInfo.Username
	profile.Nickname = userInfo.Realname
	if len(profile.Nickname) == 0 {
		profile.Nickname = userInfo.Username
	}
	profile.Email = userInfo.Email
	if len(userInfo.Portrait) > 0 {
		profile.Avatar = "https://himg.bdimg.com/sys/portrait/item/" + userInfo.Portrait
	}
	profile.Home = "https://passport.baidu.com"
	return
}
