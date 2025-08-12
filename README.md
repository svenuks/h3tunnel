# h3tunnel - Secure QUIC Tunnel

`h3tunnel` is a client-server (C/S) application designed to securely and transparently tunnel TCP traffic over the QUIC protocol.

## 🚀 Features

*   **Secure & Transparent TCP Tunneling over QUIC**: Encrypts and forwards TCP traffic without modifying application data.
*   **Client-Server Architecture with Password Authentication**: Flexible deployment with secure access control.
*   **Flexible Multi-Tunnel Configuration**: Define multiple local-to-remote port mappings.
*   **Cross-Platform**: Supports Linux, macOS (Intel/Apple Silicon), and Windows builds.
*   **Golang Implementation**: High-performance and easy to deploy.
*   **Debian Package Support**: Streamlined deployment on Debian/Ubuntu systems.

## 🛠️ Build & Usage

### Prerequisites

*   Go (version 1.24.6 or higher)
*   `git` (for cloning the repository)
*   `dpkg-dev` and `debhelper` (only for building `.deb` packages)

### Build from Source

1.  **Clone the repository** (assuming you've uploaded it to GitHub):
    ```bash
    git clone https://github.com/your-username/h3tunnel.git
    cd h3tunnel
    ```
2.  **Build executables**:
    ```bash
    # Build Linux (amd64) version
    GOOS=linux GOARCH=amd64 go build -o bin/h3tunnel-client-linux-amd64 ./cmd/client
    GOOS=linux GOARCH=amd64 go build -o bin/h3tunnel-server-linux-amd64 ./cmd/server

    # Build macOS (amd64) version
    GOOS=darwin GOARCH=amd64 go build -o bin/h3tunnel-client-darwin-amd64 ./cmd/client
    GOOS=darwin GOARCH=amd64 go build -o bin/h3tunnel-server-darwin-amd64 ./cmd/server

    # Build macOS (arm64) version
    GOOS=darwin GOARCH=arm64 go build -o bin/h3tunnel-client-darwin-arm64 ./cmd/client
    GOOS=darwin GOARCH=arm64 go build -o bin/h3tunnel-server-darwin-arm64 ./cmd/server

    # Build Windows (amd64) version
    GOOS=windows GOARCH=amd64 go build -o bin/h3tunnel-client-windows-amd64.exe ./cmd/client
    GOOS=windows GOARCH=amd64 go build -o bin/h3tunnel-server-windows-amd64.exe ./cmd/server
    ```
    Built executables will be in the `bin/` directory.

### Configuration

Configure the application using `config.client.json` and `config.server.json`.

*   **`config.server.json`**:
    ```json
    {
      "listen_address": "0.0.0.0:4430",
      "password": "your-secret-password"
    }
    ```
    `listen_address`: Address and port for the server to listen for QUIC connections.
    `password`: Password for client authentication.

*   **`config.client.json`**:
    ```json
    {
      "server_address": "remote-server-ip:4433",
      "password": "your-secret-password",
      "tunnels": [
        {
          "local_address": "127.0.0.1:1080",
          "remote_address": "127.0.0.1:1080"
        }
      ]
    }
    ```
    `server_address`: Address and port of the proxy server.
    `password`: Password for authenticating with the server.
    `tunnels`: An array defining multiple tunnel mappings.
        `local_address`: Local TCP address and port for the client to listen on.
        `remote_address`: The final remote target address and port for traffic forwarding.

### Running the Application

1.  **Start the Server**:
    ```bash
    ./bin/h3tunnel-server # or other platform executable
    ```
2.  **Start the Client**:
    ```bash
    ./bin/h3tunnel-client # or other platform executable
    ```
## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).
