# Remote Overlay MVP

This repository contains a minimal skeleton for a control-plane and client agent
based on WireGuard/QUIC overlay networking. The structure is intentionally
simple to serve as a starting point.

## Directory Layout

- `db/` – SQL schema for control-plane storage.
- `proto/` – gRPC service definitions.
- `server/` – control plane server entry points.
- `client/` – client agent implementation.

## Building

```
go fmt ./...
go build ./...
```

## Testing

Run unit tests for the server and agent components:

```
go test ./...
```

## Running

Launch the control plane server:

```
go run server/cmd/control/main.go
```

Start the agent pointing at the server:

```
go run client/cmd/agent/main.go -addr localhost:50051 -token <token>
```

