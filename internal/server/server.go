package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/AdeshDeshmukh/crimson/internal/resp"
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

	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)

	for {
		value, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", clientAddr)
			} else {
				log.Printf("Error reading from %s: %v", clientAddr, err)
			}
			return
		}

		response := s.processCommand(value)

		err = writer.Write(response)
		if err != nil {
			log.Printf("Error writing to %s: %v", clientAddr, err)
			return
		}
	}
}

func (s *Server) processCommand(value resp.Value) resp.Value {
	if value.Type != resp.ARRAY || len(value.Array) == 0 {
		return resp.Value{Type: resp.ERROR, Str: "ERR invalid command format"}
	}

	command := value.Array[0].Bulk
	args := value.Array[1:]

	switch command {
	case "ping", "PING":
		return s.handlePing(args)
	case "set", "SET":
		return s.handleSet(args)
	case "get", "GET":
		return s.handleGet(args)
	case "del", "DEL":
		return s.handleDel(args)
	case "exists", "EXISTS":
		return s.handleExists(args)
	default:
		return resp.Value{Type: resp.ERROR, Str: fmt.Sprintf("ERR unknown command '%s'", command)}
	}
}

func (s *Server) handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.STRING, Str: "PONG"}
	}
	return resp.Value{Type: resp.STRING, Str: args[0].Bulk}
}

func (s *Server) handleSet(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Type: resp.ERROR, Str: "ERR SET requires key and value"}
	}
	s.store.Set(args[0].Bulk, args[1].Bulk)
	return resp.Value{Type: resp.STRING, Str: "OK"}
}

func (s *Server) handleGet(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, Str: "ERR GET requires key"}
	}
	value, exists := s.store.Get(args[0].Bulk)
	if !exists {
		return resp.Value{Type: resp.NULL}
	}
	return resp.Value{Type: resp.BULK, Bulk: value}
}

func (s *Server) handleDel(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, Str: "ERR DEL requires key"}
	}
	if s.store.Del(args[0].Bulk) {
		return resp.Value{Type: resp.INTEGER, Num: 1}
	}
	return resp.Value{Type: resp.INTEGER, Num: 0}
}

func (s *Server) handleExists(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, Str: "ERR EXISTS requires key"}
	}
	if s.store.Exists(args[0].Bulk) {
		return resp.Value{Type: resp.INTEGER, Num: 1}
	}
	return resp.Value{Type: resp.INTEGER, Num: 0}
}