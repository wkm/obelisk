package resp

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Naive implementation of a type parser.
func parse(t reflect.Kind, line string) (remaining string, value reflect.Value, err error) {
	switch t {
	case reflect.Int:
		return parseInt(line)

	case reflect.Float64:
		return parseFloat(line)

	case reflect.String:
		return parseString(line)

	case reflect.Slice:
		return parseSlice(line)

	default:
		err = errors.New(fmt.Sprintf("Unsupported kind: %v", t))
		return
	}
}

func nextToken(line string) (token, rem string, err error) {
	// consume leading whitespace
	start := 0
	for ; start < len(line); start++ {
		switch line[start] {
		case ' ':
			continue
		}
		break
	}

	// consume non-whitespace
	stop := start
Whitespace:
	for ; stop < len(line); stop++ {
		switch line[stop] {
		case ' ', '\n':
			break Whitespace
		}
	}

	return line[start:stop], line[stop:], nil
}

func parseInt(line string) (rem string, val reflect.Value, err error) {
	var intval int64
	var tok string

	tok, rem, err = nextToken(line)
	intval, err = strconv.ParseInt(tok, 10, 32)

	if err != nil {
		return
	}

	val = reflect.ValueOf(int(intval))
	return
}

func parseFloat(line string) (rem string, val reflect.Value, err error) {
	var floatval float64
	var tok string

	tok, rem, err = nextToken(line)
	floatval, err = strconv.ParseFloat(tok, 32)

	if err != nil {
		return
	}

	val = reflect.ValueOf(floatval)
	return
}

func parseString(line string) (rem string, val reflect.Value, err error) {
	var tok string
	tok, rem, err = nextToken(line)
	if err != nil {
		return
	}

	val = reflect.ValueOf(tok)
	return
}

func parseSlice(line string) (rem string, val reflect.Value, err error) {
	strType := reflect.TypeOf("")
	slice := reflect.MakeSlice(reflect.SliceOf(strType), 0, 0)
	rem = line

	for {
		prev := rem
		rem, val, err = parseString(rem)
		if rem == prev {
			break
		}

		slice = reflect.Append(slice, val)
	}

	return rem, slice, err
}
