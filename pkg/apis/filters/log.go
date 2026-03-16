package filters

import (
	"github.com/efucloud/eauth/pkg/config"
	restful "github.com/emicklei/go-restful/v3"
)

func Log(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
	config.Logger.Debugf("request uri: %s response code: %d", req.Request.URL.Path, resp.StatusCode())
}
