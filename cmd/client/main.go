package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"sync"

	"github.com/quic-go/quic-go"
	"h3tunnel/internal/auth"
	"h3tunnel/internal/config"
	"h3tunnel/internal/proxy"
)

func main() {
	log.Println("Starting client...")

	cfg, err := config.LoadClientConfig("config.client.json")
	if err != nil {
		log.Fatalf("Failed to load client config: %v", err)
	}

	log.Printf("Attempting to connect to server at %s...", cfg.ServerAddress)

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-tunnel-example"},
	}

	conn, err := quic.DialAddr(context.Background(), cfg.ServerAddress, tlsConf, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	log.Printf("Successfully established QUIC connection to server: %s", conn.RemoteAddr())

	if err := auth.PerformClientAuth(context.Background(), conn, cfg.Password); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	log.Println("Authentication successful. Starting listeners...")

	var wg sync.WaitGroup
	for _, tunnel := range cfg.Tunnels {
		wg.Add(1)
		go func(tunnel config.TunnelMapping) {
			defer wg.Done()
			listenAndProxy(conn, tunnel)
		}(tunnel)
	}

	// Wait for all listeners to finish (which is never, in this case)
	wg.Wait()
}

func listenAndProxy(conn quic.Connection, tunnel config.TunnelMapping) {
	listener, err := net.Listen("tcp", tunnel.LocalAddress)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", tunnel.LocalAddress, err)
		return
	}
	defer listener.Close()
	log.Printf("Listening on %s, forwarding to %s", tunnel.LocalAddress, tunnel.RemoteAddress)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection on %s: %v", tunnel.LocalAddress, err)
			continue
		}

		go handleLocalConnection(conn, localConn, tunnel.RemoteAddress)
	}
}

func handleLocalConnection(conn quic.Connection, localConn net.Conn, remoteAddr string) {
	log.Printf("Accepted connection from %s for target %s", localConn.RemoteAddr(), remoteAddr)

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("Failed to open stream for %s: %v", remoteAddr, err)
		localConn.Close()
		return
	}

	// Send the remote address to the server, followed by a newline character
	_, err = stream.Write([]byte(remoteAddr + "\n"))
	if err != nil {
		log.Printf("Failed to write remote address to stream: %v", err)
		localConn.Close()
		stream.Close()
		return
	}

	// Relay data
	proxy.Relay(localConn, stream)
	log.Printf("Connection from %s for target %s closed.", localConn.RemoteAddr(), remoteAddr)
}
