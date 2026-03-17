package services

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"github.com/efucloud/eauth/pkg/providers"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
	"github.com/pquerna/otp/totp"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
	"time"
)

type OAuthService struct {
}

func (svc *OAuthService) ApplicationAuthCode(ctx context.Context, userId uint, model dtos.OidcCodeRequest) (result dtos.OidcCodeResponse, errorData common.ErrorData) {
	var (
		user                          dtos.UserDetail
		app                           dtos.ApplicationDetail
		codeChallenge                 string
		codeChallengeMethodNormalized string
	)
	if model.ResponseType != "code" {
		errorData.Err = fmt.Errorf("response_type  must be code ")
		return
	}
	codeChallenge, codeChallengeMethodNormalized, errorData.Err = normalizePKCEParams(model.GetCodeChallenge(), model.GetCodeChallengeMethod())
	if errorData.IsNotNil() {
		return
	}

	userSvc := UserService{}
	user, errorData = userSvc.GetUserByID(ctx, userId)
	if errorData.IsNotNil() {
		return
	}
	appSvc := ApplicationService{}
	app, errorData = appSvc.GetApplicationByClientId(ctx, model.ClientId)
	if errorData.IsNotNil() {
		return
	}
	if !app.RedirectUriMatch(model.RedirectUri) {
		errorData.Err = fmt.Errorf("redirect_uri: %s is not right", model.RedirectUri)
		return
	}

	result.Code = common.NewSecureID(16)
	result.State = model.State
	result.RedirectUri = fmt.Sprintf("%s?code=%s&state=%s", model.RedirectUri, result.Code, model.State)
	authApp := dtos.AppAuthRecordCreate{
		ApplicationId:       app.ID,
		Code:                result.Code,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethodNormalized,
		UserId:              user.ID,
	}
	appAuthSvc := AppAuthRecordService{}
	_, errorData = appAuthSvc.AddAppAuthRecord(ctx, authApp)

	return
}
func (svc *OAuthService) UpdateUserAvatar(ctx context.Context, userId uint, avatarAddress string) (errorData common.ErrorData) {
	userSvc := UserService{}
	return userSvc.UpdateUserAvatar(ctx, userId, avatarAddress)
}
func (svc *OAuthService) ChangeSelfInformation(ctx context.Context, model dtos.UserUpdate) (userinfo dtos.AuthedUserInfo, errorData common.ErrorData) {
	var (
		user dtos.UserDetail
	)
	model.Default(ctx)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		errorData.MsgCode = config.MsgCodeRequestDataInvalid
		config.Logger.Errorf("user: %d update failed, err: %s", model.ID, errorData.Err.Error())
		return
	}
	userSvc := UserService{}
	user, errorData = userSvc.UpdateUser(ctx, model)
	if errorData.IsNotNil() {
		return
	}
	copyByJSON(user, &userinfo)
	userinfo.Role = user.Role
	userinfo.HasPassword = len(user.PasswordStore) > 0
	return userinfo, errorData
}
func (svc *OAuthService) Userinfo(ctx context.Context, userId uint) (userinfo dtos.AuthedUserInfo, errorData common.ErrorData) {

	var (
		user dtos.UserDetail
	)
	userSvc := UserService{}
	user, errorData = userSvc.GetUserByID(ctx, userId)
	if errorData.IsNotNil() {
		config.Logger.Errorf("can't get user by id: %d, err: %s", userId, errorData.Err.Error())
		return userinfo, errorData
	}
	if !user.Enable {
		config.Logger.Errorf("can't get user by id: %d, user id disabled", userId)
		return userinfo, errorData
	}
	copyByJSON(user, &userinfo)
	userinfo.Role = user.Role
	userinfo.HasPassword = len(user.PasswordStore) > 0
	userinfo.ID = user.ID
	userinfo.Avatar = fmt.Sprintf("%s%s", config.ApplicationConfig.ServerAddress, userinfo.Avatar)
	return
}

