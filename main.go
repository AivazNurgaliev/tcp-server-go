package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener

	go s.acceptLoop()

	<-s.quitch

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("new connection established: ", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}

		msg := buffer[:n]
		fmt.Printf("%s", string(msg))
	}
}

func main() {
	server := NewServer(":8000")
	log.Fatal(server.Start())
}
