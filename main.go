package main

import (
	"EchoDNS/server"
	"flag"
	"fmt"
)

func main() {
	var remote string
	var hosts string
	var debug bool
	usage()
	flag.StringVar(&remote, "r", "114.114.114.114", "remote DNS server address")
	flag.StringVar(&hosts, "h", "hosts", "hosts file")
	flag.BoolVar(&debug, "d", false, "print debug info")
	flag.Usage()
	flag.Parse()

	server.Serve(remote, hosts, debug)
}

func usage() {
	fmt.Println("EchoDNS - A lightweight DNS server written in Go")
}
