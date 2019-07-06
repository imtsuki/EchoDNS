package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	usage()
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:53")
	listener, _ := net.ListenUDP("udp", addr)
	defer func() {
		_ = listener.Close()
	}()
	ch := make(chan UDPPacket)
	go func() {
		for {
			data := make([]byte, 512)
			len, addr, _ := listener.ReadFromUDP(data)
			fmt.Println("Received Request:", addr)
			ch <- UDPPacket{addr, data[:len]}
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
			len, addr, _ := socket.ReadFromUDP(result)
			listener.WriteToUDP(result[:len], packet.addr)
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
