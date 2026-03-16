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

type WeiboProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
}

type weiboTokenInfo struct {
	UID string `json:"uid"`
}

type weiboUserInfo struct {
	ID              int64  `json:"id"`
	IDStr           string `json:"idstr"`
	ScreenName      string `json:"screen_name"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
	AvatarLarge     string `json:"avatar_large"`
	ProfileURL      string `json:"profile_url"`
	URL             string `json:"url"`
}

var _ IdProvider = &WeiboProvider{}

func NewWeiboProvider(ctx context.Context, redirectURL string, provider dtos.ProviderOidcDetail) *WeiboProvider {
	pro := &WeiboProvider{
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

func (pro *WeiboProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}

func (pro *WeiboProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	uid := fmt.Sprintf("%v", token.Extra("uid"))
	if uid == "<nil>" || len(uid) == 0 {
		uid, errorData.Err = pro.getUIDByToken(token.AccessToken)
		if errorData.IsNotNil() {
			return
		}
	}
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = "https://api.weibo.com/2/users/show.json"
	}
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	params.Set("uid", uid)
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
		errorData.Err = fmt.Errorf("request weibo userinfo failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var userInfo weiboUserInfo
	if err = jsoniter.Unmarshal(data, &userInfo); err != nil {
		errorData.Err = err
		return
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Provider = config.ProviderOidcWeibo
	profile.LoginID = userInfo.IDStr
	if len(profile.LoginID) == 0 {
		profile.LoginID = fmt.Sprintf("%d", userInfo.ID)
	}
	profile.Username = userInfo.ScreenName
	profile.Nickname = userInfo.Name
	profile.Avatar = userInfo.AvatarLarge
	if len(profile.Avatar) == 0 {
		profile.Avatar = userInfo.ProfileImageURL
	}
	if len(userInfo.URL) > 0 {
		profile.Home = userInfo.URL
	} else if len(userInfo.ProfileURL) > 0 {
		profile.Home = "https://weibo.com/" + userInfo.ProfileURL
	}
	return
}

func (pro *WeiboProvider) getUIDByToken(accessToken string) (uid string, err error) {
	formData := url.Values{}
	formData.Set("access_token", accessToken)
	resp, err := http.PostForm("https://api.weibo.com/oauth2/get_token_info", formData)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("request weibo token info failed, status: %d, body: %s", resp.StatusCode, string(data))
		return
	}
	var tokenInfo weiboTokenInfo
	if err = jsoniter.Unmarshal(data, &tokenInfo); err != nil {
		return
	}
	uid = tokenInfo.UID
	if len(uid) == 0 {
		err = fmt.Errorf("can't find uid from weibo token response")
	}
	return
}
