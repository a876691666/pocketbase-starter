package proxy

import (
	api "pb/core/api"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var proxyPath = "/proxy"

// ProxyModule 是代理API的模块。
var ProxyModule = &api.Module{}

// Init 初始化代理模块。
func Init(app *pocketbase.PocketBase) {

	handler := func(c *core.RequestEvent) error {
		path := c.Request.URL.Path

		// 裁剪path
		path = strings.TrimPrefix(path, proxyPath)

		record := FindProxyByPath(app, path)

		if record == nil {
			return c.JSON(404, "not found")
		}

		err := RunProxy(app, c, record)

		if err != nil {
			return c.JSON(500, err)
		}

		return err
	}

	ProxyModule.Register("GET", "/{path...}", handler)
	ProxyModule.Register("POST", "/{path...}", handler)
	ProxyModule.Register("PUT", "/{path...}", handler)
	ProxyModule.Register("DELETE", "/{path...}", handler)
	ProxyModule.Register("PATCH", "/{path...}", handler)
	ProxyModule.Register("OPTIONS", "/{path...}", handler)

	ProxyModule.Init(app, proxyPath)
}
