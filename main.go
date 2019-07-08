package main

import (
	"EchoDNS/protocol"
	"EchoDNS/server"
	"fmt"
	"net"
)

func main() {
	usage()
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:53")
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()

	ch := make(chan UDPPacket)
	go func() {
		for {
			data := make([]byte, 512)
			size, addr, _ := listener.ReadFromUDP(data)
			var message protocol.Message
			message.Decode(data[:size], 0)
			fmt.Println("Query:", message)
			ch <- UDPPacket{addr, message}
		}
	}()

	for {
		packet := <-ch
		go func() {
			response := server.Resolve(packet.message)
			_, err := listener.WriteToUDP(response, packet.addr)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

type UDPPacket struct {
	addr    *net.UDPAddr
	message protocol.Message
}

func usage() {
	fmt.Println("EchoDNS - A lightweight DNS server written in Go")
}
