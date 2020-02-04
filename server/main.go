package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type server struct {
	l       net.Listener
	clients []net.Conn
	mu      sync.Mutex
}

func (s *server) listen() error {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	s.l = l
	return nil
}

func (s *server) accept() (net.Conn, error) {
	conn, err := s.l.Accept()
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients = append(s.clients, conn)
	return conn, nil
}

func (s *server) receive(conn net.Conn) {
	for {
		b := make([]byte, 1024)
		_, err := conn.Read(b)
		if err == io.EOF {
			// クライアントが接続切断
			s.delete(conn)
			fmt.Println("クライアント離脱")
			return
		} else if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(conn.RemoteAddr(), string(b))
		s.broadcast(string(b), conn)
	}
}

func (s *server) delete(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, c := range s.clients {
		if c == conn {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			break
		}
	}
}

func (s *server) broadcast(msg string, conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.clients {
		if c == conn {
			continue
		}
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Print("broadcast failed")
		}
	}
}

//
//func (s *server) close() error {
//	err := s.l.Close()
//	return err
//}

func main() {
	s := new(server)

	// サーバリッスン開始
	err := s.listen()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// クライアントからの接続待ち
		conn, err := s.accept()
		if err != nil {
			log.Fatal(err)
		}

		// クライアントからのメッセージ受信待ち
		go s.receive(conn)
	}
}
