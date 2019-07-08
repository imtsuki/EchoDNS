package main

import (
	"EchoDNS/protocol"
	"fmt"
	"net"
	"time"
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
			var message protocol.Message
			message.Decode(data[:size], 0)
			fmt.Println("Query:", message)
			ch <- UDPPacket{addr, message}
		}
	}()

	for {
		packet := <-ch
		go func() {
			if len(packet.message.Questions) > 0 {
				if packet.message.Questions[0].Name.Domain == "bupt.edu.cn." {
					response := packet.message
					response.Header.MessageType = protocol.Response
					response.Header.RecursionDesired = true
					response.Header.RecursionAvailable = true
					if packet.message.Questions[0].Type == protocol.TypeA {
						answer := protocol.Resource{
							Name: protocol.Name{
								Compressed: true,
							},
							Type:   protocol.TypeA,
							Class:  protocol.ClassINET,
							TTL:    53,
							Length: 4,
							Data: &protocol.AResource{
								IP: [4]byte{10, 3, 8, 216},
							},
						}
						response.Answers = append(response.Answers, answer)
					}
					listener.WriteToUDP(response.Encode(), packet.addr)
					return
				}
			}
			addr, _ := net.ResolveUDPAddr("udp", "114.114.114.114:53")
			socket, _ := net.DialUDP("udp", nil, addr)
			_ = socket.SetDeadline(time.Now().Add(time.Duration(time.Second * 2)))
			_, _ = socket.Write(packet.message.RawPacket)
			result := make([]byte, 512)
			size, addr, _ := socket.ReadFromUDP(result)
			listener.WriteToUDP(result[:size], packet.addr)
			fmt.Println("Response to:", packet.addr)
			socket.Close()
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
