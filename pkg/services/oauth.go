package services

import (
	"bytes"
	"compress/flate"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
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

func (svc *OAuthService) ApplicationAuthCode(ctx context.Context, userId string, model dtos.OidcCodeRequest) (result dtos.OidcCodeResponse, errorData common.ErrorData) {
	var (
		user                          dtos.UserDetail
		app                           dtos.ApplicationDetail
		codeChallenge                 string
		codeChallengeMethodNormalized string
	)
	if model.ResponseType != "code" {
		errorData.Err = fmt.Errorf("response_type must be code ")
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
		errorData.Err = fmt.Errorf("rule: %s redirect_uri: %s is not right, database: %s ", app.RedirectUriMatchType, model.RedirectUri, app.RedirectUri)
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
		Nonce:               model.Nonce,
	}
	appAuthSvc := AppAuthRecordService{}
	_, errorData = appAuthSvc.AddAppAuthRecord(ctx, authApp)

	return
}
func (svc *OAuthService) UpdateUserAvatar(ctx context.Context, userId string, avatarAddress string) (errorData common.ErrorData) {
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
		config.Logger.Errorf("user: %s update failed, err: %s", model.ID, errorData.Err.Error())
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
func (svc *OAuthService) Userinfo(ctx context.Context, userId string) (userinfo dtos.AuthedUserInfo, errorData common.ErrorData) {

	var (
		user dtos.UserDetail
	)
	userSvc := UserService{}
	user, errorData = userSvc.GetUserByID(ctx, userId)
	if errorData.IsNotNil() {
		config.Logger.Errorf("can't get user by id: %s, err: %s", userId, errorData.Err.Error())
		return userinfo, errorData
	}
	if !user.Enable {
		config.Logger.Errorf("can't get user by id: %s, user id disabled", userId)
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

	if len(authRecord.Nonce) > 0 {
		ctx = context.WithValue(ctx, config.RequestNonce, authRecord.Nonce)
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
	if len(t.ID) > 0 {
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
		if len(mfa.ID) == 0 {
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
	claims.Subject = user.ID
	claims.Audience = []string{}
	if nonce := ctx.Value(config.RequestNonce); nonce != nil {
		claims.Nonce = fmt.Sprintf("%v", nonce)
	}
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
			go userSvc.UpdateUserMFa(ctx, []string{model.UserId}, true)
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

	if len(loginParam.BindId) > 0 {
		user, errorData = userSvc.GetUserByID(ctx, loginParam.BindId)
		if errorData.IsNotNil() {
			return
		}
		authProfileSvc.AddUserAuthProfile(ctx, dtos.UserAuthProfileCreate{
			UserId:         user.ID,
			Provider:       loginParam.Provider,
			LoginID:        profile.LoginID,
			LoginName:      profile.Username,
			Nickname:       profile.Nickname,
			Enable:         true,
			Avatar:         profile.Avatar,
			Home:           profile.Home,
			Properties:     profile.Properties,
			LatestUsedTime: time.Now().Format(time.DateTime),
		})
	}

	//判断是否存在以及是否允许自注册
	if len(authedProfile.ID) > 0 {
		user, errorData = userSvc.GetUserByID(ctx, authedProfile.UserId)
		if errorData.IsNotNil() {
			return
		}
		if !authedProfile.Enable {
			errorData.Err = fmt.Errorf("forbiden authed by: %s", loginParam.Provider)
			return
		}
		// 生成用户的第三方认证信息
		authedProfile.Nickname = profile.Nickname
		authedProfile.Avatar = profile.Avatar
		authedProfile.Home = profile.Home
		authedProfile.Enable = true
		authedProfile.UserId = user.ID
		authedProfile.Properties = profile.Properties
		var m dtos.UserAuthProfileUpdate
		copyByJSON(authedProfile, &m)
		authProfileSvc.UpdateUserAuthProfile(ctx, m)
	}

	return svc.GenerateTokenResponse(ctx, true, config.ApplicationName, user)
}

func (svc *OAuthService) LoginBySAML(ctx context.Context, loginParam dtos.LoginBySAML) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	var (
		provider      dtos.ProviderSamlDetail
		authedProfile dtos.UserAuthProfileDetail
		user          dtos.UserDetail
	)

	if len(strings.TrimSpace(loginParam.Provider)) == 0 {
		errorData.Err = fmt.Errorf("provider is required")
		return
	}
	if len(strings.TrimSpace(loginParam.GetSamlResponse())) == 0 {
		errorData.Err = fmt.Errorf("samlResponse is required")
		return
	}

	providerSvc := ProviderSamlService{}
	provider, errorData = providerSvc.GetProviderSamlByCategory(ctx, loginParam.Provider)
	if errorData.IsNotNil() {
		return
	}
	if !provider.Enable {
		errorData.Err = fmt.Errorf("provider: %s is disabled, can't use it", provider.Name)
		return
	}

	assertion, err := parseSAMLAssertion(loginParam.GetSamlResponse())
	if err != nil {
		errorData.Err = err
		return
	}

	expectedIssuer := strings.TrimSpace(provider.EntityId)
	if len(expectedIssuer) > 0 {
		issuerMatched := strings.EqualFold(strings.TrimSpace(assertion.ResponseIssuer), expectedIssuer) || strings.EqualFold(strings.TrimSpace(assertion.AssertionIssuer), expectedIssuer)
		if !issuerMatched {
			errorData.Err = fmt.Errorf("saml issuer is not matched with provider entityId")
			return
		}
	}
	if !assertion.NotOnOrAfter.IsZero() && time.Now().After(assertion.NotOnOrAfter) {
		errorData.Err = fmt.Errorf("saml assertion is expired")
		return
	}

	loginID := selectSAMLAttribute(assertion.Attributes, provider.LoginIDAttr, "NameID", "nameid", "uid", "sub")
	if len(loginID) == 0 {
		loginID = strings.TrimSpace(assertion.NameID)
	}
	if len(loginID) == 0 {
		errorData.Err = fmt.Errorf("saml response has no login id")
		return
	}
	loginName := selectSAMLAttribute(assertion.Attributes, provider.LoginNameAttr, "name", "username", "displayName")
	if len(loginName) == 0 {
		loginName = loginID
	}
	nickname := selectSAMLAttribute(assertion.Attributes, provider.NicknameAttr, "nickname", "displayName")
	if len(nickname) == 0 {
		nickname = loginName
	}
	email := selectSAMLAttribute(assertion.Attributes, provider.EmailAttr, "email", "mail")
	phone := selectSAMLAttribute(assertion.Attributes, provider.PhoneAttr, "phone_number", "phone", "mobile")
	avatar := selectSAMLAttribute(assertion.Attributes, provider.AvatarAttr, "picture", "avatar")

	properties := dtos.JsonMap{}
	for key, value := range assertion.Attributes {
		properties[key] = value
	}
	if len(assertion.NameID) > 0 {
		properties["NameID"] = assertion.NameID
	}
	if len(assertion.ResponseIssuer) > 0 {
		properties["responseIssuer"] = assertion.ResponseIssuer
	}
	if len(assertion.AssertionIssuer) > 0 {
		properties["assertionIssuer"] = assertion.AssertionIssuer
	}
	if len(loginParam.GetRelayState()) > 0 {
		properties["relayState"] = loginParam.GetRelayState()
	}

	profile := dtos.ThirdAuthProfile{
		Provider:   provider.Category,
		LoginID:    loginID,
		Username:   loginName,
		Nickname:   nickname,
		Avatar:     avatar,
		Properties: properties,
	}
	if len(email) > 0 {
		profile.Properties["email"] = email
	}
	if len(phone) > 0 {
		profile.Properties["phone"] = phone
	}

	authProfileSvc := UserAuthProfileService{}
	authedProfile, errorData = authProfileSvc.GetUserAuthProfileByProviderAndId(ctx, profile.Provider, profile.LoginID)
	if errorData.IsNotNil() {
		return
	}

	userSvc := UserService{}
	if len(authedProfile.ID) == 0 {
		errorData.Err = fmt.Errorf("saml user is not bound in eauth, provider: %s, loginId: %s", profile.Provider, profile.LoginID)
		return
	}
	if !authedProfile.Enable {
		errorData.Err = fmt.Errorf("forbiden authed by: %s", loginParam.Provider)
		return
	}

	authedProfile.LoginName = profile.Username
	authedProfile.Nickname = profile.Nickname
	authedProfile.Avatar = profile.Avatar
	authedProfile.Properties = profile.Properties
	authedProfile.Enable = true
	var updateModel dtos.UserAuthProfileUpdate
	copyByJSON(authedProfile, &updateModel)
	authProfileSvc.UpdateUserAuthProfile(ctx, updateModel)
	user, errorData = userSvc.GetUserByID(ctx, authedProfile.UserId)
	if errorData.IsNotNil() {
		return
	}

	return svc.GenerateTokenResponse(ctx, true, config.ApplicationName, user)
}

func (svc *OAuthService) LoginByLDAP(ctx context.Context, loginParam dtos.LoginByLDAP) (response dtos.AccessTokenResponse, errorData common.ErrorData) {
	return
}
func (svc *OAuthService) ThirdAuthMethods(ctx context.Context) (response dtos.ThirdAuthMethod) {
	oidcSvc := ProviderOidcService{}
	samlSvc := ProviderSamlService{}
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
	samlResults, _ := samlSvc.ListProviderSaml(ctx, 1, 1000, "", "", nil)
	for _, item := range samlResults.Data {
		if !item.Enable {
			continue
		}
		address, err := buildSAMLAuthnRequestURL(*item)
		if err != nil {
			config.Logger.Warnf("build saml authn request failed for provider: %s, err: %s", item.Category, err.Error())
			continue
		}
		response.Samls = append(response.Samls, dtos.SAML{
			Name:     item.Name,
			Address:  address,
			Category: item.Category,
		})
	}
	return
}

type samlParsedAssertion struct {
	NameID          string
	ResponseIssuer  string
	AssertionIssuer string
	Attributes      map[string]string
	NotOnOrAfter    time.Time
}

type samlResponseEnvelope struct {
	XMLName    xml.Name        `xml:"Response"`
	Issuer     string          `xml:"Issuer"`
	Assertions []samlAssertion `xml:"Assertion"`
}

type samlAssertion struct {
	Issuer              string                   `xml:"Issuer"`
	Subject             samlSubject              `xml:"Subject"`
	Conditions          samlConditions           `xml:"Conditions"`
	AttributeStatements []samlAttributeStatement `xml:"AttributeStatement"`
}

type samlSubject struct {
	NameID               samlNameID                `xml:"NameID"`
	SubjectConfirmations []samlSubjectConfirmation `xml:"SubjectConfirmation"`
}

type samlNameID struct {
	Value string `xml:",chardata"`
}

type samlSubjectConfirmation struct {
	SubjectConfirmationData samlSubjectConfirmationData `xml:"SubjectConfirmationData"`
}

type samlSubjectConfirmationData struct {
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
}

type samlConditions struct {
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
}

type samlAttributeStatement struct {
	Attributes []samlAttribute `xml:"Attribute"`
}

type samlAttribute struct {
	Name         string               `xml:"Name,attr"`
	FriendlyName string               `xml:"FriendlyName,attr"`
	Values       []samlAttributeValue `xml:"AttributeValue"`
}

type samlAttributeValue struct {
	Value string `xml:",chardata"`
}

type samlAuthnRequest struct {
	XMLName                     xml.Name          `xml:"samlp:AuthnRequest"`
	XMLNSSAMLP                  string            `xml:"xmlns:samlp,attr"`
	XMLNSSAML                   string            `xml:"xmlns:saml,attr"`
	ID                          string            `xml:"ID,attr"`
	Version                     string            `xml:"Version,attr"`
	IssueInstant                string            `xml:"IssueInstant,attr"`
	Destination                 string            `xml:"Destination,attr,omitempty"`
	AssertionConsumerServiceURL string            `xml:"AssertionConsumerServiceURL,attr,omitempty"`
	ProtocolBinding             string            `xml:"ProtocolBinding,attr,omitempty"`
	Issuer                      samlAuthnIssuer   `xml:"saml:Issuer"`
	NameIDPolicy                *samlNameIDPolicy `xml:"samlp:NameIDPolicy,omitempty"`
}

type samlAuthnIssuer struct {
	Value string `xml:",chardata"`
}

type samlNameIDPolicy struct {
	AllowCreate string `xml:"AllowCreate,attr,omitempty"`
}

func buildSAMLAuthnRequestURL(provider dtos.ProviderSamlDetail) (address string, err error) {
	request := samlAuthnRequest{
		XMLNSSAMLP:                  "urn:oasis:names:tc:SAML:2.0:protocol",
		XMLNSSAML:                   "urn:oasis:names:tc:SAML:2.0:assertion",
		ID:                          "_" + common.NewSecureID(16),
		Version:                     "2.0",
		IssueInstant:                time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Destination:                 strings.TrimSpace(provider.SsoURL),
		AssertionConsumerServiceURL: strings.TrimSpace(provider.AcsURL),
		ProtocolBinding:             "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
		Issuer: samlAuthnIssuer{
			Value: strings.TrimSpace(provider.AcsURL),
		},
		NameIDPolicy: &samlNameIDPolicy{
			AllowCreate: "true",
		},
	}
	if len(strings.TrimSpace(provider.EntityId)) > 0 {
		request.Issuer.Value = strings.TrimSpace(provider.EntityId)
	}

	rawXML, marshalErr := xml.Marshal(request)
	if marshalErr != nil {
		return "", marshalErr
	}
	var compressed bytes.Buffer
	writer, writerErr := flate.NewWriter(&compressed, flate.DefaultCompression)
	if writerErr != nil {
		return "", writerErr
	}
	_, writeErr := writer.Write(rawXML)
	if writeErr != nil {
		_ = writer.Close()
		return "", writeErr
	}
	if closeErr := writer.Close(); closeErr != nil {
		return "", closeErr
	}

	params := url.Values{}
	params.Set("SAMLRequest", base64.StdEncoding.EncodeToString(compressed.Bytes()))
	params.Set("RelayState", common.NewSecureID(16))

	if strings.Contains(provider.SsoURL, "?") {
		return provider.SsoURL + "&" + params.Encode(), nil
	}
	return provider.SsoURL + "?" + params.Encode(), nil
}

func parseSAMLAssertion(encodedResponse string) (result samlParsedAssertion, err error) {
	encodedResponse = strings.TrimSpace(encodedResponse)
	if len(encodedResponse) == 0 {
		err = fmt.Errorf("saml response is empty")
		return
	}
	decoded, decodeErr := base64.StdEncoding.DecodeString(encodedResponse)
	if decodeErr != nil {
		err = fmt.Errorf("decode saml response failed: %w", decodeErr)
		return
	}
	var envelope samlResponseEnvelope
	if unmarshalErr := xml.Unmarshal(decoded, &envelope); unmarshalErr != nil {
		err = fmt.Errorf("parse saml response xml failed: %w", unmarshalErr)
		return
	}
	if len(envelope.Assertions) == 0 {
		err = fmt.Errorf("saml response has no assertion")
		return
	}

	assertion := envelope.Assertions[0]
	result.ResponseIssuer = strings.TrimSpace(envelope.Issuer)
	result.AssertionIssuer = strings.TrimSpace(assertion.Issuer)
	result.NameID = strings.TrimSpace(assertion.Subject.NameID.Value)
	result.Attributes = map[string]string{}
	for _, statement := range assertion.AttributeStatements {
		for _, attr := range statement.Attributes {
			value := ""
			for _, item := range attr.Values {
				if len(strings.TrimSpace(item.Value)) > 0 {
					value = strings.TrimSpace(item.Value)
					break
				}
			}
			if len(value) == 0 {
				continue
			}
			name := strings.TrimSpace(attr.Name)
			if len(name) > 0 {
				result.Attributes[name] = value
			}
			friendlyName := strings.TrimSpace(attr.FriendlyName)
			if len(friendlyName) > 0 {
				result.Attributes[friendlyName] = value
			}
		}
	}

	result.NotOnOrAfter = parseSAMLTime(assertion.Conditions.NotOnOrAfter)
	if result.NotOnOrAfter.IsZero() {
		for _, item := range assertion.Subject.SubjectConfirmations {
			result.NotOnOrAfter = parseSAMLTime(item.SubjectConfirmationData.NotOnOrAfter)
			if !result.NotOnOrAfter.IsZero() {
				break
			}
		}
	}
	return
}

func selectSAMLAttribute(attributes map[string]string, keys ...string) string {
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if len(key) == 0 {
			continue
		}
		if strings.EqualFold(key, "nameid") {
			if value, ok := attributes["NameID"]; ok && len(strings.TrimSpace(value)) > 0 {
				return strings.TrimSpace(value)
			}
			continue
		}
		if value, ok := attributes[key]; ok && len(strings.TrimSpace(value)) > 0 {
			return strings.TrimSpace(value)
		}
		for attrKey, attrValue := range attributes {
			if strings.EqualFold(strings.TrimSpace(attrKey), key) && len(strings.TrimSpace(attrValue)) > 0 {
				return strings.TrimSpace(attrValue)
			}
		}
	}
	return ""
}

func parseSAMLTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999Z",
		"2006-01-02T15:04:05Z",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
