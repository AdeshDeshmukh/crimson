# In your crimson directory
cat > README.md << 'EOF'
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