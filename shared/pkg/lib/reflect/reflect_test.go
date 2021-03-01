package reflect_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/lib/reflect"
)

type TestStruct struct {
	A bool   `db:"A-db" json:"A-json"`
	B int    `db:"B-db"`
	C string `json:"C-json"`
	D func()
}

func TestReflect_Tags(t *testing.T) {
	assert := require.New(t)

	tags := reflect.Tags(TestStruct{}, "db")
	assert.Equal([]string{"A-db", "B-db"}, tags)

	tags = reflect.Tags(TestStruct{}, "json")
	assert.Equal([]string{"A-json", "C-json"}, tags)
}

func TestReflect_Fields(t *testing.T) {
	assert := require.New(t)

	fields := reflect.Fields(TestStruct{}, "db")
	assert.Equal(map[string]string{"A-db": "A", "B-db": "B"}, fields)

	fields = reflect.Fields(TestStruct{}, "json")
	assert.Equal(map[string]string{"A-json": "A", "C-json": "C"}, fields)
}

func TestReflect_Values(t *testing.T) {
	assert := require.New(t)

	testStructInstance := TestStruct{
		A: true,
		B: 31415,
		C: "hello, world!",
		D: func() {},
	}

	values := reflect.Values(testStructInstance, "db")
	assert.Equal([]interface{}{testStructInstance.A, testStructInstance.B}, values)

	values = reflect.Values(testStructInstance, "json")
	assert.Equal([]interface{}{testStructInstance.A, testStructInstance.C}, values)
}
