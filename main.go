package main

import (
	"fmt"
	"net"
	"time"
	"EchoDNS/protocol"
)

type Server struct {
	conn *net.UDPConn
}

func (server Server) run() {

}

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
			var header protocol.Header
			
			fmt.Println("Received Request:", header.Decode(data))
			ch <- UDPPacket{addr, data[:size]}
		}
	}()

	for {
		packet := <-ch
		go func() {
			addr, _ := net.ResolveUDPAddr("udp", "114.114.114.114:53")
			socket, _ := net.DialUDP("udp", nil, addr)
			_ = socket.SetDeadline(time.Now().Add(time.Duration(time.Second * 2)))
			_, _ = socket.Write(packet.data)
			result := make([]byte, 512)
			size, addr, _ := socket.ReadFromUDP(result)
			listener.WriteToUDP(result[:size], packet.addr)
			fmt.Println("Response to:", packet.addr)
			socket.Close()
		}()
	}
}

type UDPPacket struct {
	addr *net.UDPAddr
	data []byte
}

func usage() {
	fmt.Println("EchoDNS - A lightweight DNS server written in Go")
}
