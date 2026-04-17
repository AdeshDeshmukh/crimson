package store

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrWrongType   = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrNotInteger  = errors.New("ERR value is not an integer or out of range")
)

type Store struct {
	strings map[string]string
	lists   map[string][]string
	sets    map[string]map[string]bool
	hashes  map[string]map[string]string
	expiry  map[string]time.Time
	mu      sync.RWMutex
}

func New() *Store {
	s := &Store{
		strings: make(map[string]string),
		lists:   make(map[string][]string),
		sets:    make(map[string]map[string]bool),
		hashes:  make(map[string]map[string]string),
		expiry:  make(map[string]time.Time),
	}
	go s.cleanupLoop()
	return s
}

func (s *Store) isExpired(key string) bool {
	expireAt, exists := s.expiry[key]
	if !exists {
		return false
	}
	return time.Now().After(expireAt)
}

func (s *Store) deleteKey(key string) {
	delete(s.strings, key)
	delete(s.lists, key)
	delete(s.sets, key)
	delete(s.hashes, key)
	delete(s.expiry, key)
}

func (s *Store) cleanupLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for key := range s.expiry {
			if s.isExpired(key) {
				s.deleteKey(key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Store) SetExpiry(key string, duration time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, strExists := s.strings[key]
	_, listExists := s.lists[key]
	_, setExists := s.sets[key]
	_, hashExists := s.hashes[key]
	if !strExists && !listExists && !setExists && !hashExists {
		return false
	}
	s.expiry[key] = time.Now().Add(duration)
	return true
}

func (s *Store) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, strExists := s.strings[key]
	_, listExists := s.lists[key]
	_, setExists := s.sets[key]
	_, hashExists := s.hashes[key]
	if !strExists && !listExists && !setExists && !hashExists {
		return -2
	}
	expireAt, hasExpiry := s.expiry[key]
	if !hasExpiry {
		return -1
	}
	remaining := time.Until(expireAt)
	if remaining <= 0 {
		return -2
	}
	return int(remaining.Seconds())
}

func (s *Store) Persist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.expiry[key]
	if exists {
		delete(s.expiry, key)
	}
	return exists
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strings[key] = value
	if ttl > 0 {
		s.expiry[key] = time.Now().Add(ttl)
	} else {
		delete(s.expiry, key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return "", false
	}
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
	s.deleteKey(key)
	return exists
}

func (s *Store) Exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return false
	}
	_, strExists := s.strings[key]
	_, listExists := s.lists[key]
	_, setExists := s.sets[key]
	_, hashExists := s.hashes[key]
	return strExists || listExists || setExists || hashExists
}

func (s *Store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
	}
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
	if s.isExpired(key) {
		s.deleteKey(key)
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	values := make([]string, len(keys))
	for i, key := range keys {
		if s.isExpired(key) {
			s.deleteKey(key)
			continue
		}
		if val, exists := s.strings[key]; exists {
			values[i] = val
		}
	}
	return values
}

func (s *Store) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
	}
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
	if s.isExpired(key) {
		s.deleteKey(key)
	}
	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}
	s.lists[key] = append(s.lists[key], values...)
	return len(s.lists[key]), nil
}

func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return "", false, nil
	}
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
	if s.isExpired(key) {
		s.deleteKey(key)
		return "", false, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return []string{}, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return 0, nil
	}
	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}
	return len(s.lists[key]), nil
}

func (s *Store) SAdd(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
	}
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
	if s.isExpired(key) {
		s.deleteKey(key)
		return 0, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return false, nil
	}
	if _, exists := s.strings[key]; exists {
		return false, ErrWrongType
	}
	return s.sets[key][member], nil
}

func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return []string{}, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return 0, nil
	}
	if _, exists := s.strings[key]; exists {
		return 0, ErrWrongType
	}
	return len(s.sets[key]), nil
}

func (s *Store) HSet(key string, fields map[string]string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return "", false, nil
	}
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
	if s.isExpired(key) {
		s.deleteKey(key)
		return 0, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return map[string]string{}, nil
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return false, nil
	}
	if _, exists := s.strings[key]; exists {
		return false, ErrWrongType
	}
	_, exists := s.hashes[key][field]
	return exists, nil
}

