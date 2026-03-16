package providers

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/models/dtos"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/efucloud/eauth/pkg/config"
	"golang.org/x/oauth2"
)

// https://open.dingtalk.com/document/isvapp/obtain-identity-credentials
type dingTalkAccessToken struct {
	ErrCode     int    `json:"code"`
	ErrMsg      string `json:"message"`
	AccessToken string `json:"accessToken"` //Interface call credentials
	ExpiresIn   int64  `json:"expireIn"`    //access_token interface call credential timeout time, unit (seconds)
}
type dingTalkUserResponse struct {
	Nick      string `json:"nick"`
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	AvatarUrl string `json:"avatarUrl"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	StateCode string `json:"stateCode"`
}
type DingTalkProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	//redirectUrl string
	scopes []string
}

var _ IdProvider = &DingTalkProvider{}

func NewDingTalkProvider(ctx context.Context, redirectUrl string, provider dtos.ProviderOidcDetail) *DingTalkProvider {
	pro := &DingTalkProvider{
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
		RedirectURL: redirectUrl,
		Scopes:      provider.Scopes,
	}
	return pro
}
func (pro *DingTalkProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	config.Logger.Infof("oauth cfg %+v", pro.oauthCfg)
	pTokenParams := &struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		Code         string `json:"code"`
		GrantType    string `json:"grantType"`
	}{pro.provider.ClientId, pro.provider.ClientSecret, code, "authorization_code"}

	data, err := pro.postWithBody(pTokenParams, pro.provider.TokenEndpoint)
	if err != nil {
		errorData.Err = err
		return token, errorData
	}

	pToken := &dingTalkAccessToken{}
	err = jsoniter.Unmarshal(data, pToken)
	if err != nil {
		errorData.Err = err
		return token, errorData
	}
	if pToken.ErrCode != 0 {
		errorData.Err = fmt.Errorf("pToken.Errcode = %d, pToken.Errmsg = %s", pToken.ErrCode, pToken.ErrMsg)
		return token, errorData
	}
	token = &oauth2.Token{
		AccessToken: pToken.AccessToken,
		Expiry:      time.Unix(time.Now().Unix()+pToken.ExpiresIn, 0),
	}
	return token, errorData
}
func (pro *DingTalkProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	profile, errorData.Err = pro.getUserinfoFormAPI(ctx, token)
	return
}
func (pro *DingTalkProvider) getUserinfoFormAPI(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, err error) {
	dtUserInfo := &dingTalkUserResponse{}
	accessToken := token.AccessToken
	userInfoEndpoint := pro.provider.UserinfoEndpoint
	if len(userInfoEndpoint) == 0 {
		userInfoEndpoint = pro.provider.AuthorizationEndpoint
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return profile, err
	}
	req.Header.Add("x-acs-dingtalk-access-token", accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return profile, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return profile, err
	}

	err = jsoniter.Unmarshal(data, dtUserInfo)
	if err != nil {
		return profile, err
	}
	profile.Username = dtUserInfo.Nick
	profile.Nickname = dtUserInfo.Nick
	profile.Phone = dtUserInfo.Mobile
	profile.Email = dtUserInfo.Email
	profile.Avatar = dtUserInfo.AvatarUrl
	profile.Provider = config.ProviderOidcDingTalk
	profile.LoginID = dtUserInfo.UnionId
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	return profile, err
}

func (pro *DingTalkProvider) postWithBody(body interface{}, url string) ([]byte, error) {
	bs, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(string(bs))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(url, "application/json;charset=UTF-8", r)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	return data, nil
}
