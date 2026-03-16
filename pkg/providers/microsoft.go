package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type MicrosoftProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type microsoftUserInfo struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
}

var _ IdProvider = &MicrosoftProvider{}

func NewMicrosoftProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *MicrosoftProvider {
	pro := &MicrosoftProvider{
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

func (pro *MicrosoftProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}

func (pro *MicrosoftProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://graph.microsoft.com/v1.0/me"
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
		errorData.Err = fmt.Errorf("request microsoft userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfo microsoftUserInfo
	if err = jsoniter.Unmarshal(data, &userInfo); err != nil {
		errorData.Err = err
		return
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcMicrosoft
	profile.LoginID = userInfo.ID
	profile.Nickname = userInfo.DisplayName
	profile.Username = userInfo.UserPrincipalName
	if len(profile.Username) == 0 {
		profile.Username = userInfo.Mail
	}
	profile.Email = userInfo.Mail
	if len(profile.Email) == 0 && strings.Contains(userInfo.UserPrincipalName, "@") {
		profile.Email = userInfo.UserPrincipalName
	}
	profile.Home = "https://portal.office.com"
	return
}
