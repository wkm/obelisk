package server

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/wkm/obelisk/lib/resp"
)

func TestServer(t *testing.T) {
	t.Parallel()

	closedConnection := false

	tempdir, err := ioutil.TempDir(os.TempDir(), "obelisk")
	if err != nil {
		t.Fatal(err)
	}

	s := new(App)
	s.Config.StoreDir = tempdir
	s.Start()

	// Start listening
	ln, err := net.Listen("tcp", "127.0.0.1:6666")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer ln.Close()

	go func() {
		t.Logf("Listening for connection")

		// Accept a connection
		conn, err := ln.Accept()
		if closedConnection {
			return
		}

		if err != nil {
			t.Fatal(err.Error())
		}

		t.Logf("Accepted connection: %v", conn)
		resp.Listen(s, conn, conn)
	}()

	// Open a connection to the server
	t.Logf("Dialing server")
	conn, err := net.Dial("tcp", ":6666")
	if err != nil {
		t.Fatalf("Couldn't connect to server: %v", err)
	}
	defer conn.Close()

	t.Logf("Connection opened")
	cw := bufio.NewWriter(conn)
	cr := bufio.NewReader(conn)

	// Test commands
	testcases := []struct {
		in, out string
	}{
		{"kvset a val", "+OK"},
		{"kvget a", "$3\nval"},
		{"kvset a", "+OK"},
		{"kvget a", "$0\n"},

		{"declare a712371 host/h1/service/s1•get service/s1/host/h1•get", "+OK"},
		{"declare a712372 host/h1/service/s1•set service/s1/host/h1•set", "+OK"},
		{"schema a712371 counter op 'number of get commands'", "+OK"},
		{"schema a712372 counter op 'number of set commands'", "+OK"},
		{"record a712371 2014-05-12 19", "+OK"},
		{"record a712371 2014-05-13 19", "+OK"},

		{
			"tagchildren host/h1/service",
			`*2
$32
host/h1/service/s1•get/a712371
$32
host/h1/service/s1•set/a712372`,
		},
	}

	for _, tc := range testcases {
		t.Logf("in > %#v", tc.in)
		cw.WriteString(tc.in)
		cw.WriteString("\n")
		cw.Flush()

		respOut := strings.Replace(tc.out+"\n", "\n", "\r\n", -1)

		// Consume lines over the connection
		t.Logf("out= %#v", respOut)
		lineCount := strings.Count(respOut, "\r\n")
		var bb bytes.Buffer
		for i := 0; i < lineCount; i++ {
			out, err := cr.ReadString('\n')
			if err != nil {
				t.Errorf("Unexpected error=%v", err)
			}
			bb.WriteString(out)
		}

		if bb.String() != respOut {
			t.Errorf("expected out=%#v; but got out=%#v and err=%#v", respOut, bb.String(), err)
		}
	}

	closedConnection = true
}
