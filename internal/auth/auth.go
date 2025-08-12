package auth

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/quic-go/quic-go"
)

const (
	authStreamTimeout = 5 * time.Second
	AuthSuccess       = 0x00
	AuthFailure       = 0x01
)

// PerformClientAuth performs the client-side authentication handshake.
func PerformClientAuth(ctx context.Context, conn quic.Connection, password string) error {
	ctx, cancel := context.WithTimeout(ctx, authStreamTimeout)
	defer cancel()

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to open auth stream: %w", err)
	}

	// Use a goroutine to close the stream to avoid blocking on writing the password
	// in case the server is not reading.
	go func() {
		<-ctx.Done()
		stream.Close()
	}()

	_, err = stream.Write([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to write password: %w", err)
	}

	// Signal that we are done writing.
	// The actual close will happen when the context is done.
	stream.Close()

	resp := make([]byte, 1)
	_, err = io.ReadFull(stream, resp)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp[0] != AuthSuccess {
		return fmt.Errorf("authentication failed from server")
	}

	return nil
}

// HandleServerAuth performs the server-side authentication handshake.
func HandleServerAuth(ctx context.Context, conn quic.Connection, expectedPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, authStreamTimeout)
	defer cancel()

	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to accept auth stream: %w", err)
	}
	defer stream.Close()

	passBytes, err := io.ReadAll(stream)
	if err != nil {
		return fmt.Errorf("failed to read password from stream: %w", err)
	}

	if string(passBytes) != expectedPassword {
		_, _ = stream.Write([]byte{AuthFailure})
		return fmt.Errorf("invalid password received from %s", conn.RemoteAddr())
	}

	_, err = stream.Write([]byte{AuthSuccess})
	if err != nil {
		return fmt.Errorf("failed to write auth success response: %w", err)
	}

	return nil
}