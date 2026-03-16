package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

var _ IdProvider = &GoogleProvider{}

func NewGoogleProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *GoogleProvider {
	pro := &GoogleProvider{
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

func (pro *GoogleProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}

func (pro *GoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"
	}
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(userInfoEndpoint)
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
		errorData.Err = fmt.Errorf("request google userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfo googleUserInfo
	if err = jsoniter.Unmarshal(data, &userInfo); err != nil {
		errorData.Err = err
		return
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcGoogle
	profile.LoginID = userInfo.Sub
	profile.Username = userInfo.Email
	if len(profile.Username) == 0 {
		profile.Username = userInfo.Name
	}
	profile.Nickname = userInfo.Name
	profile.Email = userInfo.Email
	profile.Avatar = userInfo.Picture
	if len(userInfo.Sub) > 0 {
		profile.Home = "https://profiles.google.com/" + userInfo.Sub
	}
	return
}
