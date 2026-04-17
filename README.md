# In your crimson directory
<div align="center">

# 🔴 Crimson

### A Production-Grade Redis Clone Built in Go

*Understanding databases from first principles*

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge)](CONTRIBUTING.md)

[Features](#-features) • [Installation](#-installation) • [Architecture](#-architecture) • [Commands](#-commands) • [Roadmap](#-roadmap)

</div>

---

## 📖 Table of Contents

- [What is Crimson?](#what-is-crimson)
- [Why Build This?](#why-build-this)
- [Features](#-features)
- [Quick Start](#-quick-start)
- [Architecture](#-architecture)
- [Supported Commands](#-supported-commands)
- [Usage Examples](#-usage-examples)
- [How It Works](#-how-it-works)
- [Project Structure](#-project-structure)
- [Testing](#-testing)
- [Performance](#-performance)
- [Future Roadmap](#-future-roadmap)
- [Contributing](#-contributing)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)

---

## What is Crimson?

**Crimson** is a **Redis-compatible in-memory data store** built entirely from scratch in Go. It implements the full Redis RESP (REdis Serialization Protocol), making it compatible with any standard Redis client.

This project is built to deeply understand:
- 🔌 **Network Programming** - TCP servers, socket programming, concurrent connections
- 📡 **Wire Protocols** - Binary protocol design and implementation
- 🗄️ **Database Internals** - Storage engines, data structures, indexing
- 💾 **Persistence** - Write-ahead logging (AOF), durability guarantees
- ⚡ **Concurrent Systems** - Goroutines, mutexes, race conditions
- 🏗️ **System Design** - Building production-grade distributed systems

> **Note:** This is an educational project demonstrating database fundamentals. For production use, please use [official Redis](https://redis.io/).

---

## Why Build This?

"Don't build applications. Build products. Build systems."
- Anuj Bhaiya


Most students build:
- ❌ Todo apps
- ❌ Weather apps  
- ❌ CRUD APIs

This project builds:
- ✅ A real database server
- ✅ A binary wire protocol
- ✅ Concurrent data structures
- ✅ Persistence mechanisms
- ✅ Real-time messaging systems

**Learning by building the tools you use daily.**

---

## ✨ Features

### **Core Features**

| Feature | Status | Description |
|---------|--------|-------------|
| **TCP Server** | ✅ Complete | High-performance concurrent server |
| **RESP Protocol** | ✅ Complete | Full Redis wire protocol implementation |
| **Data Types** | ✅ Complete | Strings, Lists, Sets, Hashes |
| **TTL/Expiry** | ✅ Complete | Automatic key expiration with background cleanup |
| **AOF Persistence** | ✅ Complete | Append-only file for data durability |
| **Pub/Sub** | ✅ Complete | Real-time publish/subscribe messaging |
| **Transactions** | ✅ Complete | ACID transactions via MULTI/EXEC |

### **Commands Implemented: 30+**

<details>
<summary><b>String Commands (9)</b></summary>

- `PING` - Test connection
- `SET` - Set key to value (with EX/PX options)
- `GET` - Get value of key
- `DEL` - Delete key
- `EXISTS` - Check if key exists
- `INCR` - Increment integer value
- `DECR` - Decrement integer value
- `MSET` - Set multiple keys
- `MGET` - Get multiple keys

</details>

<details>
<summary><b>List Commands (6)</b></summary>

- `LPUSH` - Push to list head
- `RPUSH` - Push to list tail
- `LPOP` - Pop from list head
- `RPOP` - Pop from list tail
- `LRANGE` - Get range of elements
- `LLEN` - Get list length

</details>

<details>
<summary><b>Set Commands (5)</b></summary>

- `SADD` - Add member to set
- `SREM` - Remove member from set
- `SMEMBERS` - Get all members
- `SISMEMBER` - Check membership
- `SCARD` - Get set size

</details>

<details>
<summary><b>Hash Commands (5)</b></summary>

- `HSET` - Set hash field
- `HGET` - Get hash field
- `HDEL` - Delete hash field
- `HGETALL` - Get all fields
- `HEXISTS` - Check if field exists

</details>

<details>
<summary><b>TTL Commands (3)</b></summary>

- `EXPIRE` - Set key expiry
- `TTL` - Get remaining time
- `PERSIST` - Remove expiry

</details>

<details>
<summary><b>Pub/Sub Commands (2)</b></summary>

- `SUBSCRIBE` - Subscribe to channel
- `PUBLISH` - Publish message

</details>

<details>
<summary><b>Transaction Commands (3)</b></summary>

- `MULTI` - Start transaction
- `EXEC` - Execute transaction
- `DISCARD` - Cancel transaction

</details>

---

## 🚀 Quick Start

### **Prerequisites**

- Go 1.22 or higher
- Git
- redis-cli (for testing)

### **Installation**

```bash
# Clone repository
git clone https://github.com/AdeshDeshmukh/crimson.git
cd crimson

# Build
make build

# Run
make run

Server will start on port 6379 (Redis default).

Connect with redis-cli

redis-cli -p 6379

127.0.0.1:6379> PING
PONG

127.0.0.1:6379> SET name "Adesh"
OK

127.0.0.1:6379> GET name
"Adesh"

Connect with Go

package main

import (
    "github.com/go-redis/redis/v8"
    "context"
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    ctx := context.Background()
    
    client.Set(ctx, "key", "value", 0)
    val, _ := client.Get(ctx, "key").Result()
    println(val) // Output: value
}


🏗️ Architecture

┌─────────────────────────────────────────────────────────────┐
│                         CLIENT                              │
│                  (redis-cli / go-redis / etc)               │
└────────────────────────┬────────────────────────────────────┘
                         │ TCP Connection
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    TCP SERVER LAYER                         │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Connection 1 │  │ Connection 2 │  │ Connection N │     │
│  │  (goroutine) │  │  (goroutine) │  │  (goroutine) │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │             │
└─────────┼──────────────────┼──────────────────┼─────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   RESP PROTOCOL LAYER                       │
│                                                             │
│  ┌──────────────┐              ┌──────────────┐            │
│  │    Parser    │              │    Writer    │            │
│  │ (bytes → Go) │              │ (Go → bytes) │            │
│  └──────┬───────┘              └──────▲───────┘            │
└─────────┼──────────────────────────────┼───────────────────┘
          │                              │
          ▼                              │
┌─────────────────────────────────────────────────────────────┐
│                   COMMAND EXECUTOR                          │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  String  │  │   List   │  │   Set    │  │   Hash   │   │
│  │ Handlers │  │ Handlers │  │ Handlers │  │ Handlers │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        ▼             ▼             ▼             ▼
┌─────────────────────────────────────────────────────────────┐
│                      DATA STORE                             │
│                                                             │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────┐│
│  │  Strings   │ │   Lists    │ │    Sets    │ │  Hashes  ││
│  │map[k]v     │ │map[k][]v   │ │map[k]set   │ │map[k]map ││
│  └────────────┘ └────────────┘ └────────────┘ └──────────┘│
│                                                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │            Expiry Map (TTL tracking)               │    │
│  │            map[key]expiryTimestamp                 │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
│  🔒 Protected by sync.RWMutex (thread-safe)                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   PERSISTENCE LAYER                         │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                 AOF (Append Only File)               │  │
│  │                                                      │  │
│  │  Every write → appended to crimson.aof              │  │
│  │  On restart → replay all commands                   │  │
│  │                                                      │  │
│  │  *3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nAdesh\r\n   │  │
│  │  *3\r\n$3\r\nSET\r\n$3\r\nage\r\n$2\r\n19\r\n       │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

                         ALSO

┌─────────────────────────────────────────────────────────────┐
│                      PUB/SUB SYSTEM                         │
│                                                             │
│  Channel Map:  map[channelName][]*Subscriber                │
│                                                             │
│  Publisher → finds subscribers → broadcasts message         │
│  Subscribers listen on Go channels (buffered, size 100)    │
└─────────────────────────────────────────────────────────────┘

📋 Supported Commands
String Operations

SET key value [EX seconds] [PX milliseconds]
GET key
DEL key
EXISTS key
INCR key
DECR key
MSET key1 value1 key2 value2 ...
MGET key1 key2 key3 ...

List Operations

LPUSH key value [value ...]
RPUSH key value [value ...]
LPOP key
RPOP key
LRANGE key start stop
LLEN key

Set Operations

SADD key member [member ...]
SREM key member [member ...]
SMEMBERS key
SISMEMBER key member
SCARD key

Hash Operations

HSET key field value [field value ...]
HGET key field
HDEL key field [field ...]
HGETALL key
HEXISTS key field

TTL Operations

EXPIRE key seconds
TTL key
PERSIST key

Pub/Sub

SUBSCRIBE channel [channel ...]
PUBLISH channel message

Transactions

MULTI
<commands...>
EXEC
DISCARD

💡 Usage Examples
Session Management

# Set session with 30 minute expiry
SET session:user123 "eyJhbGc..." EX 1800

# Check remaining time
TTL session:user123
# (integer) 1799

# Remove expiry
PERSIST session:user123

Real-time Chat

# Terminal 1 - Subscriber
SUBSCRIBE chat:room1

# Terminal 2 - Publisher
PUBLISH chat:room1 "Hello everyone!"
# (integer) 1  ← number of subscribers

# Terminal 1 receives:
# 1) "message"
# 2) "chat:room1"
# 3) "Hello everyone!"

Atomic Bank Transfer

SET account:alice 1000
SET account:bob 500

MULTI
DECR account:alice
INCR account:bob
EXEC

# Both operations succeed or both fail
# No partial transfers

🔍 How It Works

1. TCP Server

// Accept connections on port 6379
listener, _ := net.Listen("tcp", ":6379")

for {
    conn, _ := listener.Accept()
    go handleConnection(conn)  // Each client in own goroutine
}

Why it matters: Handles thousands of concurrent clients efficiently using Go's lightweight goroutines.


2. RESP Protocol
Client sends:

*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nAdesh\r\n

Parser converts to:

Value{
    Type: ARRAY,
    Array: [
        {Type: BULK, Bulk: "SET"},
        {Type: BULK, Bulk: "name"},
        {Type: BULK, Bulk: "Adesh"},
    ]
}

Server responds:

+OK\r\n

Why it matters: Binary-safe protocol that can handle any data including nulls, newlines, and special characters.




📁 Project Structure

crimson/
├── cmd/
│   └── crimson/
│       └── main.go              # Entry point, server startup
├── internal/
│   ├── aof/
│   │   └── aof.go              # Append-only file persistence
│   ├── pubsub/
│   │   └── pubsub.go           # Publish/Subscribe system
│   ├── resp/
│   │   ├── parser.go           # RESP protocol parser
│   │   └── writer.go           # RESP protocol writer
│   ├── server/
│   │   └── server.go           # TCP server, command routing
│   └── store/
│       └── store.go            # In-memory data structures
├── docs/
│   └── ARCHITECTURE.md         # Detailed architecture docs
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
├── README.md                   # This file
├── LICENSE                     # MIT License
└── crimson.aof                 # Persistence file (auto-created)



🧪 Testing

Manual Testing

# Start server
make run

# In another terminal
redis-cli -p 6379

# Run commands
PING
SET test "value"
GET test

Automated Testing

# Run all tests
make test

# Run with coverage
make test-coverage

# Run with race detector
go test -race ./...

Benchmarking

# Using redis-benchmark
redis-benchmark -p 6379 -t set,get -n 100000 -q

# Custom benchmark
make bench


🗺️ Future Roadmap
Phase 8: Advanced Commands
 KEYS pattern - Pattern-based key search
 SCAN cursor - Iterative key scanning
 TYPE key - Get key type
 RENAME key newkey - Rename keys
 SORT - Sort list/set/hash values
Phase 9: More Data Types
 Sorted Sets - ZADD, ZRANGE, ZRANK, ZINCRBY
 Bitmaps - SETBIT, GETBIT, BITCOUNT
 HyperLogLog - PFADD, PFCOUNT, PFMERGE
 Streams - XADD, XREAD, XRANGE
Phase 10: Server Management
 INFO - Server statistics
 CONFIG GET/SET - Runtime configuration
 DBSIZE - Number of keys
 FLUSHDB - Clear database
 SAVE / BGSAVE - Manual snapshots
Phase 11: Advanced Features
 RDB Snapshots - Point-in-time backups
 Lua Scripting - EVAL, EVALSHA
 Pipelining - Batch command execution
 Blocking Operations - BLPOP, BRPOP
 Geospatial - GEOADD, GEORADIUS
Phase 12: Replication
 Master-Replica setup
 Asynchronous replication
 Replica promotion
 Read scaling
Phase 13: Clustering
 Data sharding across nodes
 Hash slot allocation
 Cluster discovery
 Automatic failover
Phase 14: Observability
 Prometheus metrics export
 Logging levels (debug, info, warn, error)
 Slow log tracking
 Connection pooling stats
 Memory profiling
Phase 15: Security
 Password authentication (AUTH command)
 ACL (Access Control Lists)
 TLS/SSL support
 Command renaming/disabling



🤝 Contributing
Contributions are welcome! Here's how you can help:

Ways to Contribute
Bug Reports - Found a bug? Open an issue
Feature Requests - Have an idea? Suggest it
Code Contributions - Submit a PR
Documentation - Improve docs
Testing - Add test cases


Development Setup

# Fork the repository
git clone https://github.com/YOUR_USERNAME/crimson.git
cd crimson

# Create a branch
git checkout -b feature/your-feature

# Make changes
# ...

# Test
make test

# Commit
git commit -m "feat: add awesome feature"

# Push
git push origin feature/your-feature

# Open PR on GitHub


📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

MIT License

Copyright (c) 2026 Adesh Deshmukh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.


🙏 Acknowledgments
Inspiration
Redis - The amazing in-memory database this project emulates
Build Your Own Redis - Excellent learning resource
CodeCrafters - For the Redis challenge
Anuj Bhaiya - For the guidance to build products, not just applications
Learning Resources
Redis Protocol Specification
Redis Internals
Designing Data-Intensive Applications by Martin Kleppmann
Go Concurrency Patterns
Tools & Libraries
Go - The language that makes this possible
redis-cli - For testing
Homebrew - Package manager for macOS

👨‍💻 Author
Adesh Deshmukh

🎓 B.Tech in Electronics & Telecommunications Engineering
🏫 SGGS Institute of Engineering and Technology, Nanded
💻 Passionate about Systems Programming, Databases, and Distributed Systems
🏆 900+ problems solved across competitive programming platforms
Connect:

GitHub: @AdeshDeshmukh
LinkedIn: Adesh Deshmukh
LeetCode: Rating 1528
Codeforces: Rating 1114
Email: adeshkd123@gmail.com
💬 Testimonials
"Understanding Redis by building it from scratch is one of the best ways to learn database internals. Crimson demonstrates a deep understanding of systems programming."

"The code is clean, well-structured, and professionally organized. Great learning resource for anyone wanting to understand how databases work."

<div align="center">
⭐ Star This Project
If you found this helpful, please consider giving it a star!

Built with ❤️ and lots of ☕ by Adesh Deshmukh

⬆ Back to Top

</div>
🔗 Related Projects
go-redis - Go Redis client
redis - Official Redis source
KeyDB - Multi-threaded Redis fork
Dragonfly - Modern Redis alternative

📚 Blog Posts
Building Crimson: Part 1 - TCP Server
Building Crimson: Part 2 - RESP Protocol
Building Crimson: Part 3 - Concurrency
Building Crimson: Part 4 - Persistence
