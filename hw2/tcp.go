package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go server()
	time.Sleep(5 * time.Second)
	client()
}

func server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server received: %s\n", buffer[:n])

	_, err = conn.Write([]byte("SYN-ACK"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Server sent: SYN-ACK")

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server received: %s\n", buffer[:n])
}

func client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("SYN"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Client sent: SYN")

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client received: %s\n", buffer[:n])

	_, err = conn.Write([]byte("ACK"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Client sent: ACK")

	message := "Hello, World!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	fmt.Println("Client sent:", message)
}
