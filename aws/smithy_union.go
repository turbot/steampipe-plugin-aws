package aws

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// AWS SDK v2 models smithy unions as an interface with one concrete member
// type per variant, named <Union>Member<Variant> and holding a single Value
// field. Marshalled directly they serialise as {"Value": ...}, which hides the
// variant name from SQL. This file flattens them to {"<Variant>": ...}.

var smithyUnionMemberRe = regexp.MustCompile(`^[A-Z]\w*Member([A-Z]\w*)$`)

var timeType = reflect.TypeOf(time.Time{})

// flattenSmithyUnionsTransform is a transform step that applies
// flattenSmithyUnions to the current value.
func flattenSmithyUnionsTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return flattenSmithyUnions(d.Value), nil
}

// flattenSmithyUnions recursively converts a value into plain maps and slices,
// replacing every smithy union member with a single-key map named after the
// variant. Non-union structs keep their exported field names. A nil input, nil
// pointer or nil interface returns nil.
func flattenSmithyUnions(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return flattenReflectValue(reflect.ValueOf(v))
}

func flattenReflectValue(rv reflect.Value) interface{} {
	if !rv.IsValid() {
		return nil
	}

	switch rv.Kind() {
	case reflect.Interface, reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return flattenReflectValue(rv.Elem())

	case reflect.Struct:
		if variant, ok := smithyUnionVariant(rv); ok {
			return map[string]interface{}{variant: flattenReflectValue(rv.FieldByName("Value"))}
		}
		if rv.Type() == timeType {
			return rv.Interface()
		}
		out := map[string]interface{}{}
		t := rv.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			out[f.Name] = flattenReflectValue(rv.Field(i))
		}
		return out

	case reflect.Slice:
		if rv.IsNil() {
			return nil
		}
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return rv.Interface()
		}
		out := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = flattenReflectValue(rv.Index(i))
		}
		return out

	case reflect.Map:
		if rv.IsNil() {
			return nil
		}
		out := map[string]interface{}{}
		iter := rv.MapRange()
		for iter.Next() {
			out[fmt.Sprint(iter.Key().Interface())] = flattenReflectValue(iter.Value())
		}
		return out

	default:
		return rv.Interface()
	}
}

// smithyUnionVariant reports whether rv is a smithy union member struct and,
// if so, returns the variant name. A member is recognised by its type name
// (<Union>Member<Variant>) and by having Value as its only exported field.
func smithyUnionVariant(rv reflect.Value) (string, bool) {
	t := rv.Type()
	m := smithyUnionMemberRe.FindStringSubmatch(t.Name())
	if m == nil {
		return "", false
	}
	exported := 0
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			exported++
		}
	}
	if exported != 1 {
		return "", false
	}
	if _, ok := t.FieldByName("Value"); !ok {
		return "", false
	}
	return m[1], true
}
