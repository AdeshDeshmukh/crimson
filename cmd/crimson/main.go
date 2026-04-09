package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdeshDeshmukh/crimson/internal/server"
)

const (
	Version = "0.1.0"
	Banner  = `
    ╔═══════════════════════════════════════╗
    ║                                       ║
    ║        🔴  C R I M S O N  🔴           ║
    ║                                       ║
    ║   Redis-compatible data store in Go   ║
    ║            Version %s                 ║
    ║                                       ║
    ╚═══════════════════════════════════════╝
    `
)

func main() {
	fmt.Printf(Banner, Version)

	srv := server.New(":6379")

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Shutting down gracefully...")
	log.Println("Server stopped")
}