func (s *Store) Keys(pattern string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []string
	for key := range s.strings {
		if !s.isExpired(key) && matchPattern(pattern, key) {
			keys = append(keys, key)
		}
	}
	for key := range s.lists {
		if !s.isExpired(key) && matchPattern(pattern, key) {
			keys = append(keys, key)
		}
	}
	for key := range s.sets {
		if !s.isExpired(key) && matchPattern(pattern, key) {
			keys = append(keys, key)
		}
	}
	for key := range s.hashes {
		if !s.isExpired(key) && matchPattern(pattern, key) {
			keys = append(keys, key)
		}
	}
	return keys
}

func matchPattern(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	return matchGlob(pattern, key)
}

func matchGlob(pattern, str string) bool {
	p := 0
	s := 0
	for p < len(pattern) && s < len(str) {
		if pattern[p] == '*' {
			for p < len(pattern) && pattern[p] == '*' {
				p++
			}
			if p == len(pattern) {
				return true
			}
			for s < len(str) {
				if matchGlob(pattern[p:], str[s:]) {
					return true
				}
				s++
			}
			return false
		} else if pattern[p] == '?' {
			p++
			s++
		} else if pattern[p] == str[s] {
			p++
			s++
		} else {
			return false
		}
	}
	for p < len(pattern) && pattern[p] == '*' {
		p++
	}
	return p == len(pattern) && s == len(str)
}

func (s *Store) Scan(cursor int, pattern string, count int) (int, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	allKeys := make([]string, 0)
	for key := range s.strings {
		if !s.isExpired(key) {
			allKeys = append(allKeys, key)
		}
	}
	for key := range s.lists {
		if !s.isExpired(key) {
			allKeys = append(allKeys, key)
		}
	}
	for key := range s.sets {
		if !s.isExpired(key) {
			allKeys = append(allKeys, key)
		}
	}
	for key := range s.hashes {
		if !s.isExpired(key) {
			allKeys = append(allKeys, key)
		}
	}
	if count <= 0 {
		count = 10
	}
	if cursor >= len(allKeys) {
		return 0, []string{}
	}
	end := cursor + count
	if end > len(allKeys) {
		end = len(allKeys)
	}
	chunk := allKeys[cursor:end]
	var result []string
	for _, key := range chunk {
		if matchPattern(pattern, key) {
			result = append(result, key)
		}
	}
	nextCursor := end
	if nextCursor >= len(allKeys) {
		nextCursor = 0
	}
	return nextCursor, result
}

func (s *Store) Type(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(key) {
		return "none"
	}
	if _, exists := s.strings[key]; exists {
		return "string"
	}
	if _, exists := s.lists[key]; exists {
		return "list"
	}
	if _, exists := s.sets[key]; exists {
		return "set"
	}
	if _, exists := s.hashes[key]; exists {
		return "hash"
	}
	return "none"
}

func (s *Store) Rename(key, newKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isExpired(key) {
		s.deleteKey(key)
		return ErrKeyNotFound
	}
	if _, exists := s.strings[key]; exists {
		s.strings[newKey] = s.strings[key]
		delete(s.strings, key)
		if expiry, hasExpiry := s.expiry[key]; hasExpiry {
			s.expiry[newKey] = expiry
			delete(s.expiry, key)
		}
		return nil
	}
	if _, exists := s.lists[key]; exists {
		s.lists[newKey] = s.lists[key]
		delete(s.lists, key)
		if expiry, hasExpiry := s.expiry[key]; hasExpiry {
			s.expiry[newKey] = expiry
			delete(s.expiry, key)
		}
		return nil
	}
	if _, exists := s.sets[key]; exists {
		s.sets[newKey] = s.sets[key]
		delete(s.sets, key)
		if expiry, hasExpiry := s.expiry[key]; hasExpiry {
			s.expiry[newKey] = expiry
			delete(s.expiry, key)
		}
		return nil
	}
	if _, exists := s.hashes[key]; exists {
		s.hashes[newKey] = s.hashes[key]
		delete(s.hashes, key)
		if expiry, hasExpiry := s.expiry[key]; hasExpiry {
			s.expiry[newKey] = expiry
			delete(s.expiry, key)
		}
		return nil
	}
	return ErrKeyNotFound
}

func (s *Store) DBSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.strings) + len(s.lists) + len(s.sets) + len(s.hashes)
}

func (s *Store) FlushDB() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strings = make(map[string]string)
	s.lists = make(map[string][]string)
	s.sets = make(map[string]map[string]bool)
	s.hashes = make(map[string]map[string]string)
	s.expiry = make(map[string]time.Time)
}