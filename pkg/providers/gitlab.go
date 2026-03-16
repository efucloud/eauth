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

type GitlabProvider struct {
	ctx      context.Context
	provider *dtos.ProviderOidcDetail
	oauthCfg *oauth2.Config
	//redirectUrl string
	scopes []string
}

type gitlabUserInfo struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
	Username        string      `json:"username"`
	State           string      `json:"state"`
	AvatarUrl       string      `json:"avatar_url"`
	WebUrl          string      `json:"web_url"`
	CreatedAt       time.Time   `json:"created_at"`
	Bio             string      `json:"bio"`
	BioHtml         string      `json:"bio_html"`
	Location        string      `json:"location"`
	PublicEmail     string      `json:"public_email"`
	Skype           string      `json:"skype"`
	Linkedin        string      `json:"linkedin"`
	Twitter         string      `json:"twitter"`
	WebsiteUrl      string      `json:"website_url"`
	Organization    string      `json:"organization"`
	JobTitle        string      `json:"job_title"`
	Pronouns        interface{} `json:"pronouns"`
	Bot             bool        `json:"bot"`
	WorkInformation interface{} `json:"work_information"`
	Followers       int         `json:"followers"`
	Following       int         `json:"following"`
	LastSignInAt    time.Time   `json:"last_sign_in_at"`
	ConfirmedAt     time.Time   `json:"confirmed_at"`
	LastActivityOn  string      `json:"last_activity_on"`
	Email           string      `json:"email"`
	ThemeId         int         `json:"theme_id"`
	ColorSchemeId   int         `json:"color_scheme_id"`
	ProjectsLimit   int         `json:"projects_limit"`
	CurrentSignInAt time.Time   `json:"current_sign_in_at"`
	Identities      []struct {
		Provider       string      `json:"provider"`
		ExternUid      string      `json:"extern_uid"`
		SamlProviderId interface{} `json:"saml_provider_id"`
	} `json:"identities"`
	CanCreateGroup                 bool        `json:"can_create_group"`
	CanCreateProject               bool        `json:"can_create_project"`
	TwoFactorEnabled               bool        `json:"two_factor_enabled"`
	External                       bool        `json:"external"`
	PrivateProfile                 bool        `json:"private_profile"`
	CommitEmail                    string      `json:"commit_email"`
	SharedRunnersMinutesLimit      interface{} `json:"shared_runners_minutes_limit"`
	ExtraSharedRunnersMinutesLimit interface{} `json:"extra_shared_runners_minutes_limit"`
}

var _ IdProvider = &GitlabProvider{}

func NewGitlabProvider(ctx context.Context, redirectUrl string, provider dtos.ProviderOidcDetail) *GitlabProvider {
	pro := &GitlabProvider{
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
func (pro *GitlabProvider) GetToken(ctx context.Context, code string) (token *oauth2.Token, errorData common.ErrorData) {
	token, errorData.Err = pro.oauthCfg.Exchange(ctx, code)
	return
}
func (pro *GitlabProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (authProfile dtos.ThirdAuthProfile, errorData common.ErrorData) {
	authProfile, errorData.Err = pro.getUserinfoFormAPI(ctx, token)
	return
}
func (pro *GitlabProvider) getUserinfoFormAPI(ctx context.Context, token *oauth2.Token) (profile dtos.ThirdAuthProfile, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	url := strings.TrimSuffix(pro.provider.Issuer, "/") + "/api/v4/user?access_token=" + token.AccessToken
	if len(pro.provider.UserinfoEndpoint) > 0 {
		url = pro.provider.UserinfoEndpoint + token.AccessToken
	}
	resp, err := client.Get(url)
	if err != nil {
		return profile, err
	}
	defer resp.Body.Close()
	var userInfo gitlabUserInfo
	data, err := io.ReadAll(resp.Body)
	err = jsoniter.Unmarshal(data, &userInfo)
	if err != nil {
		return profile, err
	}

	profile.Username = userInfo.Username
	profile.Nickname = userInfo.Name
	profile.Avatar = userInfo.AvatarUrl
	profile.Email = userInfo.Email
	profile.Home = userInfo.WebUrl
	profile.Provider = config.ProviderOidcGitlab
	profile.LoginID = fmt.Sprintf("%d", userInfo.Id)
	_ = jsoniter.Unmarshal(data, &profile.Properties)
	return profile, err
}
