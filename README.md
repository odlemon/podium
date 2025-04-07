<h1 align="center">Podium</h1>

<p align="center">
  <strong>A lightweight container management system with health checking and automatic recovery</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#why-podium">Why Podium</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#roadmap">Roadmap</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

## Features

Podium is a container management system built in Go that provides:

- **Container Lifecycle Management**: Create, start, stop, and remove containers
- **Health Checking**: Automated health monitoring of containers
- **Automatic Recovery**: Restart policies for unhealthy containers
- **REST API**: Simple HTTP API for container operations
- **Persistent Storage**: Container configurations stored in BoltDB
- **Lightweight**: Minimal resource footprint compared to full orchestration systems

## Why Podium

While tools like Kubernetes provide comprehensive container orchestration, they can be complex and resource-intensive for simpler use cases. Podium fills the gap between running containers manually with Docker and deploying a full Kubernetes cluster.

Podium is ideal for:

- Development environments where Kubernetes is overkill
- Small production deployments with basic health checking needs
- Edge computing scenarios with limited resources
- Learning container management concepts without the complexity of larger systems

## Installation

### Prerequisites

- Go 1.16 or higher
- Docker

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/podium.git
cd podium

# Build the binary
go build -o podium cmd/server/main.go

# Run the server
./podium
```

### Using Go Install

```bash
go install github.com/yourusername/podium/cmd/server@latest
```

### Platform-Specific Notes

The installation process is similar across platforms (Windows, macOS, Linux), but there are a few differences to note:

**Linux**:
- Docker socket is typically at: `unix:///var/run/docker.sock`

**macOS**:
- Docker Desktop for Mac uses the same socket path as Linux: `unix:///var/run/docker.sock`

**Windows**:
- If using Docker Desktop with WSL2 backend, use the Linux socket path
- If using Docker Desktop with Hyper-V backend, use: `npipe:////./pipe/docker_engine`
- You may need to adjust the `--docker-host` flag accordingly

## Usage

### Starting the Server

```bash
podium --port 8080 --db-path ./podium.db
```

### API Examples

#### Create a Container

```bash
curl -X POST http://localhost:8080/api/containers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web-server",
    "image": "nginx:latest",
    "ports": [{"host": 8080, "container": 80}],
    "healthCheck": {
      "type": "http",
      "endpoint": "/",
      "interval": "10s",
      "timeout": "2s",
      "retries": 3
    },
    "restartPolicy": {
      "type": "always"
    }
  }'
```

#### List Containers

```bash
curl http://localhost:8080/api/containers
```

#### Get Container Health

```bash
curl http://localhost:8080/api/containers/web-server/health
```

## Configuration

Podium can be configured using command-line flags or environment variables:

| Flag | Environment Variable | Description | Default |
|------|---------------------|-------------|---------|
| `--port` | `PODIUM_PORT` | HTTP server port | 8080 |
| `--db-path` | `PODIUM_DB_PATH` | Path to BoltDB file | ./podium.db |
| `--docker-host` | `PODIUM_DOCKER_HOST` | Docker host address | unix:///var/run/docker.sock |
| `--log-level` | `PODIUM_LOG_LEVEL` | Logging level (debug, info, warn, error) | info |

## Roadmap

### Current Focus

- ✅ Basic container management (create, start, stop, remove)
- ✅ Health checking system
- ✅ Restart policies for unhealthy containers
- ✅ BoltDB storage for container configurations

### Upcoming Features

- [ ] Volume management
- [ ] Network management
- [ ] Container resource limits
- [ ] Container logs streaming
- [ ] Authentication and authorization
- [ ] Multi-node support
- [ ] Metrics collection and visualization
- [ ] Webhook notifications for container events
- [ ] Support for Docker Compose-like configuration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
