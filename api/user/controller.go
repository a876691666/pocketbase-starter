package user

import (
	api "pb/core/api"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// UserModule 是用户API的模块。
var UserModule = &api.Module{}

// Init 初始化用户模块。
func Init(app *pocketbase.PocketBase) {
	UserModule.Register("GET", "/all", func(c *core.RequestEvent) error {
		records, err := GetAllUser(app)
		if err != nil {
			return c.JSON(500, err)
		}
		return c.JSON(200, records)
	})

	UserModule.Init(app, "/user")
}
