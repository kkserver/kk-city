package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityBatchCreateTaskResult struct {
	app.Result
	Citys []City `json:"citys,omitempty"`
}

type CityBatchCreateTask struct {
	app.Task
	Pid    int64  `json:"id"`
	Names  string `json:"names"` //名称,名称,名称,名称
	Result CityBatchCreateTaskResult
}

func (task *CityBatchCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityBatchCreateTask) GetInhertType() string {
	return "city"
}

func (task *CityBatchCreateTask) GetClientName() string {
	return "City.BatchCreate"
}
