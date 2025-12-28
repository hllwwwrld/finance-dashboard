package tables

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var columnsMap = sync.Map{}

// Columns returns columns
func columns(v any) []string {
	tp := reflect.TypeOf(v)
	name := tp.String()

	columnsMMap, ok := columnsMap.Load(name)
	if ok {
		return columnsMMap.([]string)
	}
	cols := make([]string, 0)

	if tp.Kind() == reflect.Slice {
		tp = tp.Elem()
	}

	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	if tp.Kind() != reflect.Struct {
		panic(errors.New("object must by struct or structPtr, slice of struct, not " + tp.String()))
	}

	m := make(map[string]struct{})
	for idx := 0; idx < tp.NumField(); idx++ {
		f := tp.Field(idx)
		tag := f.Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(f.Name)
		}
		if _, ok = m[tag]; ok {
			panic(fmt.Errorf("duplicate column name: %s. Check srtuct field names and db tags", tag))
		}
		cols = append(cols, tag)
		m[tag] = struct{}{}
	}

	columnsMap.Store(name, cols)

	return cols
}

func scanRows[T any](rows *sql.Rows) (res []*T, err error) {
	res = make([]*T, 0)
	for rows.Next() {
		var current T
		// todo тут наверное надо передавать каждый параметр модельки руками,
		// todo как поинтер, не получится передавать сразу модельки целиком?
		// если так, то надо придумать что-то изящное, типа передаешь структуру, в нее само все парсится
		err = rows.Scan(current)

		if err != nil {
			return nil, fmt.Errorf("GetByUserID.rows.Scan err: %w", err)
		}
		res = append(res, &current)
	}

	return res, nil
}
