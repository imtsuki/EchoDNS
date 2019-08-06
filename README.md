# EchoDNS

A lightweight DNS server written in Go.

## `package server`

实现了服务器的基本运行逻辑。

Go 语言最大的特性便是由其 `go` 关键字所指明的**协程**——在用户态进行协作式调度的轻量级线程。协程间通过 `channel` 提供的消息队列进行通信与同步，当遇到阻塞调用时，该协程便会**主动放弃** CPU，让出 CPU 供其他协程使用。由于协程十分轻量级，调度时性能损耗大大少于系统级线程，因此可十分简单且优雅地实现高并发。

`Serve()`	首先启动一个协程，用来侦听传入的 DNS 报文，并将其推入消息队列：

```go
	ch := make(chan UDPPacket)
	go func() {
		for {
			data := make([]byte, 512)
			size, addr, _ := listener.ReadFromUDP(data)
			var message Message
			message.Decode(data[:size], 0)
			ch <- UDPPacket{addr, message}
		}
	}()
```

由于主函数本身也是一个协程，因此当其从消息队列 `ch` 中尝试接收报文时，主协程便会被阻塞。

当接收到一个报文，协程恢复运行时，便立即启动一个新的协程来处理这个请求，由此实现服务器的并发。

```go
	for {
		packet := <-ch // 被阻塞
		// 启动新协程
		go func() {
			response := Resolve(packet.message)
			_, _ := listener.WriteToUDP(response, packet.addr)
		}()
	}
```

`Resolve` 函数实现对 DNS 查询报文的解析。其首先判断域名是否在程序给定的对照表中；若是，如为 `0.0.0.0` 则构造 `NXDOMAIN` 并返回，如为普通 IP 地址则构造 `ANSWER` 字段返回。

若不在对照表中，则调用 `FowardQuery()`，将报文的原始字节报文 `RawPacket` 转发至远端服务器。这里通过设置 `SetDeadline` 来对超时作出限制。由于每次转发时分配的 UDP 端口号均不一样，因此不会出现 ID 字段冲突的问题。

验收过程中，还额外增加了非法访问日志记录的功能。

## `package protocol`

所有数据结构均实现了 `Encodable` 及 `Decodable` 接口。`Encode()` 对数据结构进行序列化，构造 DNS 报文；`Decode()` 函数反序列化 DNS 报文，得到相应数据结构。

`Encodable` 接口定义了将数据结构序列化为字节序列的语义。该接口的返回值为一个字节数组的 `slice` 引用。实现了该接口的类型可简单通过调用其 `Encode()` 方法来获得对应的字节序列。例如，代表 DNS Header 字段的类型 `Header`，其实现 `Encodable` 接口如下：

```go
func (header Header) Encode() []byte {
	data := make([]byte, 12)
	binary.BigEndian.PutUint16(data[0:2], header.ID)
	data[2] |= uint8(header.MessageType) << 7
	if header.Authoritative {
		data[2] |= 1 << 2
	}
	if header.Truncation {
		data[2] |= 1 << 1
	}
	if header.RecursionDesired {
		data[2] |= 1 << 0
	}
	if header.RecursionAvailable {
		data[3] |= 1 << 7
	}
	data[2] |= (uint8(header.OpCode) & 0xF) << 3
	data[3] |= uint8(header.ResponseCode) & 0xF
	binary.BigEndian.PutUint16(data[4:6], header.QuestionCount)
	binary.BigEndian.PutUint16(data[6:8], header.AnswerCount)
	binary.BigEndian.PutUint16(data[8:10], header.NameServerCount)
	binary.BigEndian.PutUint16(data[10:12], header.AdditionalCount)
	return data
}
```

在序列化的过程中，需要注意**字节序**及**比特序**的问题。**字节序**定义了占据多个字节的数值类型，其内部的字节顺序；**比特序**则是在一个字节内部，各比特的前后顺序。一般 x86 平台上均默认为小端序，而**网络序**则均为大端序。

在网络编程的时候，字节序需要我们**手动转换**，而比特序则是由网络硬件**自动转换**的。

例如，在 RFC1035 中，DNS Header 的部分定义如下：

```
                                     1  1  1  1  1  1
       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                      ID                       |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

其中，ID 字段长度两个字节：

- 对于**字节序**，我们需要手动调整字节序，即将两个字节交换。
- 对于**比特序**，我们并未作出改动，在内存中的表示仍为小端序；而当发送到网络上时，硬件自动将小端序转换为大端序，正常发送。

而对于下方的标志位，则又是另一种情况：

- 对于**字节序**，没有占用多个字节的字段，不需要关注字节序。
- 对于**比特序**，我们以 QR 标志位作为示例。该字段定义在大端字节中的第 0 位。如果我们在编程时将其赋值到小端序的第 0 位，则在发送到网络后，其会被自动转换到第 7 位，这是错误的；所以此时我们需要手动先将比特序反转，将 QR 赋到第 7 位，这样在网络硬件自动转换后，就会回到第 0 位的正确位置。

`Decode()` 函数签名为 `Decode(data []byte, off int)` 。注意到 DNS 报文中包含不定常字段及域名压缩存储等特性，使得偏移参数成为必要。每个 `Decode()` 均会返回新的偏移量。

其中，最富有趣味的是 `NAME` 字段的解析。高两位的 00 指明了下一个串的长度；11 则指明了两个字节共同构成一个指针，实现域名的压缩存储。本程序可实现压缩存储的完全解析；但囿于时间紧迫，仅实现了域名与 QUESTION 字段中相同的压缩存储的序列化。

# 总结

由于各字段在 RFC 中均有明确定义，且采用了测试驱动的方法进行开发，因此编写过程中较为简单。过程为先实现转发功能，使得 DNS 服务器一开始就可以运行；然后逐渐加入、替换掉字段的解析、报文的构造等，保证了开发的良好体验。唯一遇到的 BUG 是解析域名指针时，偏移量设置错误的问题，不过马上就发觉到了。

DNS 协议的设计在妥协于众多条件的背景下，以其简单与可拓展，经久不衰，支撑起了整个互联网。某种意义上，这正是 KISS 原则的一种体现：实现功能的最小子集，而由 Resource Record 实现高拓展性，使得 DNS 不会因早期限制而在现代失去可用性。
