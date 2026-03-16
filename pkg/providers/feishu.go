package providers

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type feiShuTokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
	} `json:"data"`
}

type feiShuUserResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OpenID    string `json:"open_id"`
		UnionID   string `json:"union_id"`
		Name      string `json:"name"`
		EnName    string `json:"en_name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
	} `json:"data"`
}

type FeiShuProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

var _ IdProvider = &FeiShuProvider{}

func NewFeiShuProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *FeiShuProvider {
	pro := &FeiShuProvider{
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

func (pro *FeiShuProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	requestData := struct {
		GrantType    string `json:"grant_type"`
		Code         string `json:"code"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
	}{
		GrantType:    "authorization_code",
		Code:         code,
		ClientID:     pro.provider.ClientId,
		ClientSecret: pro.provider.ClientSecret,
		RedirectURI:  pro.oauthCfg.RedirectURL,
	}
	data, err := pro.postWithBody(requestData, pro.provider.TokenEndpoint)
	if err != nil {
		errorData.Err = err
		return
	}
	var tokenResp feiShuTokenResponse
	if err = jsoniter.Unmarshal(data, &tokenResp); err != nil {
		errorData.Err = err
		return
	}
	if tokenResp.Code != 0 {
		errorData.Err = fmt.Errorf("request feishu token failed, code: %d, msg: %s", tokenResp.Code, tokenResp.Msg)
		return
	}
	if len(tokenResp.Data.AccessToken) == 0 {
		errorData.Err = fmt.Errorf("request feishu token failed, empty access_token")
		return
	}
	token = &oauth2.Token{
		AccessToken:  tokenResp.Data.AccessToken,
		RefreshToken: tokenResp.Data.RefreshToken,
		TokenType:    tokenResp.Data.TokenType,
	}
	if tokenResp.Data.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(tokenResp.Data.ExpiresIn) * time.Second)
	}
	return
}

func (pro *FeiShuProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://open.feishu.cn/open-apis/authen/v1/user_info"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoEndpoint, nil)
	if err != nil {
		errorData.Err = err
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
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
		errorData.Err = fmt.Errorf("request feishu userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfoResp feiShuUserResponse
	if err = jsoniter.Unmarshal(data, &userInfoResp); err != nil {
		errorData.Err = err
		return
	}
	if userInfoResp.Code != 0 {
		errorData.Err = fmt.Errorf("request feishu userinfo failed, code: %d, msg: %s", userInfoResp.Code, userInfoResp.Msg)
		return
	}
	profile.Provider = config.ProviderOidcFeiShu
	profile.LoginID = userInfoResp.Data.UnionID
	if len(profile.LoginID) == 0 {
		profile.LoginID = userInfoResp.Data.OpenID
	}
	profile.Username = userInfoResp.Data.EnName
	if len(profile.Username) == 0 {
		profile.Username = userInfoResp.Data.Name
	}
	profile.Nickname = userInfoResp.Data.Name
	profile.Email = userInfoResp.Data.Email
	profile.Phone = userInfoResp.Data.Mobile
	profile.Avatar = userInfoResp.Data.AvatarURL
	profile.Home = "https://www.feishu.cn"
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	return
}

func (pro *FeiShuProvider) postWithBody(body interface{}, requestURL string) ([]byte, error) {
	bs, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(bs))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(requestURL, "application/json;charset=UTF-8", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request feishu endpoint failed, status: %d, body: %s", resp.StatusCode, string(data))
	}
	return data, nil
}
