package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
