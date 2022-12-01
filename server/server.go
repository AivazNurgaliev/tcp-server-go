package server

import (
	"fmt"
	"net"
)

type Message struct {
	From    string
	Payload []byte
}

type Server struct {
	ListenAddr string
	Listener   net.Listener
	Quitch     chan struct{}
	Msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Quitch:     make(chan struct{}),
		Msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.Listener = listener

	go s.acceptLoop()

	<-s.Quitch
	close(s.Msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.Listener.Accept()
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

		s.Msgch <- Message{
			From:    conn.RemoteAddr().String(),
			Payload: buffer[:n],
		}

	}
}
