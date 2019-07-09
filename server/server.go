package server

import (
	. "EchoDNS/protocol"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var domains = map[string]net.IP{}
var remoteServer string

type UDPPacket struct {
	addr    *net.UDPAddr
	message Message
}

func init() {
	domains["bupt.edu.cn."] = net.ParseIP("10.3.8.216")
	domains["baidu.com."] = net.ParseIP("0.0.0.0")
}

func readHosts(hosts string) {
	wd, _ := os.UserHomeDir()
	_ = os.Chdir(wd)
	file, err := os.Open(hosts)
	if err != nil {
		fmt.Printf("Load hosts file %s error\n", hosts)
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		host := strings.Fields(line)
		if len(host) == 2 {
			if !strings.HasSuffix(host[1], ".") {
				host[1] = host[1] + "."
			}
			domains[host[1]] = net.ParseIP(host[0])
		}
	}
}

func Serve(remote string, hosts string, debug bool) {
	remoteServer = remote + ":53"
	readHosts(hosts)
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
			var message Message
			message.Decode(data[:size], 0)
			fmt.Println("Query:", message)
			ch <- UDPPacket{addr, message}
		}
	}()

	for {
		packet := <-ch
		go func() {
			response := Resolve(packet.message)
			_, err := listener.WriteToUDP(response, packet.addr)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
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
	addr, _ := net.ResolveUDPAddr("udp", remoteServer)
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
