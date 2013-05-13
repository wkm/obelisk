package main

import (
	_ "circuit/load/cmd"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"flag"
	"fmt"
	"obelisk/lib/rconfig"
	"os"
)

// utility for interacting with remote configurations
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		usageAndExit()
	}

	switch args[2] {
	case "get":
		// FIXME
		panic("not implemented")
		// get()

	case "set":
		// FIXME
		panic("not implemented")
		// set()

	case "getall":
		getall()

	default:
		usageAndExit()
	}
}

func usageAndExit() {
	fmt.Printf("Usage: %s AnchorPath getall|get key|set key value", flag.Args()[0])
	os.Exit(1)
}

func connect() circuit.X {
	file, err := anchorfs.OpenFile(flag.Args()[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)\n", err)
		os.Exit(1)
	}

	x, err := circuit.TryDial(file.Owner(), rconfig.ServiceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing service (%s)\n", err)
		os.Exit(1)
	}

	return x
}

func getall() {
	x := connect()
	retrn := x.Call("GetAll")

	values := retrn[0].(map[string]string)
	for key, value := range values {
		fmt.Printf("%s:%s\n", key, value)
	}
}
