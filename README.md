<h1 align="center">Podium</h1>

<p align="center">
  <strong>A lightweight container orchestration tool with health checking and automatic recovery</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#why-podium">Why Podium</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#roadmap">Roadmap</a> •
  <a href="#contributing">Contributing</a> •
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

While tools like Kubernetes provide comprehensive container orchestration, they can be complex and resource intensive for simpler use cases. Podium aims to surpass Kubernetes in specific areas by focusing on simplicity, developer experience, and specialized use cases.

### "Kubernetes-Simple" Developer Experience
Podium eliminates the steep learning curve of Kubernetes with:
- Drastically simplified configuration (80% less verbose than Kubernetes)
- One command deployments that "just work" with sensible defaults
- Clear, intuitive CLI feedback showing exactly what's happening
- No need to understand complex concepts like pods, deployments, services, etc.

### Application-Level Health Monitoring
Beyond basic container health checks, Podium provides:
- Deep insights into application health, not just container status
- Business relevant metrics for truly understanding application performance
- Intelligent recovery actions tailored to specific failure scenarios
- Simple health dashboards for at a glance monitoring

### Single-Node Excellence
Podium is optimized for the common case of single server deployments:
- Significantly lower resource overhead than Kubernetes
- No complex networking required
- Simple backup and restore capabilities
- Perfect for small to medium applications that don't need multi-node clusters

Podium is ideal for:
- Development environments where Kubernetes is overkill
- Small production deployments with basic health checking needs
- Edge computing scenarios with limited resources
- Learning container management concepts without the complexity of larger systems
- Teams that want to focus on building applications, not managing infrastructure

## Installation

### Prerequisites

- Go 1.16 or higher
- Docker

### From Source

```bash
# Clone the repository
git clone https://github.com/odlemon/podium.git
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

### Short-Term Goals

- [ ] Simplified YAML configuration format
- [ ] One-command deployment CLI (`podium deploy app.yaml`)
- [ ] Enhanced application-level health monitoring
- [ ] Basic health dashboards
- [ ] Resource usage monitoring and limits
- [ ] Container logs streaming
- [ ] Volume management

### Medium-Term Goals

- [ ] Advanced health metrics with business-relevant indicators
- [ ] Intelligent recovery actions beyond simple restarts
- [ ] Network management and service discovery
- [ ] Simple backup and restore functionality
- [ ] Authentication and authorization
- [ ] Webhook notifications for container events

### Long-Term Vision 

- [ ] Complete "Kubernetes-Simple" developer experience
- [ ] Comprehensive application health insights
- [ ] Single-node excellence with optimized resource usage
- [ ] Support for Docker Compose-like configuration
- [ ] Multi-node support (optional, for growth scenarios)
- [ ] Metrics collection and visualization
- [ ] Edge computing optimizations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

