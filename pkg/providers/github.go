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
	"time"

	"github.com/efucloud/eauth/pkg/config"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	scopes   []string
	ctx      context.Context
}

func (pro *GithubProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (authProfile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	if len(pro.provider.UserinfoEndpoint) > 0 {
		client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
		userInfoResp, err := client.Get(pro.provider.UserinfoEndpoint)
		if err != nil {
			errorData.Err = err
			return
		}
		defer userInfoResp.Body.Close()
		data, err := io.ReadAll(userInfoResp.Body)
		if err != nil {
			errorData.Err = err
			return
		}
		var user githubUserInfo
		err = jsoniter.Unmarshal(data, &user)
		if err != nil {
			config.Logger.Infof("thirty auth provider github response data: %s", string(data))
			errorData.Err = err
			return
		}
		_ = jsoniter.Unmarshal(data, &authProfile.Properties)
		authProfile.Provider = config.ProviderOidcGithub
		authProfile.LoginID = fmt.Sprintf("%d", user.Id)
		authProfile.Username = user.Login
		authProfile.Avatar = user.AvatarUrl
		authProfile.Email = user.Email
		authProfile.Home = user.HtmlUrl

	} else {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 30 * time.Second,
		}
		url := "https://github.com/api/v4/user?access_token=" + token.AccessToken
		resp, err := client.Get(url)
		if err != nil {
			errorData.Err = err
			return
		}
		defer resp.Body.Close()
		var userInfo githubUserInfo
		data, _ := io.ReadAll(resp.Body)
		err = jsoniter.Unmarshal(data, &userInfo)
		if err != nil {
			errorData.Err = err
			return
		}
		_ = jsoniter.Unmarshal(data, &authProfile.Properties)
		authProfile.Provider = config.ProviderOidcGithub
		authProfile.LoginID = fmt.Sprintf("%d", userInfo.Id)
		authProfile.Username = userInfo.Login
		authProfile.Nickname = userInfo.Name
		authProfile.Home = userInfo.HtmlUrl

	}
	return
}

type githubUserInfo struct {
	Login                   string      `json:"login"`
	Id                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type "`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                string      `json:"location"`
	Email                   string      `json:"email"`
	Hireable                bool        `json:"hireable"`
	Bio                     string      `json:"bio"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		Collaborators int    `json:"collaborators"`
		PrivateRepos  int    `json:"private_repos"`
	} `json:"plan"`
}

var _ IdProvider = &GithubProvider{}

func NewGithubProvider(ctx context.Context, redirectUrl string, provider dtos.ProviderOidcDetail) *GithubProvider {
	pro := &GithubProvider{
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
func (pro *GithubProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}
