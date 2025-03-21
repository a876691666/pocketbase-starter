package api

import (
	"pb/api/proxy"
	"pb/api/user"

	"github.com/pocketbase/pocketbase"
)

func Init(app *pocketbase.PocketBase) {
	user.Init(app)
	proxy.Init(app)
}
