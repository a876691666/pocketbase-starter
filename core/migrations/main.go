package migrations

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Init(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		record, err := app.FindFirstRecordByData("config", "name", "isInit")

		if err := superadmin(app); err != nil {
			log.Println("初始化失败, 管理员账户创建失败")
		} else {
			log.Println("管理员账户创建成功")
		}

		if err == nil && record.GetString("value") == "true" {
			log.Println("无需初始化")
			return e.Next()
		}

		if err := collections(app); err != nil {
			log.Println("初始化失败, 集合表创建失败")
		} else {
			log.Println("集合表创建成功")
		}

		collection, err := app.FindCollectionByNameOrId("config")
		if err != nil {
			log.Println("config表不存在")
			return err
		}

		record = core.NewRecord(collection)
		record.Set("name", "isInit")
		record.Set("value", "true")

		err = app.Save(record)
		if err != nil {
			log.Println("初始化失败, 配置表更新失败")
			return err
		}

		log.Println("初始化完成")

		return e.Next()
	})
}

func collections(app *pocketbase.PocketBase) error {
	jsonData, err := os.ReadFile("collections.json")
	if err != nil {
		return err
	}
	jsonData = []byte(strings.TrimSpace(string(jsonData)))
	app.ImportCollectionsByMarshaledJSON([]byte(jsonData), true)
	return nil
}

func superadmin(app *pocketbase.PocketBase) error {

	superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return err
	}

	// 从环境变量获取管理员账户信息
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// 设置默认值
	if adminEmail == "" {
		adminEmail = "admin@pb.com"
	}
	if adminPassword == "" {
		adminPassword = "0123456789"
	}

	// 检查是否存在admin账户
	admin, err := app.FindAuthRecordByEmail(superusers, adminEmail)
	if admin == nil || err != nil {
		admin = core.NewRecord(superusers)
		admin.SetEmail(adminEmail)
		admin.SetPassword(adminPassword)

		log.Println("创建管理员账户: ", adminEmail, " / ", adminPassword)

		return app.Save(admin)
	}

	return nil
}
