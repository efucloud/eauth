package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

type WechatProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type wechatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	UnionID      string `json:"unionid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

type wechatUserResponse struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
}

var _ IdProvider = &WechatProvider{}

func NewWechatProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *WechatProvider {
	pro := &WechatProvider{
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

func (pro *WechatProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	params := url.Values{}
	params.Set("appid", pro.provider.ClientId)
	params.Set("secret", pro.provider.ClientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	requestURL := pro.provider.TokenEndpoint + "?" + params.Encode()
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
		errorData.Err = fmt.Errorf("request wechat token failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var tokenResp wechatTokenResponse
	if err = jsoniter.Unmarshal(data, &tokenResp); err != nil {
		errorData.Err = err
		return
	}
	if tokenResp.ErrCode != 0 {
		errorData.Err = fmt.Errorf("request wechat token failed, errcode: %d, errmsg: %s", tokenResp.ErrCode, tokenResp.ErrMsg)
		return
	}
	if len(tokenResp.AccessToken) == 0 {
		errorData.Err = fmt.Errorf("request wechat token failed, empty access_token")
		return
	}
	token = &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    "Bearer",
	}
	if tokenResp.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	token = token.WithExtra(map[string]interface{}{
		"openid":  tokenResp.OpenID,
		"unionid": tokenResp.UnionID,
	})
	return
}

func (pro *WechatProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	openID := fmt.Sprintf("%v", token.Extra("openid"))
	if openID == "<nil>" || len(openID) == 0 {
		errorData.Err = fmt.Errorf("wechat token missing openid")
		return
	}
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://api.weixin.qq.com/sns/userinfo"
	}
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")
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
		errorData.Err = fmt.Errorf("request wechat userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userResp wechatUserResponse
	if err = jsoniter.Unmarshal(data, &userResp); err != nil {
		errorData.Err = err
		return
	}
	if userResp.ErrCode != 0 {
		errorData.Err = fmt.Errorf("request wechat userinfo failed, errcode: %d, errmsg: %s", userResp.ErrCode, userResp.ErrMsg)
		return
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcWechat
	profile.LoginID = userResp.UnionID
	if len(profile.LoginID) == 0 {
		profile.LoginID = userResp.OpenID
	}
	profile.Username = userResp.Nickname
	profile.Nickname = userResp.Nickname
	profile.Avatar = userResp.Headimgurl
	profile.Home = "https://weixin.qq.com/"
	return
}
