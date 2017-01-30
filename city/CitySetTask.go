package city

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CitySetTaskResult struct {
	app.Result
	City *City `json:"city,omitempty"`
}

type CitySetTask struct {
	app.Task
	Id        int64       `json:"id"`
	Pid       interface{} `json:"pid"`
	Name      interface{} `json:"name"`      //名称
	Code      interface{} `json:"code`       //城市代码
	Tags      interface{} `json:"tags"`      //搜索标签
	Polygon   interface{} `json:"polygon"`   //区域 lng,lat;lng,lat;lng,lat|lng,lat;lng,lat;lng,lat;
	Longitude interface{} `json:"longitude"` //经度
	Latitude  interface{} `json:"latitude"`  //纬度
	Result    CitySetTaskResult
}

func (task *CitySetTask) GetResult() interface{} {
	return &task.Result
}

func (task *CitySetTask) GetInhertType() string {
	return "city"
}

func (task *CitySetTask) GetClientName() string {
	return "City.Set"
}
