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

type WechatWorkProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type wechatWorkTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type wechatWorkCodeUserResponse struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	UserID         string `json:"UserId"`
	OpenID         string `json:"OpenId"`
	ExternalUserID string `json:"external_userid"`
}

type wechatWorkUserInfoResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	UserID  string `json:"userid"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
}

var _ IdProvider = &WechatWorkProvider{}

func NewWechatWorkProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *WechatWorkProvider {
	pro := &WechatWorkProvider{
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

func (pro *WechatWorkProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	accessToken, expiresIn, errorData := pro.getCorpAccessToken()
	if errorData.IsNotNil() {
		return
	}
	userByCode, err := pro.getUserByCode(accessToken, code)
	if err != nil {
		errorData.Err = err
		return
	}
	token = &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}
	if expiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
	token = token.WithExtra(map[string]interface{}{
		"user_id":          userByCode.UserID,
		"open_id":          userByCode.OpenID,
		"external_user_id": userByCode.ExternalUserID,
	})
	return
}

func (pro *WechatWorkProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	userID := fmt.Sprintf("%v", token.Extra("user_id"))
	if userID == "<nil>" {
		userID = ""
	}
	openID := fmt.Sprintf("%v", token.Extra("open_id"))
	if openID == "<nil>" {
		openID = ""
	}
	if len(userID) == 0 {
		profile.Provider = config.ProviderOidcWechatWork
		profile.LoginID = openID
		profile.Username = openID
		profile.Nickname = openID
		profile.Home = "https://work.weixin.qq.com/"
		return
	}

	userInfo, err := pro.getUserInfoByUserID(token.AccessToken, userID)
	if err != nil {
		errorData.Err = err
		return
	}
	profile.Provider = config.ProviderOidcWechatWork
	profile.LoginID = userInfo.UserID
	profile.Username = userInfo.UserID
	profile.Nickname = userInfo.Name
	if len(profile.Username) == 0 {
		profile.Username = userInfo.Name
	}
	profile.Avatar = userInfo.Avatar
	profile.Email = userInfo.Email
	profile.Phone = userInfo.Mobile
	profile.Home = "https://work.weixin.qq.com/"

	properties := dtos.JsonMap{}
	properties["user_id"] = userInfo.UserID
	properties["name"] = userInfo.Name
	properties["avatar"] = userInfo.Avatar
	properties["email"] = userInfo.Email
	properties["mobile"] = userInfo.Mobile
	properties["open_id"] = openID
	profile.Properties = properties
	return
}

func (pro *WechatWorkProvider) getCorpAccessToken() (accessToken string, expiresIn int64, errorData common.ErrorData) {
	params := url.Values{}
	params.Set("corpid", pro.provider.ClientId)
	params.Set("corpsecret", pro.provider.ClientSecret)
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
		errorData.Err = fmt.Errorf("request wechat work token failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var tokenResp wechatWorkTokenResponse
	if err = jsoniter.Unmarshal(data, &tokenResp); err != nil {
		errorData.Err = err
		return
	}
	if tokenResp.ErrCode != 0 {
		errorData.Err = fmt.Errorf("request wechat work token failed, errcode: %d, errmsg: %s", tokenResp.ErrCode, tokenResp.ErrMsg)
		return
	}
	accessToken = tokenResp.AccessToken
	expiresIn = tokenResp.ExpiresIn
	if len(accessToken) == 0 {
		errorData.Err = fmt.Errorf("request wechat work token failed, empty access_token")
	}
	return
}

func (pro *WechatWorkProvider) getUserByCode(corpAccessToken, code string) (result wechatWorkCodeUserResponse, err error) {
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
	}
	params := url.Values{}
	params.Set("access_token", corpAccessToken)
	params.Set("code", code)
	requestURL := userInfoEndpoint + "?" + params.Encode()
	resp, err := http.Get(requestURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("request wechat work user by code failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	if err = jsoniter.Unmarshal(data, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("request wechat work user by code failed, errcode: %d, errmsg: %s", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (pro *WechatWorkProvider) getUserInfoByUserID(corpAccessToken, userID string) (result wechatWorkUserInfoResponse, err error) {
	params := url.Values{}
	params.Set("access_token", corpAccessToken)
	params.Set("userid", userID)
	requestURL := "https://qyapi.weixin.qq.com/cgi-bin/user/get?" + params.Encode()
	resp, err := http.Get(requestURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("request wechat work user info failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	if err = jsoniter.Unmarshal(data, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("request wechat work user info failed, errcode: %d, errmsg: %s", result.ErrCode, result.ErrMsg)
	}
	return
}
