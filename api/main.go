package api

import (
	"pb/api/user"

	"github.com/pocketbase/pocketbase"
)

func Init(app *pocketbase.PocketBase) {
	user.Init(app)
}
