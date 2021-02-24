package reflect

import (
	"reflect"
)

type typeKeyPair struct {
	ty  reflect.Type
	key string
}

var tagsMemo = map[typeKeyPair][]string{}

// Tags returns the tags in a structure for a given key. Tags not containing that key are ignored. Tags and Values have the same return order.
func Tags(i interface{}, key string) []string {
	t := reflect.TypeOf(i)

	if tags, ok := tagsMemo[typeKeyPair{t, key}]; ok {
		return tags
	}

	var tags []string
	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup(key); ok {
			tags = append(tags, tag)
		}
	}

	tagsMemo[typeKeyPair{t, key}] = tags
	return tags
}

var fieldsMemo = map[typeKeyPair]map[string]string{}

// Fields returns a mapping of tags to field names, given key. Fields that don't have the provided key are ignored.
func Fields(i interface{}, key string) map[string]string {
	t := reflect.TypeOf(i)

	if fields, ok := fieldsMemo[typeKeyPair{t, key}]; ok {
		return fields
	}

	var fields map[string]string
	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup(key); ok {
			fields[tag] = t.Field(i).Name
		}
	}

	fieldsMemo[typeKeyPair{t, key}] = fields
	return fields
}

// Values returns the values in a structure for a given key. Tags not containing that key are ignored. Tags and Values have the same return order.
func Values(i interface{}, key string) []interface{} {
	v := reflect.ValueOf(i)
	fields := Fields(i, key)

	values := []interface{}{}
	for _, tag := range Tags(i, key) {
		values = append(values, v.FieldByName(fields[tag]).Interface())
	}
	return values
}