// GetTokenByAuthorizationCode 授权码模式获取token
func (svc *OAuthService) GetTokenByAuthorizationCode(ctx context.Context, req *restful.Request, resp *restful.Response, lang string) {
	var (
		model      dtos.OidcRequestToken
		errorData  common.ErrorData
		authRecord dtos.AppAuthRecordDetail
		user       dtos.UserDetail
		token      dtos.AccessTokenResponse
		app        dtos.ApplicationDetail
	)

	contentType := req.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, config.RequestForm) || req.Request.Method == "GET" || len(contentType) == 0 {
		decoder := schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)
		errorData.Err = req.Request.ParseForm()
		if errorData.IsNil() {
			formValues := req.Request.PostForm
			if len(formValues) == 0 {
				formValues = req.Request.Form
			}
			errorData.Err = decoder.Decode(&model, formValues)
		}
		if username, password, ok := req.Request.BasicAuth(); ok {
			model.ClientId = username
			model.ClientSecret = password
		}
	} else {
		errorData.Err = json.NewDecoder(req.Request.Body).Decode(&model)
	}
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	authRecordSvc := AppAuthRecordService{}
	authRecord, errorData = authRecordSvc.GetAppAuthRecordByCode(ctx, model.Code)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	userSvc := UserService{}
	user, errorData = userSvc.GetUserByID(ctx, authRecord.UserId)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	appSvc := ApplicationService{}
	app, errorData = appSvc.GetApplicationById(ctx, authRecord.ApplicationId)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	if len(model.ClientId) > 0 && model.ClientId != app.ClientId {
		errorData.Err = fmt.Errorf("client id: %s is not right", model.ClientId)
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	errorData.Err = verifyPKCE(authRecord.CodeChallenge, authRecord.CodeChallengeMethod, model.GetCodeVerifier())
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	if len(authRecord.CodeChallenge) == 0 && len(model.ClientSecret) == 0 {
		errorData.Err = fmt.Errorf("authorization code has no code_challenge, public client must send code_challenge in /oidc/code request")
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	isPKCEPublicClient := len(authRecord.CodeChallenge) > 0 && len(model.ClientSecret) == 0
	if isPKCEPublicClient {
		if len(model.ClientId) == 0 {
			errorData.Err = fmt.Errorf("client_id is required for pkce public client")
			config.Logger.Error(errorData.Err)
			errorData.Lang = lang
			common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
			return
		}
	} else if app.ClientSecret != model.ClientSecret {
		errorData.Err = fmt.Errorf("client secret: %s is not right", model.ClientSecret)
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}

	token, errorData = svc.GenerateTokenResponse(ctx, true, app.ClientId, user)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.Lang = lang
		common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
		return
	}
	common.ResponseSuccess(resp, token)

}
func (svc *OAuthService) GetTokenByPassword(ctx context.Context, req *restful.Request, resp *restful.Response, lang string) {

}

func (svc *OAuthService) GetTokenByToken(ctx context.Context, req *restful.Request, resp *restful.Response, lang string) {

}

func (svc *OAuthService) RefreshToken(ctx context.Context, req *restful.Request, resp *restful.Response, lang string) {

	var (
		token     dtos.AccessTokenResponse
		errorData common.ErrorData
	)
	refreshToken, _ := req.BodyParameter("refresh_token")
	tokenSvc := UserTokenService{}
	t, _ := tokenSvc.GetUserTokensByreRefreshToken(ctx, refreshToken)
	if t.ID > 0 {
		userSvc := UserService{}
		short, _ := userSvc.GetUserByID(ctx, t.UserId)
		token, errorData = svc.GenerateTokenResponse(ctx, false, t.ClientId, short)
		if errorData.IsNotNil() {
			config.Logger.Error(errorData.Err)
			errorData.Lang = lang
			common.ResponseErrorMessage(ctx, req, resp, config.Bundle, errorData)
			return
		}
		common.ResponseSuccess(resp, token)
	}
}
func (svc *OAuthService) GetTokenByClientCredentials(ctx context.Context, req *restful.Request, resp *restful.Response, lang string) {

}
func (svc *OAuthService) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	return
}

