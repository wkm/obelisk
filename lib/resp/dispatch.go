package resp

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Dispatch provides case-insensitivity support
type dispatch struct {
	recv      interface{}
	recvValue reflect.Value
	commands  map[string]*reflect.Method
}

func normalizeMethodName(name string) string {
	return strings.ToLower(name)
}

// newDispatch creates a dispatcher against an interface
func newDispatch(receiver interface{}) (d *dispatch, err error) {
	d = new(dispatch)
	d.recv = receiver
	d.recvValue = reflect.ValueOf(d.recv)
	d.commands = make(map[string]*reflect.Method)

	rval := reflect.TypeOf(receiver)
	for i := 0; i < rval.NumMethod(); i++ {
		method := rval.Method(i)
		d.commands[normalizeMethodName(method.Name)] = &method

		err = validateMethod(&method)
		if err != nil {
			return
		}
	}

	return
}

// Ensure the specified method takes and returns types compatible with RESP
func validateMethod(method *reflect.Method) (err error) {
	// ...
	return
}

// Call parses and executes the function call against the receiver.
func (d *dispatch) Call(line string) (res string, err error) {
	var methodName string
	methodName, line, err = nextToken(line)

	method, ok := d.commands[normalizeMethodName(methodName)]
	if !ok {
		err = fmt.Errorf("Unknown command %q", methodName)
		return
	}

	fn := method.Func

	ins := make([]reflect.Value, method.Type.NumIn())
	for i := range ins {
		t := method.Type.In(i)

		if i == 0 {
			ins[i] = d.recvValue
			continue
		}

		line, ins[i], err = parse(t.Kind(), line)
		if err != nil {
			return
		}
	}

	var bb bytes.Buffer
	outs := fn.Call(ins)

	// Methods which return null
	if len(outs) == 0 {
		_, err = writeOk(&bb)
		res = bb.String()
		return
	}

	for i := range outs {
		o := outs[i]
		_, err = write(&bb, method.Type.Out(i).Kind(), o)
		if err != nil {
			return
		}
	}

	res = bb.String()
	return
}
