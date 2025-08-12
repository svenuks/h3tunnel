package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"strings"

	"github.com/quic-go/quic-go"
	"h3tunnel/internal/auth"
	"h3tunnel/internal/config"
	"h3tunnel/internal/protocol"
	"h3tunnel/internal/proxy"
)

func main() {
	log.Println("Starting server...")

	cfg, err := config.LoadServerConfig("config.server.json")
	if err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}

	log.Printf("Server config loaded successfully")

	tlsConfig, err := protocol.GenerateTLSConfig()
	if err != nil {
		log.Fatalf("Failed to generate TLS config: %v", err)
	}

	listener, err := quic.ListenAddr(cfg.ListenAddress, tlsConfig, nil)

	if err != nil {
		log.Fatalf("Failed to start QUIC listener: %v", err)
	}

	log.Printf("Server listening on %s", listener.Addr())

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, cfg.Password)
	}
}

func handleConnection(conn quic.Connection, password string) {
	log.Printf("Accepted new connection from: %s", conn.RemoteAddr())

	if err := auth.HandleServerAuth(context.Background(), conn, password); err != nil {
		log.Printf("Authentication failed for %s: %v", conn.RemoteAddr(), err)
		conn.CloseWithError(1, "Authentication failed")
		return
	}

	log.Printf("Client %s authenticated successfully. Waiting for proxy streams...", conn.RemoteAddr())

	// After auth, loop to accept streams for proxying
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("Failed to accept stream from %s: %v", conn.RemoteAddr(), err)
			return // Connection is likely closed
		}
		go handleProxyStream(stream)
	}
}

func handleProxyStream(stream quic.Stream) {
	// The first thing on the stream should be the remote address
	reader := bufio.NewReader(stream)
	remoteAddr, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read remote address from stream: %v", err)
		stream.Close()
		return
	}
	remoteAddr = strings.TrimSpace(remoteAddr)

	log.Printf("Attempting to dial remote address: %s", remoteAddr)
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Failed to dial remote address %s: %v", remoteAddr, err)
		stream.Close()
		return
	}

	log.Printf("Successfully connected to %s. Relaying traffic...", remoteAddr)
	proxy.Relay(remoteConn, stream)
	log.Printf("Finished relaying for %s", remoteAddr)
}
