# mini-swarm-dash

A lightweight real-time monitoring dashboard for Docker Swarm clusters. Displays CPU, memory, and temperature per node along with per-container metrics and service status.

![Dashboard showing node cards with resource bars and a services table]

## Features

- **Node monitoring** — CPU usage, memory consumption, system temperature, role badges (leader/manager/worker), and status indicators
- **Container metrics** — per-container CPU % and memory usage, associated service names, expandable per-node
- **Service overview** — deployment mode (replicated/global), replica counts, image names
- **Real-time updates** — 5-second polling with pause/resume control
- **Color-coded resource bars** — green → yellow (>70%) → red (>90%)

## Architecture

The server runs on every Swarm node. The node you open the dashboard on acts as the aggregator:

1. It queries the Docker Swarm API for the list of nodes.
2. It fans out HTTP requests in parallel to each node's `/api/local-stats` endpoint.
3. Each node reads its own container stats from the Docker socket and its temperature from `/sys/class/thermal`.
4. The aggregator merges everything and serves it to the frontend.

```
Browser → GET /api/nodes (aggregator node)
               └─► GET /api/local-stats  (node-1)
               └─► GET /api/local-stats  (node-2)
               └─► GET /api/local-stats  (node-N)
```

## Stack

| Layer    | Technology                     |
|----------|-------------------------------|
| Backend  | Go 1.24, Docker SDK (Moby)    |
| Frontend | Svelte 5, TypeScript, Vite    |
| Build    | Multi-stage Docker (Alpine)   |

## Quick Start

### Docker (recommended)

Build the image and deploy it as a **global service** so every Swarm node runs it:

```bash
docker build -t mini-swarm-dash .

docker service create \
  --name mini-swarm-dash \
  --mode global \
  --publish published=8080,target=8080,protocol=tcp,mode=host \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  mini-swarm-dash
```

Open `http://<any-swarm-node>:8080` in your browser.

### Local development

**Frontend:**
```bash
cd client
pnpm install
pnpm build
```

**Backend:**
```bash
cd server
go run -tags=mock .
# Listens on :8080, serves mock data and frontend from client/dist instead of public
```

## Requirements (production)

- Docker Swarm cluster (single-node works too)
- `/var/run/docker.sock` mounted into the container
- Port `8080` open between Swarm nodes (for peer stat collection)
- Linux thermal zone sysfs (`/sys/class/thermal/`) for temperature readings (optional — skipped if absent)
