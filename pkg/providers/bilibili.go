package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type BilibiliProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type bilibiliTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Mid          int64  `json:"mid"`
}

type bilibiliTokenResponse struct {
	Code         int               `json:"code"`
	Message      string            `json:"message"`
	Msg          string            `json:"msg"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int64             `json:"expires_in"`
	Mid          int64             `json:"mid"`
	Data         bilibiliTokenData `json:"data"`
}

type bilibiliUserInfo struct {
	Mid      int64  `json:"mid"`
	Uname    string `json:"uname"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Face     string `json:"face"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	SpaceURL string `json:"space_url"`
}

type bilibiliUserResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Msg     string           `json:"msg"`
	Data    bilibiliUserInfo `json:"data"`
	Mid     int64            `json:"mid"`
	Uname   string           `json:"uname"`
	Name    string           `json:"name"`
	Face    string           `json:"face"`
}

var _ IdProvider = &BilibiliProvider{}

func NewBilibiliProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *BilibiliProvider {
	pro := &BilibiliProvider{
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

func (pro *BilibiliProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Set("client_id", pro.provider.ClientId)
	formData.Set("client_secret", pro.provider.ClientSecret)
	formData.Set("appkey", pro.provider.ClientId)
	formData.Set("appsecret", pro.provider.ClientSecret)
	formData.Set("redirect_uri", pro.oauthCfg.RedirectURL)
	resp, err := http.PostForm(pro.provider.TokenEndpoint, formData)
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
		errorData.Err = fmt.Errorf("request bilibili token failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var tokenResp bilibiliTokenResponse
	if err = jsoniter.Unmarshal(data, &tokenResp); err != nil {
		errorData.Err = err
		return
	}
	if tokenResp.Code != 0 {
		errorData.Err = fmt.Errorf("request bilibili token failed, code: %d, message: %s %s", tokenResp.Code, tokenResp.Message, tokenResp.Msg)
		return
	}
	accessToken := tokenResp.AccessToken
	if len(accessToken) == 0 {
		accessToken = tokenResp.Data.AccessToken
	}
	if len(accessToken) == 0 {
		errorData.Err = fmt.Errorf("request bilibili token failed, empty access_token")
		return
	}
	refreshToken := tokenResp.RefreshToken
	if len(refreshToken) == 0 {
		refreshToken = tokenResp.Data.RefreshToken
	}
	expiresIn := tokenResp.ExpiresIn
	if expiresIn == 0 {
		expiresIn = tokenResp.Data.ExpiresIn
	}
	token = &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}
	if expiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
	mid := tokenResp.Mid
	if mid == 0 {
		mid = tokenResp.Data.Mid
	}
	if mid > 0 {
		token = token.WithExtra(map[string]interface{}{"mid": mid})
	}
	return
}

func (pro *BilibiliProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://member.bilibili.com/arcopen/fn/user/account/info"
	}
	requestURL := userInfoEndpoint
	if !strings.Contains(userInfoEndpoint, "access_token=") {
		sep := "?"
		if strings.Contains(userInfoEndpoint, "?") {
			sep = "&"
		}
		requestURL = userInfoEndpoint + sep + "access_token=" + url.QueryEscape(token.AccessToken)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		errorData.Err = err
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
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
		errorData.Err = fmt.Errorf("request bilibili userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userResp bilibiliUserResponse
	if err = jsoniter.Unmarshal(data, &userResp); err != nil {
		errorData.Err = err
		return
	}
	if userResp.Code != 0 {
		errorData.Err = fmt.Errorf("request bilibili userinfo failed, code: %d, message: %s %s", userResp.Code, userResp.Message, userResp.Msg)
		return
	}
	user := userResp.Data
	if user.Mid == 0 {
		user.Mid = userResp.Mid
	}
	if len(user.Uname) == 0 {
		user.Uname = userResp.Uname
	}
	if len(user.Name) == 0 {
		user.Name = userResp.Name
	}
	if len(user.Face) == 0 {
		user.Face = userResp.Face
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcBilibili
	if user.Mid > 0 {
		profile.LoginID = strconv.FormatInt(user.Mid, 10)
	}
	profile.Username = user.Uname
	if len(profile.Username) == 0 {
		profile.Username = user.Name
	}
	profile.Nickname = user.Name
	if len(profile.Nickname) == 0 {
		profile.Nickname = user.Nickname
	}
	profile.Avatar = user.Face
	if len(profile.Avatar) == 0 {
		profile.Avatar = user.Avatar
	}
	profile.Email = user.Email
	if len(user.SpaceURL) > 0 {
		profile.Home = user.SpaceURL
	} else if len(profile.LoginID) > 0 {
		profile.Home = "https://space.bilibili.com/" + profile.LoginID
	}
	return
}
