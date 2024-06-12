package caddyrestrictprefix

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(RestrictPrefix{})
}

// RestrictPrefix는 URI의 일부가 주어진 접두사와 일치하는 요청을 제어하는 미들웨어
type RestrictPrefix struct {
	Prefix string `json:"prefix,omitempty"`
	logger *zap.Logger
}

// CaddyModule은 Caddy의 모듈 정보를 반환함
func (RestrictPrefix) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.restrict_prefix",
		New: func() caddy.Module { return new(RestrictPrefix) },
	}
}

// Zap 로거를 RestrictPrefix로 프로비저닝
func (p *RestrictPrefix) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger(p)
	return nil
}

// 모듈 구성에서 접두사를 검증하고 필요시 기본 접두사를 "."로 설정
func (p *RestrictPrefix) Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}
	return nil
}

// ServeHTTP는 caddyhttp.MiddlewareHandler 인터페이스를 구현
func (p RestrictPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request,
	next caddyhttp.Handler) error {
	for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			http.Error(w, "Not Found", http.StatusNotFound)
			if p.logger != nil {
				p.logger.Debug(fmt.Sprintf("restricted prefix: %q in %s", part, r.URL.Path))
			}
			return nil
		}
	}
	return next.ServeHTTP(w, r)
}

// 인터페이스 자체에 변화가 없었는지, 명시적으로 모듈이 기대한 대로
// 인터페이스를 구현하였는지 능동적으로 방어하는 것은 좋은 습관이다.
var (
	_ caddy.Provisioner           = (*RestrictPrefix)(nil)
	_ caddy.Validator             = (*RestrictPrefix)(nil)
	_ caddyhttp.MiddlewareHandler = (*RestrictPrefix)(nil)
)
