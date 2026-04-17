package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

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
		return errResponse("ERR invalid command format")
	}

	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]

	switch command {
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
	case "INCR":
		return s.handleIncr(args)
	case "DECR":
		return s.handleDecr(args)
	case "MSET":
		return s.handleMSet(args)
	case "MGET":
		return s.handleMGet(args)
	case "LPUSH":
		return s.handleLPush(args)
	case "RPUSH":
		return s.handleRPush(args)
	case "LPOP":
		return s.handleLPop(args)
	case "RPOP":
		return s.handleRPop(args)
	case "LRANGE":
		return s.handleLRange(args)
	case "LLEN":
		return s.handleLLen(args)
	case "SADD":
		return s.handleSAdd(args)
	case "SREM":
		return s.handleSRem(args)
	case "SISMEMBER":
		return s.handleSIsMember(args)
	case "SMEMBERS":
		return s.handleSMembers(args)
	case "SCARD":
		return s.handleSCard(args)
	case "HSET":
		return s.handleHSet(args)
	case "HGET":
		return s.handleHGet(args)
	case "HDEL":
		return s.handleHDel(args)
	case "HGETALL":
		return s.handleHGetAll(args)
	case "HEXISTS":
		return s.handleHExists(args)
	default:
		return errResponse(fmt.Sprintf("ERR unknown command '%s'", command))
	}
}

// ─── Helpers ────────────────────────────────────────────────────

func errResponse(msg string) resp.Value {
	return resp.Value{Type: resp.ERROR, Str: msg}
}

func okResponse() resp.Value {
	return resp.Value{Type: resp.STRING, Str: "OK"}
}

func intResponse(n int) resp.Value {
	return resp.Value{Type: resp.INTEGER, Num: n}
}

func nullResponse() resp.Value {
	return resp.Value{Type: resp.NULL}
}

func bulkResponse(s string) resp.Value {
	return resp.Value{Type: resp.BULK, Bulk: s}
}

// ─── String Handlers ────────────────────────────────────────────

func (s *Server) handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.STRING, Str: "PONG"}
	}
	return bulkResponse(args[0].Bulk)
}

func (s *Server) handleSet(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR SET requires key and value")
	}
	s.store.Set(args[0].Bulk, args[1].Bulk)
	return okResponse()
}

func (s *Server) handleGet(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR GET requires key")
	}
	value, exists := s.store.Get(args[0].Bulk)
	if !exists {
		return nullResponse()
	}
	return bulkResponse(value)
}

func (s *Server) handleDel(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR DEL requires key")
	}
	if s.store.Del(args[0].Bulk) {
		return intResponse(1)
	}
	return intResponse(0)
}

func (s *Server) handleExists(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR EXISTS requires key")
	}
	if s.store.Exists(args[0].Bulk) {
		return intResponse(1)
	}
	return intResponse(0)
}

func (s *Server) handleIncr(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR INCR requires key")
	}
	num, err := s.store.Incr(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(num)
}

func (s *Server) handleDecr(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR DECR requires key")
	}
	num, err := s.store.Decr(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(num)
}

func (s *Server) handleMSet(args []resp.Value) resp.Value {
	if len(args) < 2 || len(args)%2 != 0 {
		return errResponse("ERR MSET requires pairs of key value")
	}
	pairs := make([]string, len(args))
	for i, arg := range args {
		pairs[i] = arg.Bulk
	}
	s.store.MSet(pairs)
	return okResponse()
}

func (s *Server) handleMGet(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR MGET requires at least one key")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.Bulk
	}
	values := s.store.MGet(keys)
	array := make([]resp.Value, len(values))
	for i, val := range values {
		if val == "" {
			array[i] = nullResponse()
		} else {
			array[i] = bulkResponse(val)
		}
	}
	return resp.Value{Type: resp.ARRAY, Array: array}
}

// ─── List Handlers ──────────────────────────────────────────────

func (s *Server) handleLPush(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR LPUSH requires key and value")
	}
	key := args[0].Bulk
	values := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		values[i] = arg.Bulk
	}
	length, err := s.store.LPush(key, values...)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(length)
}

