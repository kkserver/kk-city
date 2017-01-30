package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityAvailableSetTaskResult struct {
	app.Result
}

type CityAvailableSetTask struct {
	app.Task
	Alias   string `json:"alias"`
	CityId  int64  `json:"cityId"`
	CityIds string `json:"cityIds"`
	Result  CityAvailableSetTaskResult
}

func (task *CityAvailableSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityAvailableSetTask) GetInhertType() string {
	return "city"
}

func (task *CityAvailableSetTask) GetClientName() string {
	return "City.AvailableSet"
}
