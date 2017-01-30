package city

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"strings"
)

type CityService struct {
	app.Service

	Create          *CityCreateTask
	BatchCreate     *CityBatchCreateTask
	Set             *CitySetTask
	Get             *CityTask
	Remove          *CityRemoveTask
	Query           *CityQueryTask
	AvailableSet    *CityAvailableSetTask
	AvailableRemove *CityAvailableRemoveTask
}

func (S *CityService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *CityService) HandleCityCreateTask(a ICityApp, task *CityCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := City{}

	v.Name = task.Name
	v.Pid = task.Pid
	v.Code = task.Code
	v.Oid = NewOid()
	v.Tags = task.Tags
	v.Polygon = task.Polygon
	v.Latitude = task.Latitude
	v.Longitude = task.Longitude
	v.Path = "/"

	if v.Pid != 0 {

		t := CityTask{}
		t.Id = v.Pid

		app.Handle(a, &t)

		if t.Result.City == nil {
			task.Result.Errno = ERROR_CITY_NOT_FOUND_PCITY
			task.Result.Errmsg = "Not found parent city"
			return nil
		}

		v.Path = fmt.Sprintf("%s%d/", t.Result.City.Path, v.Pid)
	}

	_, err = kk.DBInsert(db, a.GetCityTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.City = &v

	return nil
}

func (S *CityService) HandleCityBatchCreateTask(a ICityApp, task *CityBatchCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	citys := []City{}

	var path string = "/"

	if task.Pid != 0 {

		t := CityTask{}
		t.Id = task.Pid

		app.Handle(a, &t)

		if t.Result.City == nil {
			task.Result.Errno = ERROR_CITY_NOT_FOUND_PCITY
			task.Result.Errmsg = "Not found parent city"
			return nil
		}

		path = fmt.Sprintf("%s%d/", t.Result.City.Path, t.Result.City.Id)
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		for _, name := range strings.Split(task.Names, ",") {

			if name != "" {

				v := City{}
				v.Name = name
				v.Pid = task.Pid
				v.Path = path

				_, err = kk.DBInsert(db, a.GetCityTable(), a.GetPrefix(), &v)

				if err != nil {
					return err
				}
			}

		}

		return nil

	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Citys = citys

	return nil
}

func (S *CityService) HandleCitySetTask(a ICityApp, task *CitySetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_CITY_NOT_FOUND_ID
		task.Result.Errmsg = "Not found city id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := City{}

	rows, err := kk.DBQuery(db, a.GetCityTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		keys := map[string]bool{}

		if task.Name != nil {
			v.Name = dynamic.StringValue(task.Name, v.Name)
			keys["name"] = true
		}

		if task.Pid != nil {
			var pid = dynamic.IntValue(task.Pid, 0)

			v.Pid = pid

			if pid != 0 {

				t := CityTask{}
				t.Id = v.Pid

				app.Handle(a, &t)

				if t.Result.City == nil {
					task.Result.Errno = ERROR_CITY_NOT_FOUND_PCITY
					task.Result.Errmsg = "Not found parent city"
					return nil
				}

				v.Path = fmt.Sprintf("%s%d/", t.Result.City.Path, v.Pid)
			} else {
				v.Path = "/"
			}

			keys["pid"] = true
			keys["path"] = true
		}

		if task.Code != nil {
			v.Code = dynamic.StringValue(task.Code, v.Code)
			keys["code"] = true
		}

		if task.Tags != nil {
			v.Tags = dynamic.StringValue(task.Tags, v.Tags)
			keys["tags"] = true
		}

		if task.Polygon != nil {
			v.Polygon = dynamic.StringValue(task.Polygon, v.Polygon)
			keys["polygon"] = true
		}

		if task.Longitude != nil {
			v.Longitude = dynamic.FloatValue(task.Longitude, v.Longitude)
			keys["longitude"] = true
		}

		if task.Latitude != nil {
			v.Latitude = dynamic.FloatValue(task.Latitude, v.Latitude)
			keys["latitude"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetCityTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.City = &v

	} else {
		task.Result.Errno = ERROR_CITY_NOT_FOUND
		task.Result.Errmsg = "Not found city"
		return nil
	}

	return nil
}

func (S *CityService) HandleCityRemoveTask(a ICityApp, task *CityRemoveTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id != 0 {

		rows, err := db.Query(fmt.Sprintf("SELECT path FROM %s%s WHERE id=?", a.GetPrefix(), a.GetCityTable().Name), task.Id)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		if rows.Next() {

			var path interface{} = nil

			err = rows.Scan(&path)

			if err != nil {
				task.Result.Errno = ERROR_CITY
				task.Result.Errmsg = err.Error()
				return nil
			}

			_, err = kk.DBDelete(db, a.GetCityTable(), a.GetPrefix(), " WHERE id=? OR path LIKE ?", task.Id, fmt.Sprintf("%s%d/%%", dynamic.StringValue(path, "/"), task.Id))

			if err != nil {
				task.Result.Errno = ERROR_CITY
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else {
			task.Result.Errno = ERROR_CITY_NOT_FOUND
			task.Result.Errmsg = "Not found city"
			return nil
		}

	} else if task.Pid != nil {

		var pid = dynamic.IntValue(task.Pid, 0)

		if pid == 0 {
			_, err = kk.DBDelete(db, a.GetCityTable(), a.GetPrefix(), " WHERE pid=? OR path LIKE '/%'", pid)
			if err != nil {
				task.Result.Errno = ERROR_CITY
				task.Result.Errmsg = err.Error()
				return nil
			}
		} else {

			rows, err := db.Query(fmt.Sprintf("SELECT path FROM %s%s WHERE id=?", a.GetPrefix(), a.GetCityTable().Name), pid)

			if err != nil {
				task.Result.Errno = ERROR_CITY
				task.Result.Errmsg = err.Error()
				return nil
			}

			defer rows.Close()

			if rows.Next() {

				var path interface{} = nil

				err = rows.Scan(&path)

				if err != nil {
					task.Result.Errno = ERROR_CITY
					task.Result.Errmsg = err.Error()
					return nil
				}

				_, err = kk.DBDelete(db, a.GetCityTable(), a.GetPrefix(), " WHERE path LIKE ?", fmt.Sprintf("%s%d/%%", dynamic.StringValue(path, "/"), pid))

				if err != nil {
					task.Result.Errno = ERROR_CITY
					task.Result.Errmsg = err.Error()
					return nil
				}

			} else {
				task.Result.Errno = ERROR_CITY_NOT_FOUND_PCITY
				task.Result.Errmsg = "Not found parent city"
				return nil
			}
		}

	}

	return nil
}

func (S *CityService) HandleCityTask(a ICityApp, task *CityTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}
	v := City{}

	rows, err := kk.DBQuery(db, a.GetCityTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.City = &v

	} else {
		task.Result.Errno = ERROR_CITY_NOT_FOUND
		task.Result.Errmsg = "Not found city"
		return nil
	}

	return nil
}

func (S *CityService) HandleCityQueryTask(a ICityApp, task *CityQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	var citys = []City{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	if task.Alias != "" {
		sql.WriteString(fmt.Sprintf("FROM %s%s as c LEFT JOIN %s%s as a ON c.id=a.cityid WHERE a.alias=?", a.GetPrefix(), a.GetCityTable().Name, a.GetPrefix(), a.GetAvailableTable().Name))
		args = append(args, task.Alias)
	} else {
		sql.WriteString(fmt.Sprintf("FROM %s%s as c WHERE 1", a.GetPrefix(), a.GetCityTable().Name))
	}

	if task.Id != 0 {
		sql.WriteString(" AND c.id=?")
		args = append(args, task.Id)
	}

	if task.Pid != nil {
		sql.WriteString(" AND c.pid=?")
		args = append(args, task.Pid)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (c.tags LIKE ? OR c.name LIKE ?)")
		args = append(args, q, q)
	}

	if task.Alias != "" {
		sql.WriteString(" ORDER BY a.oid ASC,c.oid ASC,c.id ASC")
	} else {
		sql.WriteString(" ORDER BY c.oid ASC,c.id ASC")
	}

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {

		var counter = CityQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize

		rows, err := db.Query("SELECT COUNT(*) "+sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		total := 0

		if rows.Next() {
			err = rows.Scan(&total)
			if err != nil {
				task.Result.Errno = ERROR_CITY
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		if total%pageSize == 0 {
			counter.PageCount = total / pageSize
		} else {
			counter.PageCount = total/pageSize + 1
		}

		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var v = City{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := db.Query("SELECT c.* "+sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}

		citys = append(citys, v)
	}

	task.Result.Citys = citys

	return nil
}

func (S *CityService) HandleCityAvailableSetTask(a ICityApp, task *CityAvailableSetTask) error {

	if task.Alias == "" {
		task.Result.Errno = ERROR_CITY_NOT_FOUND_ALIAS
		task.Result.Errmsg = "Not found alias"
		return nil
	}

	cityIds := []int64{}

	if task.CityId != 0 {
		cityIds = append(cityIds, task.CityId)
	}

	for _, v := range strings.Split(task.CityIds, ",") {
		id := dynamic.IntValue(v, 0)
		if id != 0 {
			cityIds = append(cityIds, id)
		}
	}

	if len(cityIds) == 0 {
		task.Result.Errno = ERROR_CITY_NOT_FOUND_ID
		task.Result.Errmsg = "Not found cityId"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		city := City{}

		scanner := kk.NewDBScaner(&city)

		citys := map[int64]City{}

		for _, cityId := range cityIds {

			city, ok := citys[cityId]

			if !ok {

				rows, err := kk.DBQuery(tx, a.GetCityTable(), a.GetPrefix(), " WHERE id=?", cityId)

				if err != nil {
					return err
				}

				if rows.Next() {

					err = scanner.Scan(rows)

					if err != nil {
						return err
					}

					citys[cityId] = city

				} else {
					rows.Close()
					return app.NewError(ERROR_CITY_NOT_FOUND, "Not found city")
				}

				rows.Close()
			}

			count, err := kk.DBQueryCount(tx, a.GetAvailableTable(), a.GetPrefix(), " WHERE cityid=? AND alias=?", cityId, task.Alias)

			if err != nil {
				return err
			}

			if count == 0 {

				v := Available{}
				v.CityId = cityId
				v.Alias = task.Alias
				v.Oid = city.Oid

				_, err = kk.DBInsert(tx, a.GetAvailableTable(), a.GetPrefix(), &v)

				if err != nil {
					return err
				}
			}
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_CITY
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	return nil
}

func (S *CityService) HandleCityAvailableRemoveTask(a ICityApp, task *CityAvailableRemoveTask) error {

	if task.Alias == "" {
		task.Result.Errno = ERROR_CITY_NOT_FOUND_ALIAS
		task.Result.Errmsg = "Not found alias"
		return nil
	}

	cityIds := []int64{}

	if task.CityId != 0 {
		cityIds = append(cityIds, task.CityId)
	}

	for _, v := range strings.Split(task.CityIds, ",") {
		id := dynamic.IntValue(v, 0)
		if id != 0 {
			cityIds = append(cityIds, id)
		}
	}

	if len(cityIds) == 0 {
		task.Result.Errno = ERROR_CITY_NOT_FOUND_ID
		task.Result.Errmsg = "Not found cityId"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	var sql = bytes.NewBuffer(nil)
	var args = []interface{}{}

	sql.WriteString(" WHERE alias=? AND cityid IN (")

	for i, cityId := range cityIds {
		if i != 0 {
			sql.WriteString(",")
		}
		sql.WriteString("?")
		args = append(args, cityId)
	}

	sql.WriteString(")")

	_, err = kk.DBDelete(db, a.GetAvailableTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_CITY
		task.Result.Errmsg = err.Error()
		return nil
	}

	return nil
}
