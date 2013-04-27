package timestore

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// load up values from a dump
func (s *Store) Load(r io.Reader) error {
	// we don't need to lock because our individual inserts will
	statLoad.Incr()
	buf := bufio.NewReader(r)
	for {
		str, err := buf.ReadString('\n')
		switch err {
		case io.EOF:
			break
		default:
			return err
		}

		line := strings.SplitN(str, ":", 2)
		if len(line) != 2 {
			return errors.New("invalid syntax in buffer: " + str)
		}

		key, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			return err
		}

		components := strings.Split(line[1], ",")
		for _, pair := range components {
			if pair == "" {
				// this is the final part of components
				break
			}

			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				return errors.New("invalid syntax in pair: " + pair)
			}

			time, err := strconv.ParseUint(parts[0], 10, 64)
			if err != nil {
				return err
			}

			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}

			// time to insert this badboy
			s.Insert(key, time, value)
		}
	}

	return nil
}
