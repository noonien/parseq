package parseq

import (
	"bytes"
	"errors"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Errors map[string]error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	keys := []string{}
	for key := range e {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for i, key := range keys {
		if i > 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(key + ": " + e[key].Error())
	}

	return buf.String()
}

func Unmarshal(q url.Values, v interface{}) error {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("parseq: must be pointer")
	}
	rv = rv.Elem()

	if rv.Kind() != reflect.Struct {
		return errors.New("parseq: must be struct")
	}

	errs := make(Errors)

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := rv.Field(i)

		name := fieldName(t.Field(i))
		if name == "-" {
			continue
		}

		val, ok := q[name]
		if !ok {
			continue
		}

		err := unmarshal(val, f)
		if err != nil {
			errs[name] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func fieldName(f reflect.StructField) string {
	name, ok := nameFromTag(f, "query")
	if ok && len(name) > 0 {
		return name
	}

	name, ok = nameFromTag(f, "json")
	if ok && len(name) > 0 {
		return name
	}

	return f.Name
}

func nameFromTag(f reflect.StructField, tag string) (string, bool) {
	tag, ok := f.Tag.Lookup("query")
	if !ok {
		return "", false
	}

	parts := strings.SplitN(tag, ",", 2)
	return strings.TrimSpace(parts[0]), true
}

func unmarshal(val []string, f reflect.Value) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(val[0])

	case reflect.Int, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val[0], 10, 64)
		if err != nil {
			return errors.New("not a valid integer")
		}
		f.SetInt(i)

	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(val[0], 10, 64)
		if err != nil {
			return errors.New("not a valid unsigned integer")
		}
		f.SetUint(u)

	case reflect.Slice:
		slice := reflect.MakeSlice(f.Type().Elem(), len(val), len(val))

		for i, item := range val {
			item = strings.TrimSpace(item)
			err := unmarshal([]string{item}, slice.Index(i))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
