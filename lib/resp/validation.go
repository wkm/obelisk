package resp

import (
	"fmt"
	"reflect"
)

// Possible validation errors
var (
	ErrInvalidOutputSize     = fmt.Errorf("Methods with more than two output values are not supported")
	ErrUnsupportedInputKind  = fmt.Errorf("Method has an unsupported input kind or type")
	ErrUnsupportedOutputKind = fmt.Errorf("Method has an unsupported output kind or type")
	ErrNonErrorSecondOutput  = fmt.Errorf("Two output methods are only supported with a second output of error")
)

// SupportedOutputKinds is a "hashset" of simple kinds which RESP can format for
// output. RESP also supports slices/arrays of these kinds as well as interfaces
// which implement an Error()string method.
var SupportedOutputKinds = map[reflect.Kind]struct{}{
	reflect.Uint64: {},
	reflect.Uint32: {},
	reflect.Uint16: {},
	reflect.Uint8:  {},
	reflect.Uint:   {},

	reflect.Int64: {},
	reflect.Int32: {},
	reflect.Int16: {},
	reflect.Int8:  {},
	reflect.Int:   {},

	reflect.String: {},
}

// SupportedInputKinds is a "hashset" of simple kinds which RESP can parse for
// input. RESP also supports slices/arrays of these kinds.
//
// The input language is significantly more limited to only support very general
// types since RESP is very loosely typed.
var SupportedInputKinds = map[reflect.Kind]struct{}{
	reflect.String: {},
	reflect.Array:  {},
}

// Ensure the specified method takes and returns types compatible with RESP
func validateMethod(method *reflect.Method) (err error) {
	if method.Type.NumIn() < 1 {
		// Ensure the method has a receiver. Since this is an internal method
		// within RESP it should never be triggered.
		panic("Can only validate methods with a reciever")
	}

	if method.Type.NumOut() > 2 {
		return ErrInvalidOutputSize
	}

	for i := 1; i < method.Type.NumIn(); i++ {
		mt := method.Type.In(i)
		println("in: ", i, "kind: ", mt)
		_, ok := SupportedInputKinds[mt.Kind()]
		if !ok {
			return ErrUnsupportedInputKind
		}
	}

	for i := 0; i < method.Type.NumOut(); i++ {
		mt := method.Type.Out(i)
		_, ok := SupportedOutputKinds[mt.Kind()]
		if !ok {
			return ErrUnsupportedOutputKind
		}
	}

	// ...
	return
}
