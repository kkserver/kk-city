package city

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"time"
)

type City struct {
	Id        int64   `json:"id"`
	Pid       int64   `json:"pid"`       //上级城市
	Name      string  `json:"name"`      //名称
	Path      string  `json:"path"`      // pid/pid
	Code      string  `json:"code`       //城市代码
	Tags      string  `json:"tags"`      //搜索标签
	Polygon   string  `json:"polygon"`   //区域 lng,lat;lng,lat;lng,lat|lng,lat;lng,lat;lng,lat;
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Oid       int64   `json:"oid"`
}

type Available struct {
	Id     int64  `json:"id"`
	Alias  string `json:"alias"`
	CityId int64  `json:"cityId"`
	Oid    int64  `json:"oid"`
}

type ICityApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetCityTable() *kk.DBTable
	GetAvailableTable() *kk.DBTable
}

type CityApp struct {
	app.App

	DB *app.DBConfig

	Remote *remote.Service

	City *CityService

	CityTable      kk.DBTable
	AvailableTable kk.DBTable
}

const twepoch = int64(1424016000000)

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func NewOid() int64 {
	return milliseconds() - twepoch
}

func (C *CityApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *CityApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *CityApp) GetCityTable() *kk.DBTable {
	return &C.CityTable
}

func (C *CityApp) GetAvailableTable() *kk.DBTable {
	return &C.AvailableTable
}