// GenerateTokenResponse clientID可能为空，当登录eauth本身时即为空
func (svc *OAuthService) GenerateTokenResponse(ctx context.Context, needMfa bool, clientId string, user dtos.UserDetail) (response dtos.AccessTokenResponse, errorData common.ErrorData) {

	if config.ApplicationConfig.LoginConfig.MFA && needMfa {
		mfaSvc := MultiFactorAuthService{}
		var (
			mfa dtos.MultiFactorAuthDetail
		)
		mfa = mfaSvc.GetUserMultiFactorAuthByUserId(ctx, user.ID)
		if mfa.ID == 0 {
			mfa, _ = mfaSvc.AddMultiFactorAuth(ctx, user.ID)
		}
		response.Mfa = true
		if mfa.Status == "unbound" {
			response.Image = mfa.Image
			response.Secret = mfa.Secret
		}
		response.ID = user.ID
		response.Username = user.Username
		response.Nickname = user.Nickname
		response.Email = user.Email
		response.Phone = user.Phone
		return
	}

	var claimsID string
	claimsID = common.NewID()
	nowTime := time.Now()
	response.AccessToken, errorData = svc.newAccessToken(ctx, clientId, nowTime, claimsID, user)
	if errorData.IsNotNil() {
		config.Logger.Errorf("create access token failed for  user: %s, err: %s", user.Username, errorData.Err.Error())
		return
	}
	response.IDToken, errorData = svc.newIDToken(ctx, clientId, nowTime, claimsID, user)
	if errorData.IsNotNil() {
		config.Logger.Errorf("create access token failed for  user: %s, err: %s", user.Username, errorData.Err.Error())
		return
	}
	response.RefreshToken = common.NewSecureID(20)
	response.ExpiresIn = int64(config.ApplicationConfig.TokenPeriod) * 3600
	response.TokenType = "bearer"
	response.Timestamp = nowTime.Unix() + int64(config.ApplicationConfig.TokenPeriod)*3600
	//todo store  user token
	if errorData.IsNil() {
		token := dtos.UserTokenCreate{
			UserId:       user.ID,
			ClientId:     clientId,
			Expired:      response.Timestamp,
			ExpiredTime:  time.Unix(response.Timestamp, 0),
			RefreshToken: response.RefreshToken,
			ClaimsID:     claimsID,
			SessionKey:   "",
			Token:        response.AccessToken,
		}
		tokenSvc := UserTokenService{}
		tokenSvc.AddUserToken(ctx, token)
	}
	return response, errorData
}

func (svc *OAuthService) newAccessToken(ctx context.Context, clientId string, nowTime time.Time, claimsID string, user dtos.UserDetail) (token string, errorData common.ErrorData) {
	token, errorData = svc.newIDToken(ctx, clientId, nowTime, claimsID, user)
	return token, errorData
}
func (svc *OAuthService) GenerateBaseClaims(ctx context.Context, nowTime time.Time, claimsID string, user dtos.UserDetail) (claims dtos.UserClaims, errorData common.ErrorData) {
	expireTime := nowTime.Add(time.Duration(config.ApplicationConfig.TokenPeriod) * time.Hour)
	claims.Username = user.Username
	claims.Nickname = user.Nickname
	claims.Role = user.Role
	claims.Id = user.ID
	claims.Email = user.Email
	claims.Phone = user.Phone
	claims.ID = claimsID
	claims.Subject = fmt.Sprintf("%d", user.ID)
	claims.Audience = []string{}
	claims.ExpiresAt = jwt.NewNumericDate(expireTime)
	claims.NotBefore = jwt.NewNumericDate(nowTime)
	claims.IssuedAt = jwt.NewNumericDate(nowTime)
	return
}
func (svc *OAuthService) newIDToken(ctx context.Context, clientId string, nowTime time.Time, claimsID string, user dtos.UserDetail) (token string, errorData common.ErrorData) {
	var claims dtos.UserClaims
	claims, errorData = svc.GenerateBaseClaims(ctx, nowTime, claimsID, user)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return token, errorData
	}
	claims.Issuer = config.ApplicationConfig.ServerAddress
	claims.ID = claimsID
	claims.Audience = append(claims.Audience, clientId)
	if clientId != config.ApplicationName {
		claims.Audience = append(claims.Audience, config.ApplicationName)
	}
	//使用第一个证书
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	var key *rsa.PrivateKey
	sys := ConfigService{}
	_, private, _ := sys.GetRsaKeys(ctx)
	key, errorData.Err = jwt.ParseRSAPrivateKeyFromPEM([]byte(private))
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return token, errorData
	}
	token, errorData.Err = tok.SignedString(key)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return token, errorData
	}
	return token, errorData
}

