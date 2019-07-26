package util

import (
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"

	"github.com/shuge/goproxy/whitelist"
)

func ProxyOnReqInWhitelist(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	splits := strings.Split(r.RemoteAddr, ":")
	var clientIpAddr string
	if len(splits) >= 1 {
		clientIpAddr = splits[0]
	}
	if clientIpAddr != "" && !whitelist.InWhitelist(clientIpAddr) {
		return r, goproxy.NewResponse(r,
			goproxy.ContentTypeText,
			http.StatusForbidden,
			"not in whitelist\n",
		)
	}
	return r, nil
}
