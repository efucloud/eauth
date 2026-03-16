package filters

import (
	"github.com/emicklei/go-restful/v3"
	"strings"
)

const bearer = "bearer "

func GetRequestToken(key string, req *restful.Request) (token string) {
	authHeader := req.HeaderParameter(key)
	if len(authHeader) > 0 {
		if len(bearer) > len(authHeader) {
			return
		}
		token = authHeader[len(bearer):]
		tokenPrefix := authHeader[:len(bearer)]
		if strings.ToLower(tokenPrefix) != bearer {
			return
		}
	}
	return token
}
