package bsonify

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// Accepts a struct or map and returns a bson.M for mongo update $set operation.
// The struct or map can contain pointers, interfaces, maps, and structs.
// Map keys should be strings. BSON tags with omitempty are supported.
func SetUpdateM(v any) bson.M {
	rv := reflect.ValueOf(v)
	dereference(rv)

	switch rv.Kind() {
	case reflect.Struct, reflect.Map:
		upd := setUpdateD(nil, rv)
		res := make(bson.M, len(upd))
		for _, e := range upd {
			res[e.Key] = e.Value
		}
		return res
	default:
		panic(fmt.Errorf("unsupported argument type: %T", v))
	}

}

// Accepts a struct or map and returns a bson.D for mongo update $set operation.
// The struct or map can contain pointers, interfaces, maps, and structs.
// Map keys should be strings. BSON tags with omitempty are supported.
func SetUpdateD(v any) bson.D {
	rv := reflect.ValueOf(v)
	dereference(rv)

	switch rv.Kind() {
	case reflect.Struct, reflect.Map:
		return setUpdateD(nil, rv)
	default:
		panic(fmt.Errorf("unsupported argument type: %T", v))
	}
}

// Accepts a struct or map and returns a bson.D that resembles the original
// object structure. Map keys should be strings. The struct or map can contain
// pointers, interfaces, maps, and structs. BSON tags with omitempty are
// supported.
func Dump(v any) bson.D {
	rv := reflect.ValueOf(v)
	dereference(rv)

	switch rv.Kind() {
	case reflect.Struct, reflect.Map:
		return dump(rv)
	default:
		panic(fmt.Errorf("unsupported argument type: %T", v))
	}
}

func setUpdateD(path []string, v reflect.Value) bson.D {

	var (
		res = bson.D{}
		add func(string, reflect.Value)
		p   = func(k string) string { return strings.Join(append(path, k), ".") }
	)

	add = func(k string, v reflect.Value) {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if v.IsNil() {
				res = append(res, bson.E{Key: p(k), Value: v.Interface()})
			} else {
				add(k, v.Elem())
			}

		case reflect.Map, reflect.Struct:
			res = append(res, setUpdateD(append(path, k), v)...)

		default:
			res = append(res, bson.E{Key: p(k), Value: v.Interface()})
		}
	}

	walk(add, v)

	return res
}

func dump(v reflect.Value) bson.D {
	var (
		res = bson.D{}
		add func(string, reflect.Value)
	)

	add = func(k string, v reflect.Value) {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if v.IsNil() {
				res = append(res, bson.E{Key: k, Value: v.Interface()})
			} else {
				add(k, v.Elem())
			}

		case reflect.Map, reflect.Struct:
			res = append(res, bson.E{Key: k, Value: dump(v)})

		default:
			res = append(res, bson.E{Key: k, Value: v.Interface()})
		}
	}

	walk(add, v)

	return res
}

func walk(add func(k string, v reflect.Value), v reflect.Value) {
	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			return
		}

		it := v.MapRange()
		for it.Next() {
			k := it.Key()
			if k.Kind() != reflect.String {
				panic(fmt.Errorf("map key type should be string: %T", k))
			}
			add(k.String(), it.Value())
		}

	case reflect.Struct:
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			var n string
			f := v.Field(i)
			ft := typ.Field(i)
			if tag, ok := ft.Tag.Lookup("bson"); ok {
				if tag == "-" {
					continue
				}

				s := strings.Split(tag, ",")
				if slices.Contains(s[1:], "omitempty") && f.IsZero() {
					continue
				}

				n = s[0]
			} else {
				n = ft.Name
			}

			add(n, f)
		}
	}
}

func dereference(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			panic(fmt.Errorf("nil pointer or interface"))
		}
		v = v.Elem()
	}
	return v
}
