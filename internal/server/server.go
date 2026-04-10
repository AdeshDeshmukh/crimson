package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
}

func New(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", s.listenAddr, err)
	}
	s.listener = listener

	log.Printf("Server listening on %s", s.listenAddr)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// This is the expected error when Stop() is called.
				// It signals that the server should shut down gracefully.
				return nil
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		log.Println("Stopping server...")
		// Closing the listener will cause the Accept() loop in Start() to unblock and return an error.
		s.listener.Close()
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", conn.RemoteAddr())
			} else {
				log.Printf("Error reading from client %s: %v", conn.RemoteAddr(), err)
			}
			return
		}
		log.Printf("Received %d bytes from %s: %s", n, conn.RemoteAddr(), string(buf[:n]))
	}
}