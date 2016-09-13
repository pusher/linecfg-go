package linecfg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Handler is the interface implemented by objects that accept line
// key-value pairs. HandleEnvline must copy the data if it
// wishes to retain the data after returning.
type Handler interface {
	HandleLinePair(key, val string) error
}

// Parses and extract the given line into `v'
//
// v must be a pointer to a struct or implement the Handler interface
func Decode(line string, v interface{}) (err error) {
	h, ok := v.(Handler)
	if !ok {
		h, err = NewStructHandler(v)
		if err != nil {
			return err
		}
	}
	return Scanner(line, h)
}

// Extracts the value from the given environment variable into the `v'
//
// v must be a pointer to a struct or implement the Handler interface
func Getenv(key string, v interface{}) (err error) {
	return Decode(os.Getenv(key), v)
}

// The default handler used by `Decode'
type StructHandler struct {
	rv reflect.Value
}

func NewStructHandler(v interface{}) (Handler, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, ErrInvalidType{reflect.TypeOf(v)}
	}
	return &StructHandler{rv: rv}, nil
}

func (h *StructHandler) HandleLinePair(key, val string) error {
	el := h.rv.Elem()
	found := false
	for i := 0; i < el.NumField(); i++ {
		fv := el.Field(i)
		ft := el.Type().Field(i)
		switch {
		case ft.Name == key:
		case ft.Tag.Get("linecfg") == key: // FIXME: split on ","
		case strings.EqualFold(ft.Name, key):
		default:
			continue
		}
		found = true
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				t := fv.Type().Elem()
				v := reflect.New(t)
				fv.Set(v)
				fv = v
			}
			fv = fv.Elem()
		}
		switch fv.Interface().(type) {
		case time.Duration:
			d, err := time.ParseDuration(val)
			if err != nil {
				return &ErrUnmarshalType{val, fv.Type()}
			}
			fv.Set(reflect.ValueOf(d))
		case string:
			fv.SetString(val)
		case bool:
			fv.SetBool(true)
		default:
			switch {
			case reflect.Int <= fv.Kind() && fv.Kind() <= reflect.Int64:
				v, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				fv.SetInt(v)
			case reflect.Uint32 <= fv.Kind() && fv.Kind() <= reflect.Uint64:
				v, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					return err
				}
				fv.SetUint(v)
			case reflect.Float32 <= fv.Kind() && fv.Kind() <= reflect.Float64:
				v, err := strconv.ParseFloat(val, 10)
				if err != nil {
					return err
				}
				fv.SetFloat(v)
			default:
				return ErrUnmarshalType{val, fv.Type()}
			}
		}

	}
	if !found {
		return ErrKeyNotFound{key}
	}
	return nil
}

// A very simple scanner that passes `key=value` pairs back to the handler.
//
// Used by the Decode() function
func Scanner(line string, h Handler) error {
	fields := strings.Fields(line)

	for _, field := range fields {
		kv := strings.SplitN(field, "=", 2)
		if len(kv) != 2 {
			return ErrBadField{line, field}
		}
		err := h.HandleLinePair(kv[0], kv[1])
		if err != nil {
			return err
		}
	}

	return nil
}

type ErrBadField struct {
	Line  string
	Field string
}

func (e ErrBadField) Error() string {
	return fmt.Sprintf("linecfg: bad field '%s' in '%s'", e.Field, e.Line)
}

type ErrKeyNotFound struct {
	Key string
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("linecfg: unknown field '%s'", e.Key)
}

type ErrInvalidType struct {
	Type reflect.Type
}

func (e ErrInvalidType) Error() string {
	if e.Type == nil {
		return "linecfg: invalid type: nil"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "linecfg: invalid type: non-pointer " + e.Type.String()
	}
	return "linecfg: invalid type: nil for " + e.Type.String()
}

// An UnmarshalTypeError describes a logfmt value that was
// not appropriate for a value of a specific Go type.
type ErrUnmarshalType struct {
	Value string       // the logfmt value
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e ErrUnmarshalType) Error() string {
	return "linecfg: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
