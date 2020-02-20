package utility

import (
	"errors"
	"fmt"
	"reflect"
)

// InSlice ... return boolean by checking if value in slice
func InSlice(val interface{}, slice interface{}) (exists bool, index int, err error) {
	exists = false
	index = -1

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	default:
		m := fmt.Sprintf("second argument to InSlice must be slice; have %T", slice)
		err = errors.New(m)
		return
	}

	return
}
