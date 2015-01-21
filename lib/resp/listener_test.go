package resp

import (
	"bufio"
	"bytes"
	"net"
	"strings"
	"testing"
)

func TestWriteRequest(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		q   []interface{}
		out string
	}{
		{[]interface{}{"LLEN", "mylist"}, "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"},
	}

	for _, tc := range testcases {
		t.Logf("q:%v out:%#v", tc.q, tc.out)

		var bb bytes.Buffer
		nn, err := WriteRequest(&bb, tc.q...)
		out := bb.String()
		if err != nil {
			t.Errorf("unexpected error %#v", err)
		}

		if nn != len(tc.out) {
			t.Errorf("invalid number of bytes written")
		}

		if out != tc.out {
			t.Errorf("expected %#v; got %#v", tc.out, out)
		}
	}
}

// Start a simple test service to validate RESP validation
func TestListen(t *testing.T) {
	t.Parallel()

	k := newKeyval() // reuse keyval from dispatch_test.go

	ln, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer ln.Close()

	// Start listening
	go func() {
		t.Logf("Listening for connection")
		for {
			// Accept connections
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal(err.Error())
			}

			t.Logf("Accepted connection: %v", conn)
			Listen(k, conn, conn)
		}
	}()

	// Open a connection to the server
	t.Logf("Dialing server")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		t.Fatalf("Couldn't connect to server: %v", err)
	}

	t.Logf("Connection opened")
	cw := bufio.NewWriter(conn)
	cr := bufio.NewReader(conn)

	// Test textual commands
	testcases := []struct{ in, out string }{
		{"gerald", "-Error: Unknown command \"gerald\"\r\n"},
		{"put k1 12", "+OK\r\n"},
		{"put k2 14", "+OK\r\n"},
		{"get k1", "$2\r\n12\r\n"},
		{"get k2", "$2\r\n14\r\n"},
		{"get unknown", "$0\r\n"},
	}
	for _, tc := range testcases {
		t.Logf("in > %#v", tc.in)
		cw.WriteString(tc.in)
		cw.WriteString("\n")
		cw.Flush()

		// Consume lines over the connection
		t.Logf("out= %#v", tc.out)
		lineCount := strings.Count(tc.out, "\r\n")
		var bb bytes.Buffer
		for i := 0; i < lineCount; i++ {
			out, err := cr.ReadString('\n')
			if err != nil {
				t.Errorf("Unexpected error=%v", err)
			}
			bb.WriteString(out)
		}

		if bb.String() != tc.out {
			t.Errorf("expected out=%#v; but got out=%#v and err=%#v", tc.out, bb.String(), err)
		}
	}
}
