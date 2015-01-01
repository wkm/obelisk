package resp

import (
	"errors"
	"io"
	"reflect"
)

var (
	// ErrInvalidSize is returned if the given size prefix is not a positive integer.
	ErrInvalidSize = errors.New("invalid size")

	// ErrInvalidPrefix is returned if the argument prefix is not of a known type.
	ErrInvalidPrefix = errors.New("invalid prefix")

	// ErrInvalidType is returned if an invalid type is specified in the request.
	ErrInvalidType = errors.New("invalid type")
)

// Parse a RESP request into its constituent types. Must begin with a '*'
func readRequest(w io.ByteReader) (values []reflect.Value, err error) {
	b, err := w.ReadByte()
	if err != nil {
		return nil, err
	}

	if b != '*' {
		return nil, ErrInvalidType
	}

	sz, err := readSize(w)
	if err != nil {
		return
	}

	values = make([]reflect.Value, sz)
	for i := 0; i < sz; i++ {
		n, err := w.ReadByte()
		if err != nil {
			return nil, err
		}

		switch n {
		// ... handle other types
		case '$':
			nsz, err := readSize(w)
			if err != nil {
				return nil, err
			}

			str, err := readString(w, nsz)
			if err != nil {
				return nil, err
			}

			values[i] = reflect.ValueOf(str)

		default:
			return nil, ErrInvalidPrefix
		}
	}

	return
}

// Very fast parsing of non-negative integers. Copied from the Redis RESP implementation.
func readSize(w io.ByteReader) (sz int, err error) {
	for {
		n, err := w.ReadByte()

		if err != nil {
			return sz, err
		}

		if n == '\r' {
			break
		}

		sz *= 10
		sz += int(n - '0')
	}

	if sz < 0 {
		err = ErrInvalidSize
		return
	}

	consumeSep(w)
	return
}

// Consume a string of the given size from the reader into a newly allocated array.
func readString(w io.ByteReader, sz int) (str string, err error) {
	bb := make([]byte, sz)
	for i := 0; i < sz; i++ {
		bb[i], err = w.ReadByte()
		if err != nil {
			return
		}
	}
	str = string(bb)
	consumeSep(w)
	return
}

// Consume the '\r\n' separator
func consumeSep(w io.ByteReader) (err error) {
	n, err := w.ReadByte()
	if err != nil {
		return
	}

	if n == '\r' {
		_, err = w.ReadByte()
		if err != nil {
			return
		}
	}

	return
}
