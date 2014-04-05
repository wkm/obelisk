package resp

import (
	"log"
	"reflect"
	"strings"
)

// Dispatch provides case-insensitivity support
type dispatch struct {
	commands map[string]*reflect.Method
}

func normalizeMethodName(name string) string {
	return strings.ToLower(name)
}

func newDispatch(reciever interface{}) (d *dispatch, err error) {
	d = new(dispatch)
	d.commands = make(map[string]*reflect.Method)

	rval := reflect.TypeOf(reciever)
	log.Printf("rval: %v", rval)
	log.Printf(" -- # methods: %v", rval.NumMethod())
	for i := 0; i < rval.NumMethod(); i++ {
		method := rval.Method(i)
		d.commands[normalizeMethodName(method.Name)] = &method
		log.Printf("%s: %v", method.Name, method)

		// method.Func.Call(in)
	}

	return
}
