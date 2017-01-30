package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityAvailableRemoveTaskResult struct {
	app.Result
}

type CityAvailableRemoveTask struct {
	app.Task
	Alias   string `json:"alias"`
	CityId  int64  `json:"cityId"`
	CityIds string `json:"cityIds"`
	Result  CityAvailableRemoveTaskResult
}

func (task *CityAvailableRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityAvailableRemoveTask) GetInhertType() string {
	return "city"
}

func (task *CityAvailableRemoveTask) GetClientName() string {
	return "City.AvailableRemove"
}
