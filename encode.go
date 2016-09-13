package linecfg

import (
	"fmt"
	"reflect"
	"strings"
)

// Transforms a struct into a key=value format.
//
// All the keys are lower-cased by default.
func Encode(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return "", ErrInvalidType{reflect.TypeOf(v)}
	}

	el := rv.Elem()

	pairs := make([]string, 0, el.NumField())

	for i := 0; i < el.NumField(); i++ {
		fv := el.Field(i)
		ft := el.Type().Field(i)

		key := ft.Tag.Get("linecfg") // FIXME: split on ","
		if key == "" {
			key = strings.ToLower(ft.Name)
		}

		if fv.Interface() != reflect.Zero(fv.Type()).Interface() {
			pairs = append(pairs, fmt.Sprintf("%s=%v", key, fv.Interface()))
		}
	}

	return strings.Join(pairs, " "), nil
}
