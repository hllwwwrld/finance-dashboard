package tables

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var columnsMap = sync.Map{}

// Columns returning names of columns by 'db' attribute from go structure, made for parsing rows from postgres/squirrel
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

func allColumnsString(v any) string {
	return strings.Join(columns(v), ", ")
}

func returningAllColumns(v any) string {
	return fmt.Sprintf("RETURNING %s", strings.Join(columns(v), ", "))
}
