package common

import (
	"os"
	"reflect"
)

func ExtractTagFromStruct(s interface{}, tags ...string) map[string]map[string]string {
	t := reflect.TypeOf(s)
	t = t.Elem()
	var result = make(map[string]map[string]string)
	n := t.NumField()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		for _, tag := range tags {
			field.Tag.Get(tag)
			if _, ok := result[field.Name]; !ok {
				result[field.Name] = make(map[string]string)
			}
			result[field.Name][tag] = field.Tag.Get(tag)
		}
	}
	return result
}

func GetEnv(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}
