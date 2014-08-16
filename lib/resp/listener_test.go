package resp

import (
	"bufio"
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	k := newKeyval() // use keyval from dispatch_test.go

	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer ln.Close()

	// Start listenging
	go func() {
		t.Logf("Listening for connection")
		for {
			// Accept connections
			conn, err := ln.Accept()
			if err != nil {
				t.Fatal(err.Error())
			}

			t.Logf("Accepted a connected: %v", conn)

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

	testcases := []struct{ in, out string }{
		{"put k1 12\n", "+OK\r\n"},
		{"put k2 14\n", "+OK\r\n"},
	}
	for _, tc := range testcases {
		t.Logf("in> %s", tc.in)
		cw.WriteString(tc.in)
		cw.Flush()

		t.Logf("out> %s", tc.out)
		out, err := cr.ReadString('\n')
		if err != nil || out != tc.out {
			t.Errorf("expected out=%v; but got out=%v and err=%v", tc.out, out, err)
		}
	}
}
