package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/AdeshDeshmukh/crimson/internal/store"
)

type Server struct {
	addr     string
	store    *store.Store
	listener net.Listener
}

func New(addr string) *Server {
	return &Server{
		addr:  addr,
		store: store.New(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	log.Printf("Crimson server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", clientAddr)
			} else {
				log.Printf("Error reading from %s: %v", clientAddr, err)
			}
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		response := s.processCommand(line)

		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Printf("Error writing to %s: %v", clientAddr, err)
			return
		}
	}
}

func (s *Server) processCommand(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "ERROR: empty command"
	}

	cmd := strings.ToUpper(parts[0])
	args := parts[1:]

	switch cmd {
	case "PING":
		return s.handlePing(args)
	case "SET":
		return s.handleSet(args)
	case "GET":
		return s.handleGet(args)
	case "DEL":
		return s.handleDel(args)
	case "EXISTS":
		return s.handleExists(args)
	default:
		return fmt.Sprintf("ERROR: unknown command '%s'", cmd)
	}
}

func (s *Server) handlePing(args []string) string {
	if len(args) == 0 {
		return "PONG"
	}
	return strings.Join(args, " ")
}

func (s *Server) handleSet(args []string) string {
	if len(args) < 2 {
		return "ERROR: SET requires key and value"
	}
	key := args[0]
	value := strings.Join(args[1:], " ")
	s.store.Set(key, value)
	return "OK"
}

func (s *Server) handleGet(args []string) string {
	if len(args) < 1 {
		return "ERROR: GET requires key"
	}
	value, exists := s.store.Get(args[0])
	if !exists {
		return "nil"
	}
	return value
}

func (s *Server) handleDel(args []string) string {
	if len(args) < 1 {
		return "ERROR: DEL requires key"
	}
	if s.store.Del(args[0]) {
		return "1"
	}
	return "0"
}

func (s *Server) handleExists(args []string) string {
	if len(args) < 1 {
		return "ERROR: EXISTS requires key"
	}
	if s.store.Exists(args[0]) {
		return "1"
	}
	return "0"
}