func (svc *OAuthService) GetPublicKeys(ctx context.Context) (jwks jose.JSONWebKeySet, errorData common.ErrorData) {
	sys := ConfigService{}
	var public string
	public, _, errorData = sys.GetRsaKeys(ctx)
	if errorData.IsNotNil() {
		config.Logger.Errorf("get certs failed, err: %s", errorData.Err.Error())
		return jwks, errorData
	}
	jwks = jose.JSONWebKeySet{}
	//follows the protocol rfc 7517(draft)
	//link here: https://self-issued.info/docs/draft-ietf-jose-json-web-key.html
	//or https://datatracker.ietf.org/doc/html/draft-ietf-jose-json-web-key
	certPemBlock := []byte(public)
	certDerBlock, _ := pem.Decode(certPemBlock)
	x509Cert, _ := x509.ParseCertificate(certDerBlock.Bytes)

	var jwk jose.JSONWebKey
	jwk.Key = x509Cert.PublicKey
	jwk.Certificates = []*x509.Certificate{x509Cert}
	jwk.KeyID = config.ApplicationName
	jwk.Algorithm = "RS256"
	jwk.Use = "sig"
	jwks.Keys = append(jwks.Keys, jwk)
	return
}

func (svc *OAuthService) GetDiscoveryInfo(ctx context.Context) (discovery dtos.OpenIDConfiguration) {
	discovery = dtos.OpenIDConfiguration{
		Issuer:                                 config.ApplicationConfig.ServerAddress,
		AuthorizationEndpoint:                  config.ApplicationConfig.ServerAddress + "/oauth/authorize",
		TokenEndpoint:                          config.ApplicationConfig.ServerAddress + config.V1Prefix + "/token",
		UserinfoEndpoint:                       config.ApplicationConfig.ServerAddress + config.V1Prefix + "/userinfo",
		JwksUri:                                fmt.Sprintf("%s/api/discovery/jwks", config.ApplicationConfig.ServerAddress),
		IntrospectionEndpoint:                  fmt.Sprintf("%s/api/introspect", config.ApplicationConfig.ServerAddress),
		RevocationEndpoint:                     fmt.Sprintf("%s/api/revoke", config.ApplicationConfig.ServerAddress),
		ResponseTypesSupported:                 []string{"code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token", "none"},
		ResponseModesSupported:                 []string{"login", "code", "link"},
		GrantTypesSupported:                    []string{"password", "authorization_code"},
		SubjectTypesSupported:                  []string{"public"},
		IdTokenSigningAlgValuesSupported:       []string{"RS256"},
		ScopesSupported:                        []string{"openid", "email", "profile", "phone", "role"},
		ClaimsSupported:                        []string{"iss", "ver", "sub", "aud", "iat", "exp", "id", "type ", "email", "phone"},
		CodeChallengeMethodsSupported:          []string{"S256", "plain"},
		RequestParameterSupported:              true,
		RequestObjectSigningAlgValuesSupported: []string{"HS256", "HS384", "HS512"},
	}
	return discovery
}
func (svc *OAuthService) LoginByUsername(ctx context.Context, loginParam dtos.LoginByUsername) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		user dtos.UserDetail
	)
	userSvc := UserService{}
	user, errorData = userSvc.GetUserByUsername(ctx, loginParam.Username)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return
	}
	errorData.Err = common.ComparePassword(user.PasswordStore, loginParam.Password, config.PasswordSalt)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		errorData.MsgCode = config.MsgCodeUsernameOrPasswordNotRight
		errorData.Err = errors.New("password or username is not right")
		return
	}
	return svc.GenerateTokenResponse(ctx, true, config.ApplicationName, user)

}

