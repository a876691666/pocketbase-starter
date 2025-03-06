package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var (
	adminEmail    = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

func main() {
	// 设置默认值
	if adminEmail == "" {
		adminEmail = "admin@pocketbase-starter.com"
	}
	if adminPassword == "" {
		adminPassword = "0123456789"
	}

	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		collection, err := app.FindCollectionByNameOrId("_superusers")
		if err != nil {
			return err
		}

		// 检查是否存在admin账户
		admin, err := app.FindAuthRecordByEmail(collection, adminEmail)

		// 如果没有admin账户,创建一个默认账户
		if admin == nil || err != nil {
			admin = core.NewRecord(collection)
			admin.SetEmail(adminEmail)
			admin.SetPassword(adminPassword)

			if err := app.Save(admin); err != nil {
				return err
			}
			log.Println("已创建默认admin账户: ", adminEmail, " / ", adminPassword)
		}
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
