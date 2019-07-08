package main

import (
	"EchoDNS/server"
	"fmt"
)

func main() {
	usage()
	server.Serve()
}

func usage() {
	fmt.Println("EchoDNS - A lightweight DNS server written in Go")
}
