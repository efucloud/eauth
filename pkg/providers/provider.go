package providers

import (
	"context"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/models/dtos"
	"golang.org/x/oauth2"
)

type IdProvider interface {
	GetToken(ctx context.Context, code string) (*oauth2.Token, common.ErrorData)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (dtos.ThirdAuthProfile, common.ErrorData)
}

// Scopes represents additional data requested by the clients about the end user.
type Scopes struct {
	//The client has requested a refresh token from the server.
	OfflineAccess bool

	//The client has requested group information about the end user.
	Groups bool
}

// Identity represents the ID Token claims supported by the server.
type Identity struct {
	UserID            string
	Username          string
	PreferredUsername string
	Email             string
	EmailVerified     bool

	Groups []string

	//ConnectorData holds data used by the connector for subsequent requests after initial
	//authentication, such as access tokens for upstream provides.
	//
	//This data is never shared with end users, OAuth clients, or through the API.
	ConnectorData []byte
}

//
//const (
//	scopeOfflineAccess     = "offline_access" //Request a refresh token.
//	scopeOpenID            = "openid"
//	scopeGroups            = "groups"
//	scopeEmail             = "email"
//	scopeProfile           = "profile"
//	scopeFederatedID       = "federated:id"
//	scopeCrossClientPrefix = "audience:server:client_id:"
//)
//
//func  parseScopes(scopes []string) Scopes {
//	var s Scopes
//	for _, scope := range scopes {
//		switch scope {
//		case scopeOfflineAccess:
//			s.OfflineAccess = true
//		case scopeGroups:
//			s.Groups = true
//		}
//	}
//	return  s
//}
