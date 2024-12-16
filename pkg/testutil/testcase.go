package testutil

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func Name(obj ...any) string {
	var s string
	for i, o := range obj {
		if i == 0 {
			if val, ok := o.(int); ok {
				s += fmt.Sprintf("%02d", val) + " "
				continue
			}
		}
		s += toString(o) + " "
	}
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[_\s]+`).ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "/", "%")
	if len(s) > 64 {
		return s[:61] + "..."
	}
	return s
}

func toString(obj any) string {
	v := reflect.ValueOf(obj)
	var sb strings.Builder
	switch v.Kind() {
	case reflect.Ptr:
		sb.WriteString(toString(v.Elem()))
	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			sb.WriteString(fmt.Sprintf("%v ", toString(value.Interface())))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			sb.WriteString(fmt.Sprintf("%v ", toString(value.Interface())))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)
			if field.IsZero() || fieldType.PkgPath != "" {
				continue
			}
			sb.WriteString(toString(field.Interface()) + " ")
		}
	case reflect.Invalid:
		sb.WriteString("")
	default:
		sb.WriteString(fmt.Sprintf("%v", obj))
	}
	return strings.Trim(sb.String(), " ")
}
