package store

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrWrongType    = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrNotInteger   = errors.New("value is not an integer or out of range")
)

type Store struct {
	strings map[string]string
	lists   map[string][]string
	sets    map[string]map[string]bool
	hashes  map[string]map[string]string
	mu      sync.RWMutex
}

func New() *Store {
	return &Store{
		strings: make(map[string]string),
		lists:   make(map[string][]string),
		sets:    make(map[string]map[string]bool),
		hashes:  make(map[string]map[string]string),
	}
}

// ─── String Commands ────────────────────────────────────────────

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strings[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.strings[key]
	return value, exists
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, strExists := s.strings[key]
	_, listExists := s.lists[key]
	_, setExists := s.sets[key]
	_, hashExists := s.hashes[key]

	exists := strExists || listExists || setExists || hashExists

	delete(s.strings, key)
	delete(s.lists, key)
	delete(s.sets, key)
	delete(s.hashes, key)

	return exists
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, strExists := s.strings[key]
	_, listExists := s.lists[key]
	_, setExists := s.sets[key]
	_, hashExists := s.hashes[key]

	return strExists || listExists || setExists || hashExists
}

func (s *Store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.strings[key]
	if !exists {
		s.strings[key] = "1"
		return 1, nil
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, ErrNotInteger
	}

	num++
	s.strings[key] = strconv.Itoa(num)
	return num, nil
}

func (s *Store) Decr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.strings[key]
	if !exists {
		s.strings[key] = "-1"
		return -1, nil
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, ErrNotInteger
	}

	num--
	s.strings[key] = strconv.Itoa(num)
	return num, nil
}

func (s *Store) MSet(pairs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < len(pairs)-1; i += 2 {
		s.strings[pairs[i]] = pairs[i+1]
	}
}

func (s *Store) MGet(keys []string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values := make([]string, len(keys))
	for i, key := range keys {
		if val, exists := s.strings[key]; exists {
			values[i] = val
		}
	}
	return values
}

// ─── List Commands ──────────────────────────────────────────────

func (s *Store) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	for _, val := range values {
		s.lists[key] = append([]string{val}, s.lists[key]...)
	}

	return len(s.lists[key]), nil
}

func (s *Store) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	s.lists[key] = append(s.lists[key], values...)
	return len(s.lists[key]), nil
}

func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return "", false, ErrWrongType
	}

	list, exists := s.lists[key]
	if !exists || len(list) == 0 {
		return "", false, nil
	}

	value := list[0]
	s.lists[key] = list[1:]
	return value, true, nil
}

func (s *Store) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return "", false, ErrWrongType
	}

	list, exists := s.lists[key]
	if !exists || len(list) == 0 {
		return "", false, nil
	}

	value := list[len(list)-1]
	s.lists[key] = list[:len(list)-1]
	return value, true, nil
}

func (s *Store) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return nil, ErrWrongType
	}

	list, exists := s.lists[key]
	if !exists {
		return []string{}, nil
	}

	length := len(list)

	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return []string{}, nil
	}

	return list[start : stop+1], nil
}

func (s *Store) LLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	return len(s.lists[key]), nil
}

// ─── Set Commands ───────────────────────────────────────────────

func (s *Store) SAdd(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	if s.sets[key] == nil {
		s.sets[key] = make(map[string]bool)
	}

	added := 0
	for _, member := range members {
		if !s.sets[key][member] {
			s.sets[key][member] = true
			added++
		}
	}

	return added, nil
}

func (s *Store) SRem(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	removed := 0
	for _, member := range members {
		if s.sets[key][member] {
			delete(s.sets[key], member)
			removed++
		}
	}

	return removed, nil
}

func (s *Store) SIsMember(key, member string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return false, ErrWrongType
	}

	return s.sets[key][member], nil
}

func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return nil, ErrWrongType
	}

	members := make([]string, 0, len(s.sets[key]))
	for member := range s.sets[key] {
		members = append(members, member)
	}

	return members, nil
}

func (s *Store) SCard(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	return len(s.sets[key]), nil
}

// ─── Hash Commands ──────────────────────────────────────────────

func (s *Store) HSet(key string, fields map[string]string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	if s.hashes[key] == nil {
		s.hashes[key] = make(map[string]string)
	}

	added := 0
	for field, value := range fields {
		if _, exists := s.hashes[key][field]; !exists {
			added++
		}
		s.hashes[key][field] = value
	}

	return added, nil
}

func (s *Store) HGet(key, field string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return "", false, ErrWrongType
	}

	hash, exists := s.hashes[key]
	if !exists {
		return "", false, nil
	}

	value, exists := hash[field]
	return value, exists, nil
}

func (s *Store) HDel(key string, fields ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}

	deleted := 0
	for _, field := range fields {
		if _, exists := s.hashes[key][field]; exists {
			delete(s.hashes[key], field)
			deleted++
		}
	}

	return deleted, nil
}

func (s *Store) HGetAll(key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return nil, ErrWrongType
	}

	hash, exists := s.hashes[key]
	if !exists {
		return map[string]string{}, nil
	}

	result := make(map[string]string, len(hash))
	for k, v := range hash {
		result[k] = v
	}

	return result, nil
}

func (s *Store) HExists(key, field string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return false, ErrWrongType
	}

	_, exists := s.hashes[key][field]
	return exists, nil
}