package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func NewQueryStringDecoder() Decoder {
	return &queryStringDecoder{}
}

type queryStringDecoder struct {
}

func (p *queryStringDecoder) BodyDecoder() BodyDecoder {
	return nil
}

func (p *queryStringDecoder) Decode(_ map[string]string, request *http.Request, target interface{}) error {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("invalid type in pointer: %s", t.Elem().Kind().String())
		}
		return p.decodeStruct(request, t.Elem(), v.Elem())
	default:
		return fmt.Errorf("invalid type for decoding: %s", t.Kind().String())
	}
}

func (p *queryStringDecoder) decodeStruct(request *http.Request, t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldV := v.Field(i)
		if routingName, ok := f.Tag.Lookup("query"); ok {
			if val := request.URL.Query().Get(routingName); val != "" {
				if !fieldV.CanSet() {
					return fmt.Errorf("cannot set field %s", f.Name)
				}
				switch f.Type.Kind() {
				case reflect.String:
					fieldV.SetString(val)
				case reflect.Bool:
					switch strings.ToLower(val) {
					case "1":
						fallthrough
					case "true":
						fallthrough
					case "yes":
						fieldV.SetBool(true)
					case "0":
						fallthrough
					case "false":
						fallthrough
					case "no":
						fieldV.SetBool(false)
					default:
						return fmt.Errorf("invalid value for field %s: %s", f.Name, val)
					}
				case reflect.Int:
					fallthrough
				case reflect.Int8:
					fallthrough
				case reflect.Int16:
					fallthrough
				case reflect.Int32:
					fallthrough
				case reflect.Int64:
					intVal, err := strconv.ParseInt(val, 10, 64)
					if err != nil {
						return fmt.Errorf("failed to parse int value for field %s: %s", f.Name, val)
					}
					fieldV.SetInt(intVal)
				case reflect.Uint:
					fallthrough
				case reflect.Uint8:
					fallthrough
				case reflect.Uint16:
					fallthrough
				case reflect.Uint32:
					fallthrough
				case reflect.Uint64:
					intVal, err := strconv.ParseUint(val, 10, 64)
					if err != nil {
						return fmt.Errorf("failed to parse uint value for field %s: %s", f.Name, val)
					}
					fieldV.SetUint(intVal)
				default:
					return fmt.Errorf("unsupported field type for field %s: %s", f.Name, f.Type.Name())
				}

			}
		}
		if f.Type.Kind() == reflect.Struct {
			if err := p.decodeStruct(request, f.Type, fieldV); err != nil {
				return fmt.Errorf("failed to decode field %s (%w)", f.Name, err)
			}
		}
	}
	return nil
}
