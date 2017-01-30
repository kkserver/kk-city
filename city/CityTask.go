package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityTaskResult struct {
	app.Result
	City *City `json:"city,omitempty"`
}

type CityTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result CityTaskResult
}

func (task *CityTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityTask) GetInhertType() string {
	return "city"
}

func (task *CityTask) GetClientName() string {
	return "City.Get"
}
