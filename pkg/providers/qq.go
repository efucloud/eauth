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

type QQProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type qqOpenIDResponse struct {
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid"`
}

type qqUserInfo struct {
	Ret          int    `json:"ret"`
	Msg          string `json:"msg"`
	Nickname     string `json:"nickname"`
	Figureurl    string `json:"figureurl"`
	Figureurl1   string `json:"figureurl_1"`
	Figureurl2   string `json:"figureurl_2"`
	FigureurlQQ1 string `json:"figureurl_qq_1"`
	FigureurlQQ2 string `json:"figureurl_qq_2"`
	Gender       string `json:"gender"`
}

var _ IdProvider = &QQProvider{}

func NewQQProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *QQProvider {
	pro := &QQProvider{
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

func (pro *QQProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", pro.provider.ClientId)
	formData.Set("client_secret", pro.provider.ClientSecret)
	formData.Set("code", code)
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
	body := strings.TrimSpace(string(data))
	if resp.StatusCode >= http.StatusBadRequest {
		errorData.Err = fmt.Errorf("request qq token failed, status: %d, body: %s", resp.StatusCode, body)
		return
	}
	if strings.Contains(body, "callback(") && strings.Contains(body, "error") {
		errorData.Err = fmt.Errorf("request qq token failed, response: %s", body)
		return
	}
	queryData, err := url.ParseQuery(body)
	if err != nil {
		errorData.Err = err
		return
	}
	accessToken := queryData.Get("access_token")
	if len(accessToken) == 0 {
		errorData.Err = fmt.Errorf("qq token response does not contain access_token: %s", body)
		return
	}
	token = &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: queryData.Get("refresh_token"),
	}
	expiresIn := queryData.Get("expires_in")
	if len(expiresIn) > 0 {
		expires, parseErr := strconv.ParseInt(expiresIn, 10, 64)
		if parseErr == nil && expires > 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	}
	return
}

func (pro *QQProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	var openID string
	openID, errorData = pro.getOpenIDByAccessToken(token.AccessToken)
	if errorData.IsNotNil() {
		return
	}
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://graph.qq.com/user/get_user_info"
	}
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	params.Set("oauth_consumer_key", pro.provider.ClientId)
	params.Set("openid", openID)
	requestURL := userInfoEndpoint + "?" + params.Encode()
	resp, err := http.Get(requestURL)
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
		errorData.Err = fmt.Errorf("request qq userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfo qqUserInfo
	if err = jsoniter.Unmarshal(data, &userInfo); err != nil {
		errorData.Err = err
		return
	}
	if userInfo.Ret != 0 {
		errorData.Err = fmt.Errorf("request qq userinfo failed, ret: %d, msg: %s", userInfo.Ret, userInfo.Msg)
		return
	}
	var properties dtos.JsonMap
	_ = jsoniter.Unmarshal(data, &properties)
	if properties == nil {
		properties = dtos.JsonMap{}
	}
	properties["openid"] = openID

	profile.Provider = config.ProviderOidcQq
	profile.LoginID = openID
	profile.Username = userInfo.Nickname
	profile.Nickname = userInfo.Nickname
	profile.Avatar = userInfo.FigureurlQQ2
	if len(profile.Avatar) == 0 {
		profile.Avatar = userInfo.Figureurl2
	}
	if len(profile.Avatar) == 0 {
		profile.Avatar = userInfo.FigureurlQQ1
	}
	profile.Home = "https://qzone.qq.com/"
	profile.Properties = properties
	return
}

func (pro *QQProvider) getOpenIDByAccessToken(accessToken string) (openID string, errorData common.ErrorData) {
	requestURL := "https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(accessToken)
	resp, err := http.Get(requestURL)
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
	body := strings.TrimSpace(string(data))
	if resp.StatusCode >= http.StatusBadRequest {
		errorData.Err = fmt.Errorf("request qq openid failed, status: %d, body: %s", resp.StatusCode, body)
		return
	}
	if strings.Contains(body, "callback(") {
		start := strings.Index(body, "{")
		end := strings.LastIndex(body, "}")
		if start >= 0 && end > start {
			body = body[start : end+1]
		}
	}
	var openIDResp qqOpenIDResponse
	if err = jsoniter.Unmarshal([]byte(body), &openIDResp); err != nil {
		errorData.Err = err
		return
	}
	openID = openIDResp.OpenID
	if len(openID) == 0 {
		errorData.Err = fmt.Errorf("qq openid response does not contain openid")
	}
	return
}
