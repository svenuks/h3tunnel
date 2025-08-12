package proxy

import (
	"io"
	"net"

	"github.com/quic-go/quic-go"
)

// Relay copies data bidirectionally between a local TCP connection and a remote QUIC stream.
func Relay(local net.Conn, remote quic.Stream) {
	defer local.Close()
	defer remote.Close()

	// Goroutine to copy data from local to remote
	go func() {
		_, err := io.Copy(remote, local)
		if err != nil {
			// This error is expected if the other side closes the connection, so we don't log it as a fatal error.
			// log.Printf("Error copying from local to remote: %v", err)
		}
	}()

	// Copy data from remote to local in the main goroutine
	_, err := io.Copy(local, remote)
	if err != nil {
		// log.Printf("Error copying from remote to local: %v", err)
	}
}