func (svc *OAuthService) MfaValidate(ctx context.Context, model dtos.MfaCode) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		user dtos.UserDetail
		mfa  dtos.MultiFactorAuthDetail
	)

	userSvc := UserService{}
	user, errorData = userSvc.GetUserByID(ctx, model.UserId)
	if errorData.IsNotNil() {
		return
	}
	mfaSvc := MultiFactorAuthService{}
	mfa = mfaSvc.GetUserMultiFactorAuthByUserId(ctx, model.UserId)
	if totp.Validate(model.Code, mfa.Secret) {
		if mfa.Status == "unbound" {
			var status dtos.MultiFactorAuthStatus
			status.Id = mfa.ID
			status.UpdatedAt = time.Now()
			status.Status = "bound"
			go mfaSvc.ChangeStatus(ctx, model.UserId, status)
			go userSvc.UpdateUserMFa(ctx, []uint{model.UserId}, true)
		}
		return svc.GenerateTokenResponse(ctx, false, config.ApplicationName, user)
	}
	errorData.Err = fmt.Errorf("mfa code is expired or not right")
	return

}
func (svc *OAuthService) RegisterByOIDC(ctx context.Context, model dtos.RegisterByOIDC) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		temp       dtos.UserAuthProfileTempDetail
		userCreate dtos.UserCreate
		user       dtos.UserDetail
	)
	errorData.Err = model.Validate(ctx)
	if errorData.IsNotNil() {
		config.Logger.Error(errorData.Err)
		return
	}
	tempSvc := UserAuthProfileTempService{}
	temp, errorData = tempSvc.GetUserAuthProfileTempByCode(ctx, model.Code)
	if errorData.IsNotNil() {
		return
	}
	userSvc := UserService{}
	userCreate.Enable = true
	userCreate.Username = model.Username
	userCreate.Role = config.RoleNone
	userCreate.Nickname = model.Nickname
	userCreate.Password = model.Password
	userCreate.Email = model.Email
	userCreate.Phone = model.Phone
	userCreate.Avatar = temp.Avatar
	user, errorData = userSvc.AddUser(ctx, userCreate)
	if errorData.IsNotNil() {
		return
	}
	authProfileCreate := dtos.UserAuthProfileCreate{
		UserId:     user.ID,
		Provider:   temp.Provider,
		LoginID:    temp.LoginID,
		LoginName:  temp.LoginName,
		Nickname:   temp.Nickname,
		Enable:     true,
		Avatar:     temp.Avatar,
		Properties: temp.Properties,
		Home:       temp.Home,
	}
	authProfileSvc := UserAuthProfileService{}
	_, errorData = authProfileSvc.AddUserAuthProfile(ctx, authProfileCreate)
	if errorData.IsNotNil() {
		return
	}
	return svc.GenerateTokenResponse(ctx, true, config.ApplicationName, user)
}
func (svc *OAuthService) LoginByOIDC(ctx context.Context, loginParam dtos.LoginByOIDC) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		idProvider    providers.IdProvider
		provider      dtos.ProviderOidcDetail
		oauth2Token   *oauth2.Token
		profile       dtos.ThirdAuthProfile
		authedProfile dtos.UserAuthProfileDetail
		user          dtos.UserDetail
	)

	providerSvc := ProviderOidcService{}
	provider, errorData = providerSvc.GetProviderOidcByCategory(ctx, loginParam.Provider)
	if errorData.IsNotNil() {
		return
	}
	if !provider.Enable {
		errorData.Err = fmt.Errorf("provicer: %s is disabled, can't use it", provider.Name)
		return
	}

	switch loginParam.Provider {
	case config.ProviderOidcGithub:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewGithubProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcGitlab:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewGitlabProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcGoogle:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewGoogleProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcMicrosoft:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewMicrosoftProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcWechat:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewWechatProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcWechatWork:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewWechatWorkProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcDingTalk:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewDingTalkProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcGitee:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewGiteeProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcWeibo:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewWeiboProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcQq:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewQQProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcFeiShu:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewFeiShuProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcBilibili:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewBilibiliProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcTiktok:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewTiktokProvider(ctx, loginParam.RedirectUri, providerDetail)
	case config.ProviderOidcBaidu:
		var providerDetail dtos.ProviderOidcDetail
		copyByJSON(provider, &providerDetail)
		idProvider = providers.NewBaiduProvider(ctx, loginParam.RedirectUri, providerDetail)
	default:
		errorData.Err = fmt.Errorf("provider: %s is not supported", loginParam.Provider)
		return
	}
	oauth2Token, errorData = idProvider.GetToken(ctx, loginParam.Code)
	if errorData.IsNotNil() {
		return
	}
	profile, errorData = idProvider.GetUserInfo(ctx, oauth2Token)
	if errorData.IsNotNil() {
		return
	}
	authProfileSvc := UserAuthProfileService{}
	authedProfile, errorData = authProfileSvc.GetUserAuthProfileByProviderAndId(ctx, profile.Provider, profile.LoginID)
	if errorData.IsNotNil() {
		return
	}
	userSvc := UserService{}

	//判断是否存在以及是否允许自注册
	if authedProfile.ID > 0 {
		if !authedProfile.Enable {
			errorData.Err = fmt.Errorf("forbiden authed by: %s", loginParam.Provider)
			return
		}
		// 生成用户的第三方认证信息
		authedProfile.Nickname = profile.Nickname
		authedProfile.Avatar = profile.Avatar
		authedProfile.Home = profile.Home
		authedProfile.Enable = true
		authedProfile.Properties = profile.Properties
		var m dtos.UserAuthProfileUpdate
		copyByJSON(authedProfile, &m)
		authProfileSvc.UpdateUserAuthProfile(ctx, m)
		user, errorData = userSvc.GetUserByID(ctx, authedProfile.UserId)
		if errorData.IsNotNil() {
			return
		}
	}

	return svc.GenerateTokenResponse(ctx, true, config.ApplicationName, user)
}

