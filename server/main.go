package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	_, err = conn.Write([]byte("World"))
	if err != nil {
		log.Fatal(err)
	}
}
