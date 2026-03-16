package filters

import (
	"strings"

	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	restful "github.com/emicklei/go-restful/v3"
)

// I18n default language is zh
func I18n(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	req.SetAttribute(config.RequestLanguage, common.I18nZH)
	lang := req.Request.FormValue("lang")
	if len(lang) > 0 {
		req.SetAttribute(config.RequestLanguage, strings.ToLower(lang))
	} else {
		accept := req.HeaderParameter("Accept-Language")
		if len(accept) > 0 {
			req.SetAttribute(config.RequestLanguage, strings.ToLower(accept[:2]))
		}
	}
	chain.ProcessFilter(req, resp)
}
