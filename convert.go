package csvdecoder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// convertAssignValues copies to dest the value in src, converting it if possible.
// An error is returned if the conversion is not possible.
// dest is expected to be a non-nil pointer type.
func convertAssignValue(dest interface{}, src string) error {
	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errNotPtr
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if src == "" {
		// zero values should leave the dest untouched
		return nil
	}

	// check if the destination implements the Decoder interface
	if decoder, ok := dest.(Interface); ok {
		return decoder.DecodeField(src)
	}

	var sv reflect.Value

	// simple cases without reflect
	switch d := dest.(type) {
	case *string:
		*d = src
		return nil
	case *[]byte:
		*d = []byte(src)
		return nil
	case *bool:
		bv, err := strconv.ParseBool(src)
		if err == nil {
			*d = bv
		}
		return err
	case *interface{}:
		*d = sv.Interface()
		return nil
	}

	// cases with reflect
	sv = reflect.ValueOf(src)
	dv := reflect.Indirect(dpv)

	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		dv.Set(sv)
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.Ptr:
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssignValue(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := strconv.ParseInt(src, 10, dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64, err := strconv.ParseUint(src, 10, dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		f64, err := strconv.ParseFloat(src, dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetFloat(f64)
		return nil
	case reflect.Slice, reflect.Array:
		objType := dv.Type()
		obj := reflect.New(objType).Interface()

		if err := json.NewDecoder(strings.NewReader(src)).Decode(&obj); err != nil {
			return fmt.Errorf("could not parse %s as JSON array: %w", src, err)
		}

		dv.Set(reflect.ValueOf(obj).Elem())
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing type %T into type %T", src, dest)
}
