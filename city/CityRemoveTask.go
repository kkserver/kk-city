package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityRemoveTaskResult struct {
	app.Result
}

type CityRemoveTask struct {
	app.Task
	Id     int64       `json:"id"`
	Pid    interface{} `json:"pid"`
	Result CityRemoveTaskResult
}

func (task *CityRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityRemoveTask) GetInhertType() string {
	return "city"
}

func (task *CityRemoveTask) GetClientName() string {
	return "City.Remove"
}
