package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type TiktokProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type tiktokTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"open_id"`
	UnionID      string `json:"union_id"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type tiktokTokenResponse struct {
	ErrorCode    int             `json:"error_code"`
	Description  string          `json:"description"`
	Message      string          `json:"message"`
	Code         int             `json:"code"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	OpenID       string          `json:"open_id"`
	UnionID      string          `json:"union_id"`
	ExpiresIn    int64           `json:"expires_in"`
	Data         tiktokTokenData `json:"data"`
}

type tiktokUserData struct {
	OpenID    string `json:"open_id"`
	UnionID   string `json:"union_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AvatarURL string `json:"avatar_url"`
}

type tiktokUserResponse struct {
	ErrorCode   int            `json:"error_code"`
	Description string         `json:"description"`
	Message     string         `json:"message"`
	Code        int            `json:"code"`
	Data        tiktokUserData `json:"data"`
	OpenID      string         `json:"open_id"`
	UnionID     string         `json:"union_id"`
	Nickname    string         `json:"nickname"`
	Avatar      string         `json:"avatar"`
}

var _ IdProvider = &TiktokProvider{}

func NewTiktokProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *TiktokProvider {
	pro := &TiktokProvider{
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

func (pro *TiktokProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	formData := url.Values{}
	formData.Set("client_key", pro.provider.ClientId)
	formData.Set("client_secret", pro.provider.ClientSecret)
	formData.Set("client_id", pro.provider.ClientId)
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")
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
		errorData.Err = fmt.Errorf("request tiktok token failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var tokenResp tiktokTokenResponse
	if err = jsoniter.Unmarshal(data, &tokenResp); err != nil {
		errorData.Err = err
		return
	}
	if tokenResp.ErrorCode != 0 || tokenResp.Code != 0 {
		errorData.Err = fmt.Errorf("request tiktok token failed, code: %d/%d, message: %s %s", tokenResp.ErrorCode, tokenResp.Code, tokenResp.Description, tokenResp.Message)
		return
	}
	accessToken := tokenResp.AccessToken
	if len(accessToken) == 0 {
		accessToken = tokenResp.Data.AccessToken
	}
	if len(accessToken) == 0 {
		errorData.Err = fmt.Errorf("request tiktok token failed, empty access_token")
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
	openID := tokenResp.OpenID
	if len(openID) == 0 {
		openID = tokenResp.Data.OpenID
	}
	unionID := tokenResp.UnionID
	if len(unionID) == 0 {
		unionID = tokenResp.Data.UnionID
	}
	token = &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}
	if expiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
	token = token.WithExtra(map[string]interface{}{
		"open_id":  openID,
		"union_id": unionID,
	})
	return
}

func (pro *TiktokProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	openID := fmt.Sprintf("%v", token.Extra("open_id"))
	if openID == "<nil>" {
		openID = ""
	}
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://open.douyin.com/oauth/userinfo/"
	}
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	if len(openID) > 0 {
		params.Set("open_id", openID)
	}
	requestURL := userInfoEndpoint
	if strings.Contains(userInfoEndpoint, "?") {
		requestURL = userInfoEndpoint + "&" + params.Encode()
	} else {
		requestURL = userInfoEndpoint + "?" + params.Encode()
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
		errorData.Err = fmt.Errorf("request tiktok userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userResp tiktokUserResponse
	if err = jsoniter.Unmarshal(data, &userResp); err != nil {
		errorData.Err = err
		return
	}
	if userResp.ErrorCode != 0 || userResp.Code != 0 {
		errorData.Err = fmt.Errorf("request tiktok userinfo failed, code: %d/%d, message: %s %s", userResp.ErrorCode, userResp.Code, userResp.Description, userResp.Message)
		return
	}
	user := userResp.Data
	if len(user.OpenID) == 0 {
		user.OpenID = userResp.OpenID
	}
	if len(user.UnionID) == 0 {
		user.UnionID = userResp.UnionID
	}
	if len(user.Nickname) == 0 {
		user.Nickname = userResp.Nickname
	}
	if len(user.Avatar) == 0 {
		user.Avatar = userResp.Avatar
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcTiktok
	profile.LoginID = user.UnionID
	if len(profile.LoginID) == 0 {
		profile.LoginID = user.OpenID
	}
	profile.Username = user.Nickname
	profile.Nickname = user.Nickname
	profile.Avatar = user.Avatar
	if len(profile.Avatar) == 0 {
		profile.Avatar = user.AvatarURL
	}
	profile.Home = "https://www.douyin.com/user/" + profile.LoginID
	return
}
