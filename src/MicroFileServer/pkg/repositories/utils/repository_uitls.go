package utils

import (
	"context"
	"fmt"
	"reflect"
)


// CreateCheckTypeFunc
//
// Create func that check type of value (type, pointer og this type, slice and slice of pointers)
//
// If type is correct execute given func
//
// If it slice execute for all values this func
func CreateCheckTypeFunc(
	t reflect.Type,
	fun func(context.Context, interface{}) error,
) func(context.Context, interface{}) error {
	return func(ctx context.Context, v interface{}) error {
		typeOfV := reflect.TypeOf(v)
		tPtr := reflect.PtrTo(t)

		if typeOfV.AssignableTo(t) {
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(tPtr) {
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(reflect.SliceOf(t)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Addr().Interface()
				if err := fun(ctx, value); err != nil {
					return err
				}
			}
		} else if typeOfV.AssignableTo(reflect.SliceOf(tPtr)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Interface()
				if err := fun(ctx, value); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf(
				"Unexpected Type %T, Expected %s or %s or %s or %s",
				v, t, reflect.PtrTo(t), reflect.SliceOf(t), reflect.SliceOf(tPtr),
			)
		}

		return nil
	}
}

func GetStructType(Type interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(Type)
	switch t.Kind() {
	case reflect.Struct:
		break
	default:
		return nil, fmt.Errorf("You give %s of %s expect %s", t.Kind(), t, reflect.Struct)
	}

	return t, nil
}