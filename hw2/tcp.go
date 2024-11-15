package main

import (
	"fmt"
	"strings"
	"time"
)

var (
	server = NewConnection("SERVER")
	client = NewConnection("CLIENT")
)

func main() {
	fmt.Println("Starting TCP connection simulation...")
	go serverListen()
	time.Sleep(time.Second)
	clientListen()
}

func serverListen() {
	server.Accept(*client)

	for {
		packet := server.Receive()
		fmt.Println("Received: ", string(packet.Data))
	}
}

func clientListen() {
	client.Connect(*server)
	client.Send([]byte("Hello, World!"))
	client.Close()
}

type Flags struct {
	ACK bool
	SYN bool
	RST bool
	FIN bool
}

type Packet struct {
	Source string
	Dest   string
	Seq    uint32
	Ack    uint32
	Flags  Flags
	Data   []byte
}

type Connection struct {
	Name    string
	Target  chan Packet
	Conn    chan Packet
	Packets []Packet
}

func (f *Flags) String() string {
	flags := make([]string, 0)

	if f.ACK {
		flags = append(flags, "ACK")
	}
	if f.SYN {
		flags = append(flags, "SYN")
	}
	if f.RST {
		flags = append(flags, "RST")
	}
	if f.FIN {
		flags = append(flags, "FIN")
	}

	return strings.Join(flags, ", ")
}

func NewConnection(name string) *Connection {
	return &Connection{
		Name:    name,
		Target:  make(chan Packet, 1),
		Conn:    make(chan Packet, 1),
		Packets: make([]Packet, 0),
	}
}

func (c *Connection) send(packet Packet) {
	// This is to simulate network delay
	time.Sleep(time.Second)

	fmt.Printf(
		"%s --> %s [%s] Seq=%d Ack=%d Len=%d\n",
		packet.Source,
		packet.Dest,
		packet.Flags.String(),
		packet.Seq,
		packet.Ack,
		len(packet.Data),
	)
	c.Target <- packet
}

func (c *Connection) receive() Packet {
	packet := <-c.Conn

	// This is to simulate network delay
	time.Sleep(time.Second)

	fmt.Printf(
		"%s <-- %s [%s] Seq=%d Ack=%d Len=%d\n",
		packet.Dest,
		packet.Source,
		packet.Flags.String(),
		packet.Seq,
		packet.Ack,
		len(packet.Data),
	)
	c.Packets = append(c.Packets, packet)
	return packet
}

func (c *Connection) Accept(conn Connection) Packet {
	packet := c.receive()
	if !packet.Flags.SYN {
		c.Reset()
	}

	c.Target = conn.Conn
	c.send(Packet{
		Source: c.Name,
		Dest:   conn.Name,
		Seq:    packet.Seq + 1,
		Ack:    packet.Seq,
		Flags: Flags{
			SYN: true,
			ACK: true,
		},
	})

	packet = c.receive()
	if !packet.Flags.ACK {
		c.Reset()
	}

	return c.receive()
}

func (c *Connection) Connect(conn Connection) {
	c.Target = conn.Conn
	c.send(Packet{
		Source: c.Name,
		Dest:   conn.Name,
		Seq:    0,
		Ack:    0,
		Flags: Flags{
			SYN: true,
		},
	})

	packet := c.receive()
	if !packet.Flags.SYN && !packet.Flags.ACK {
		c.Reset()
	}

	c.send(Packet{
		Source: c.Name,
		Dest:   conn.Name,
		Seq:    packet.Seq + 1,
		Ack:    packet.Seq,
		Flags: Flags{
			ACK: true,
		},
	})
}

func (c *Connection) Receive() Packet {
	packet := c.receive()

	if packet.Flags.RST {
		c.Reset()
	}
	if packet.Flags.FIN {
		c.Close()
	}

	return packet
}

func (c *Connection) Close() {
	lastPacket := c.Packets[len(c.Packets)-1]

	c.send(Packet{
		Source: lastPacket.Source,
		Dest:   lastPacket.Dest,
		Seq:    lastPacket.Seq + 1,
		Ack:    lastPacket.Seq,
		Flags: Flags{
			FIN: true,
			ACK: true,
		},
	})
}

func (c *Connection) Reset() {
	lastPacket := c.Packets[len(c.Packets)-1]

	c.send(Packet{
		Seq: lastPacket.Seq + 1,
		Ack: lastPacket.Seq,
		Flags: Flags{
			RST: true,
			ACK: true,
		},
	})
}

func (c *Connection) Send(bytes []byte) {
	lastPacket := c.Packets[len(c.Packets)-1]

	c.send(Packet{
		Source: lastPacket.Source,
		Dest:   lastPacket.Dest,
		Seq:    lastPacket.Seq + 1,
		Ack:    lastPacket.Seq,
		Flags: Flags{
			ACK: true,
		},
		Data: bytes,
	})
}
