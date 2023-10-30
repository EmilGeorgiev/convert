package convert

import (
	"reflect"
	"time"
)

// SrcToDst convert src to dst. The function get values of the fields
// of type src and assign them to the fields of type dst. The method will
// assign values from src to dst only to this fields that have equal names.
// The dst must be a pointer value, otherwise it wil panic.
// For example:
//
//			src := struct{
//	     	ID int
//				Name string
//			}{
//	 		ID: 123,
//				Name: "Emil",
//			}
//
//			dst := struct{
//	     	Identifier int
//				Name string
//			}{}
//
//			SrcToDst(src, &dst)
//
// The result will be:
//
//	dst = struct{
//		ID int
//		Name string
//	}{0, "Emil"}
//
// Also the function can convert pointer fields to non pointer fields and vice versa.
func SrcToDst(src interface{}, dst interface{}) {
	srcValue := getNonPointerValue(reflect.ValueOf(src))
	dstValue := getNonPointerValue(reflect.ValueOf(dst))
	if dstValue.Kind() == reflect.Interface {
		dstValue = dstValue.Elem()
	}

	if dstValue.Kind() == reflect.Struct {
		setStructures(srcValue, dstValue)
		return
	}

	setPrimitives(srcValue, dstValue)
}

func setStructures(srcValue reflect.Value, dstValue reflect.Value) {
	nonPtrDstStruct := dstValue
	if dstValue.Kind() == reflect.Ptr {
		nonPtrDstStruct = dstValue.Elem()
		if nonPtrDstStruct.Kind() == reflect.Invalid {
			nonPtrDstStruct = reflect.New(dstValue.Type().Elem()).Elem()
		}
	}

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		if srcField.Kind() == reflect.Ptr && srcField.IsZero() {
			continue
		}

		if srcField.Kind() == reflect.Slice && srcField.IsZero() {
			continue
		}

		srcFieldName := srcValue.Type().Field(i).Name
		dstField := nonPtrDstStruct.FieldByName(srcFieldName)
		if !dstField.IsValid() {
			continue
		}

		setFields(srcField, dstField)
	}

	value := nonPtrDstStruct
	if dstValue.Kind() == reflect.Ptr {
		value = reflect.New(value.Type())
		value.Elem().Set(nonPtrDstStruct)
	}

	if dstValue.CanSet() {
		dstValue.Set(value)
	}

}

func setFields(srcField reflect.Value, dstField reflect.Value) {
	nonPtrSrcField := srcField
	if srcField.Kind() == reflect.Ptr {
		nonPtrSrcField = srcField.Elem()
		srcField = srcField.Elem()
	}

	nonPtrDstField := dstField.Type()
	if dstField.Kind() == reflect.Ptr {
		nonPtrDstField = nonPtrDstField.Elem()

		srcField = reflect.New(srcField.Type())
		srcField.Elem().Set(nonPtrSrcField)
	}

	if nonPtrSrcField.Kind() != nonPtrDstField.Kind() {
		return
	}

	if nonPtrSrcField.Kind() == reflect.Struct && nonPtrDstField.Kind() == reflect.Struct {
		fieldValue := nonPtrSrcField.Interface()
		if _, ok := fieldValue.(time.Time); !ok {
			setStructures(nonPtrSrcField, dstField)
			return
		}
	}

	if dstField.CanSet() {
		if !srcField.Type().ConvertibleTo(dstField.Type()) {
			return
		}

		srcField = srcField.Convert(dstField.Type())
		if dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}

func setPrimitives(src reflect.Value, dst reflect.Value) {
	if !src.Type().ConvertibleTo(dst.Type()) {
		return
	}

	src = src.Convert(dst.Type())
	if dst.CanSet() {
		dst.Set(src)
	}
}

func getNonPointerValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		v = getNonPointerValue(v)
	}

	return v
}
