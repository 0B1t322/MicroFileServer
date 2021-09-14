package typechecker_test

import (
	"reflect"
	"github.com/MicroFileServer/pkg/repositories/typechecker"
	"testing"
)

type someType struct {

}

func TestFunc_NewTypeChecker(t *testing.T) {
	f := typechecker.NewSingle(reflect.TypeOf(someType{}))

	if err := f(&[]someType{}); err != nil {
		t.Log(err)
	}
}
