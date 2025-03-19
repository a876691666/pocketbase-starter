package api

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Module struct {
	TestKey  string
	app      *pocketbase.PocketBase
	basePath string
	router   []struct {
		method     string // GET, POST, PUT, PATCH, DELETE, ANY
		path       string
		handler    func(c *core.RequestEvent) error
		middleware []func(c *core.RequestEvent) error
	}
}

func (m *Module) Init(app *pocketbase.PocketBase, basePath string) {
	m.app = app
	m.basePath = basePath
	m.app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		group := se.Router.Group(m.basePath)
		for _, route := range m.router {
			routeItem := group.Route(route.method, route.path, func(c *core.RequestEvent) error {
				return route.handler(c)
			})

			routeItem.BindFunc(func(c *core.RequestEvent) error {
				for _, middleware := range route.middleware {
					if err := middleware(c); err != nil {
						return err
					}
				}
				return c.Next()
			})

			fmt.Println("注册路由 [", routeItem.Method, "]", m.basePath+routeItem.Path)
		}
		return se.Next()
	})
}

func (m *Module) Register(
	method string,
	path string,
	handler func(c *core.RequestEvent) error,
	middleware ...func(c *core.RequestEvent) error,
) {
	m.router = append(m.router, struct {
		method     string
		path       string
		handler    func(c *core.RequestEvent) error
		middleware []func(c *core.RequestEvent) error
	}{
		method:     method,
		path:       path,
		handler:    handler,
		middleware: middleware,
	})
}
