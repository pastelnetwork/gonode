package database

import (
	"reflect"
	"strings"
)

func processEscapeString(s string) string {
	s = strings.Replace(s, `'`, `''`, -1)
	return s
}

func processInputString(s string) string {
	if s == "" {
		return "NULL"
	}
	return processEscapeString(s)
}

func safeStringQueryStruct(ps interface{}) interface{} {
	v := reflect.ValueOf(ps).Elem() // Elem() dereferences pointer
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.String:
			str := fv.String()
			fv.SetString(processEscapeString(str))
		}
	}
	return ps
}
