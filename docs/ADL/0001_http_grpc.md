# 1. Keep separate grpc and http

# Status

Accepted

## Context

To have replication and load balancing

## Decision

We replicate and distribute data in services.
Http is standalone service, that connects to several grpc.
Each grpc service has it's own pg database, that is replicated
via raft protocol

## Consequences

We cannot load in one executable both http and grpc