func (s *Server) handleRPush(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR RPUSH requires key and value")
	}
	key := args[0].Bulk
	values := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		values[i] = arg.Bulk
	}
	length, err := s.store.RPush(key, values...)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(length)
}

func (s *Server) handleLPop(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR LPOP requires key")
	}
	value, exists, err := s.store.LPop(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	if !exists {
		return nullResponse()
	}
	return bulkResponse(value)
}

func (s *Server) handleRPop(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR RPOP requires key")
	}
	value, exists, err := s.store.RPop(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	if !exists {
		return nullResponse()
	}
	return bulkResponse(value)
}

func (s *Server) handleLRange(args []resp.Value) resp.Value {
	if len(args) < 3 {
		return errResponse("ERR LRANGE requires key start stop")
	}

	start, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return errResponse("ERR start is not an integer")
	}

	stop, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return errResponse("ERR stop is not an integer")
	}

	values, err := s.store.LRange(args[0].Bulk, start, stop)
	if err != nil {
		return errResponse(err.Error())
	}

	array := make([]resp.Value, len(values))
	for i, val := range values {
		array[i] = bulkResponse(val)
	}
	return resp.Value{Type: resp.ARRAY, Array: array}
}

func (s *Server) handleLLen(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR LLEN requires key")
	}
	length, err := s.store.LLen(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(length)
}

// ─── Set Handlers ───────────────────────────────────────────────

func (s *Server) handleSAdd(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR SADD requires key and member")
	}
	key := args[0].Bulk
	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		members[i] = arg.Bulk
	}
	added, err := s.store.SAdd(key, members...)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(added)
}

func (s *Server) handleSRem(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR SREM requires key and member")
	}
	key := args[0].Bulk
	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		members[i] = arg.Bulk
	}
	removed, err := s.store.SRem(key, members...)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(removed)
}

func (s *Server) handleSIsMember(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR SISMEMBER requires key and member")
	}
	exists, err := s.store.SIsMember(args[0].Bulk, args[1].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	if exists {
		return intResponse(1)
	}
	return intResponse(0)
}

func (s *Server) handleSMembers(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR SMEMBERS requires key")
	}
	members, err := s.store.SMembers(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	array := make([]resp.Value, len(members))
	for i, member := range members {
		array[i] = bulkResponse(member)
	}
	return resp.Value{Type: resp.ARRAY, Array: array}
}

func (s *Server) handleSCard(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR SCARD requires key")
	}
	count, err := s.store.SCard(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(count)
}

// ─── Hash Handlers ──────────────────────────────────────────────

func (s *Server) handleHSet(args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 == 0 {
		return errResponse("ERR HSET requires key field value")
	}
	key := args[0].Bulk
	fields := make(map[string]string)
	for i := 1; i < len(args)-1; i += 2 {
		fields[args[i].Bulk] = args[i+1].Bulk
	}
	added, err := s.store.HSet(key, fields)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(added)
}

func (s *Server) handleHGet(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR HGET requires key and field")
	}
	value, exists, err := s.store.HGet(args[0].Bulk, args[1].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	if !exists {
		return nullResponse()
	}
	return bulkResponse(value)
}

func (s *Server) handleHDel(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR HDEL requires key and field")
	}
	key := args[0].Bulk
	fields := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		fields[i] = arg.Bulk
	}
	deleted, err := s.store.HDel(key, fields...)
	if err != nil {
		return errResponse(err.Error())
	}
	return intResponse(deleted)
}

func (s *Server) handleHGetAll(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return errResponse("ERR HGETALL requires key")
	}
	hash, err := s.store.HGetAll(args[0].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	array := make([]resp.Value, 0, len(hash)*2)
	for field, value := range hash {
		array = append(array, bulkResponse(field))
		array = append(array, bulkResponse(value))
	}
	return resp.Value{Type: resp.ARRAY, Array: array}
}

func (s *Server) handleHExists(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return errResponse("ERR HEXISTS requires key and field")
	}
	exists, err := s.store.HExists(args[0].Bulk, args[1].Bulk)
	if err != nil {
		return errResponse(err.Error())
	}
	if exists {
		return intResponse(1)
	}
	return intResponse(0)
}