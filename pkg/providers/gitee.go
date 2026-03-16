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

type GiteeProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	//redirectUrl string
	scopes []string
}

type giteeUserInfo struct {
	AvatarUrl         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsUrl         string `json:"events_url"`
	Followers         int    `json:"followers"`
	FollowersUrl      string `json:"followers_url"`
	Following         int    `json:"following"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	HtmlUrl           string `json:"html_url"`
	Id                int    `json:"id"`
	Login             string `json:"login"`
	MemberRole        string `json:"member_role"`
	Name              string `json:"name"`
	OrganizationsUrl  string `json:"organizations_url"`
	PublicGists       int    `json:"public_gists"`
	PublicRepos       int    `json:"public_repos"`
	ReceivedEventsUrl string `json:"received_events_url"`
	ReposUrl          string `json:"repos_url"`
	Stared            int    `json:"stared"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	Type              string `json:"type "`
	UpdatedAt         string `json:"updated_at"`
	Url               string `json:"url"`
	Watched           int    `json:"watched"`
	Weibo             string `json:"weibo"`
}

var _ IdProvider = &GiteeProvider{}

func NewGiteeProvider(ctx context.Context, redirectUrl string, provider dtos.ProviderOidcDetail) *GiteeProvider {
	pro := &GiteeProvider{
		ctx:      ctx,
		provider: &provider,
		oauthCfg: nil,
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
func (pro *GiteeProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}
func (pro *GiteeProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (authProfile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	authProfile, errorData.Err = pro.getUserinfoFormAPI(ctx, token)
	return
}
func (pro *GiteeProvider) getUserinfoFormAPI(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	//url := strings.TrimSuffix(provider.IssuerURL, "/") + "/api/v4/user?access_token=" + token.AccessToken
	url := strings.TrimSuffix(pro.provider.Issuer, "/") + "/api/v5/user?access_token=" + token.AccessToken
	if len(pro.provider.UserinfoEndpoint) > 0 {
		url = pro.provider.UserinfoEndpoint + token.AccessToken
	}
	resp, err := client.Get(url)
	if err != nil {
		return profile, err
	}
	defer resp.Body.Close()
	var userInfo giteeUserInfo
	data, err := io.ReadAll(resp.Body)
	//config.Logger.Infof("response data: %s", string(data))
	err = jsoniter.Unmarshal(data, &userInfo)
	if err != nil {
		return profile, err
	}
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	profile.Username = userInfo.Login
	profile.Nickname = userInfo.Name
	profile.Email = userInfo.Email
	profile.Avatar = userInfo.AvatarUrl
	profile.Home = userInfo.HtmlUrl
	profile.Provider = config.ProviderOidcGitee
	profile.LoginID = fmt.Sprintf("%d", userInfo.Id)
	return profile, err
}