func (svc *OAuthService) LoginByLDAP(ctx context.Context, loginParam dtos.LoginByLDAP) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	return
}
func (svc *OAuthService) ThirdAuthMethods(ctx context.Context) (response dtos.ThirdAuthMethod) {
	oidcSvc := ProviderOidcService{}
	response.MFA = config.ApplicationConfig.LoginConfig.MFA
	response.FaceRecognition = config.ApplicationConfig.LoginConfig.FaceRecognition
	result, _ := oidcSvc.ListProviderOidc(ctx, 0, 1000, "", "", nil)
	for _, item := range result.Data {
		if item.Enable {
			var oidc = dtos.OIDC{
				Name:     item.Name,
				Address:  item.AuthorizationEndpoint,
				Category: item.Category,
			}
			params := url.Values{}
			normalizedScopes := make([]string, 0, len(item.Scopes))
			for _, scope := range item.Scopes {
				if strings.HasPrefix(scope, "agentid:") && item.Category == config.ProviderOidcWechatWork {
					agentID := strings.TrimPrefix(scope, "agentid:")
					if len(agentID) > 0 {
						params.Add("agentid", agentID)
					}
					continue
				}
				normalizedScopes = append(normalizedScopes, scope)
			}
			switch item.Category {
			case config.ProviderOidcWechat, config.ProviderOidcWechatWork:
				params.Add("appid", item.ClientId)
			case config.ProviderOidcFeiShu:
				params.Add("app_id", item.ClientId)
			case config.ProviderOidcTiktok:
				params.Add("client_key", item.ClientId)
			case config.ProviderOidcBilibili:
				params.Add("appkey", item.ClientId)
			default:
				params.Add("client_id", item.ClientId)
			}
			params.Add("response_type", "code")
			params.Add("state", common.NewSecureID(16))
			params.Add("redirect_uri", fmt.Sprintf("%s/oauth/callback/%s", strings.TrimSuffix(config.ApplicationConfig.ServerAddress, "/"), item.Category))
			if len(normalizedScopes) > 0 {
				params.Add("scope", strings.Join(normalizedScopes, " "))
			}
			if strings.Contains(item.AuthorizationEndpoint, "?") {
				oidc.Address = item.AuthorizationEndpoint + "&" + params.Encode()
			} else {
				oidc.Address = item.AuthorizationEndpoint + "?" + params.Encode()
			}
			if item.Category == config.ProviderOidcWechat && !strings.Contains(oidc.Address, "#wechat_redirect") {
				oidc.Address += "#wechat_redirect"
			}
			response.Oidcs = append(response.Oidcs, oidc)

		}
	}
	return
}
