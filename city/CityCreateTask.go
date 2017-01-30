package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CityCreateTaskResult struct {
	app.Result
	City *City `json:"city,omitempty"`
}

type CityCreateTask struct {
	app.Task
	Pid       int64   `json:"pid"`
	Name      string  `json:"name"`      //名称
	Code      string  `json:"code`       //城市代码
	Tags      string  `json:"tags"`      //搜索标签
	Polygon   string  `json:"polygon"`   //区域 lng,lat;lng,lat;lng,lat|lng,lat;lng,lat;lng,lat;
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Result    CityCreateTaskResult
}

func (task *CityCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *CityCreateTask) GetInhertType() string {
	return "city"
}

func (task *CityCreateTask) GetClientName() string {
	return "City.Create"
}
