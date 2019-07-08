package server

import (
	. "EchoDNS/protocol"
	"fmt"
	"net"
	"time"
)

var domains = map[string]net.IP{}

func init() {
	domains["bupt.edu.cn."] = net.ParseIP("10.3.8.216")
	domains["baidu.com."] = net.ParseIP("0.0.0.0")
}

func Resolve(message Message) (responsePacket []byte) {
	if len(message.Questions) > 0 {
		ip, ok := domains[message.Questions[0].Name.Domain]
		if ok {
			response := message
			response.Header.MessageType = Response
			response.Header.RecursionDesired = true
			response.Header.RecursionAvailable = true
			if ip.Equal(net.IPv4zero) {
				response.Header.ResponseCode = NXDomain
				return response.Encode()
			}
			if message.Questions[0].Type == TypeA {
				answer := Resource{
					Name: Name{
						Compressed: true,
					},
					Type:   TypeA,
					Class:  ClassINET,
					TTL:    53,
					Length: 4,
					Data: &AResource{
						IP: ip,
					},
				}
				response.Answers = append(response.Answers, answer)
			}
			return response.Encode()
		}
	}

	responsePacket, err := ForwardQuery(message)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return responsePacket
}

func ForwardQuery(message Message) ([]byte, error) {
	addr, _ := net.ResolveUDPAddr("udp", "114.114.114.114:53")
	socket, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := socket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	err = socket.SetDeadline(time.Now().Add(time.Duration(time.Second * 2)))
	if err != nil {
		return nil, err
	}
	_, err = socket.Write(message.RawPacket)
	if err != nil {
		return nil, err
	}
	result := make([]byte, 512)
	size, addr, err := socket.ReadFromUDP(result)
	if err != nil {
		return nil, err
	}
	return result[:size], nil

}
