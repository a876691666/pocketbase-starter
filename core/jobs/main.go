package jobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Job 定义一个任务函数类型
type Job struct {
	Name string
	Cron string
	Func func(app *pocketbase.PocketBase) error
}

// 存储所有注册的任务
var registeredJobs []Job

// RegisterJob 注册一个新的任务
func RegisterJob(job Job) {
	registeredJobs = append(registeredJobs, job)
}

func Init(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		for _, job := range registeredJobs {
			RunJobs(app, job)
		}
		return e.Next()
	})
}

func RunJobs(app *pocketbase.PocketBase, job Job) {
	if job.Cron == "" {
		return
	}

	app.Cron().MustAdd(job.Name, job.Cron, func() {
		job.Func(app)
	})
}
