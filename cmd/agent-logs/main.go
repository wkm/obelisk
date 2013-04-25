package main

import (
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"obelisk/rlog"
	"os"
)

// utility for getting the logs off a remote worker
func main() {
	if len(os.Args) < 2 {
		usageAndExit()
	}

	forwardLogs()
}

func usageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s AnchorPath\n", os.Args[0])
	os.Exit(1)
}

func connect() circuit.X {
	node := os.Args[1]
	file, err := anchorfs.OpenFile(node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening %s (%s)\n", node, err)
		os.Exit(1)
	}

	x, err := circuit.TryDial(file.Owner(), rlog.ServiceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing service (%s)\n", err)
		os.Exit(1)
	}

	if x == nil {
		fmt.Fprintf(os.Stderr, "Could not open cross pointer\n")
		os.Exit(1)
	}

	return x
}

func forwardLogs() {
	x := connect()
	retrn := x.Call("FlushLog")
	buffer := retrn[0].([]byte)
	os.Stdout.Write(buffer)
}
