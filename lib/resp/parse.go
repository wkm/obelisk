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

	default:
		err = fmt.Eprintf("unsupported kind: %v", t)
		return
	}
}

func nextToken(line string) (token, rem string, err error) {
	// Consume leading whitespace
	start := 0
	for ; start < len(line); start++ {
		switch line[start] {
		case ' ':
			continue
		}
		break
	}

	// Consume as much non-whitespace as possible
	stop := start
	quoted := false

	// Test for quoted string
	if stop < len(line) && (line[stop] == '"' || line[stop] == '`') {
		quoted = true
		stop++
	}

	escaping := false
	escapedChars := 0
Whitespace:
	for ; stop < len(line); stop++ {
		// Only a single character can be escaped
		escaped := escaping
		escaping = false

		switch line[stop] {
		case '\\':
			escapedChars++
			if !escaping {
				escaping = true
			}

		case ' ', '\n':
			if !quoted {
				break Whitespace
			}

		case '"', '`':
			if !escaped {
				break Whitespace
			}
		}
	}

	if !quoted {
		return line[start:stop], line[stop:], nil
	} else if escapedChars == 0 {
		// Don't copy the string, but skip the opening and closing quotes
		return line[start+1 : stop], line[stop+1:], nil
	} else {
		newline, err := strconv.Unquote(line[start : stop+1])
		return string(newline), line[stop+1:], err
	}
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
