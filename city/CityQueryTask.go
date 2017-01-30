package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type CityQueryTaskResult struct {
	app.Result
	Counter *CityQueryCounter `json:"counter,omitempty"`
	Citys   []City            `json:"citys,omitempty"`
}

type CityQueryTask struct {
	app.Task
	Id        int64       `json:"id"`
	Pid       interface{} `json:"pid"`
	Alias     string      `json:"alias"`
	Keyword   string      `json:"q"`
	PageIndex int         `json:"p"`
	PageSize  int         `json:"size"`
	Counter   bool        `json:"counter"`
	Result    CityQueryTaskResult
}

func (task *CityQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityQueryTask) GetInhertType() string {
	return "city"
}

func (task *CityQueryTask) GetClientName() string {
	return "City.Query"
}
