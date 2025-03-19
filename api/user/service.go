package user

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func GetAllUser(app *pocketbase.PocketBase) (records []*core.Record, err error) {
	records, err = app.FindAllRecords("users")

	if err != nil {
		return nil, err
	}

	return records, nil
}
