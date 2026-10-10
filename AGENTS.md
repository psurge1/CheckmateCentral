# Distributed chess platform

## Project intent

This is an unfinished Go chess tournament platform originally implemented entirely
by the repository owner, without AI. Build on that work in small, understandable
increments. Explain changes and their tests; avoid wholesale rewrites or adding
frameworks merely to prepare for hypothetical needs.

The immediate product goal is for users to:

- Sign up, sign in, and use their account.
- Create and join tournaments, then play games in them.
- Leave a tournament before it starts.
- Pause participation or withdraw from remaining games during a tournament.

The immediate implementation priority is PostgreSQL-backed CRUD: account
storage, signup/login, and tournament creation, reading, updating, and deletion.
Establish that functionality before integrating game execution. PostgreSQL is
the chosen database, not a technology decision to reopen in the build plan.

Tournaments are **arenas**, not rounds: eligible players are repeatedly paired
until the tournament time expires. Joining adds a player to that tournament's
matchmaking pool. The pool produces batches of games, which go to a scheduler/job
queue. Eventually, the scheduler assigns each batch to a worker with light load.

The long-term target is a distributed deployment on Kubernetes using Helm.
Workers will execute game jobs (including authoritative move processing), admin
jobs (account creation, signup, login, tournament creation/joining, and related
actions), and future functionality. WebSockets are the current transport; gRPC
is the planned eventual migration. Kubernetes/cluster configuration YAML and
Helm charts must live in this repository as infrastructure as code.

## Current implementation

These are observations of the current code, not claims that the target exists:

- `game_server.go` constructs and mounts user, session, tournament, and game
  muxes, then starts HTTP on port 8080. `internal/httpapi` separates each group's
  mux registration from its handler file. Planned API handlers return 501 stubs.
  The existing `/joingame` handler lives in `games_handlers.go`; it still pairs
  connections directly using process-local slices and launches game goroutines.
- `matchmaking/fcfs_pool.go` contains the basic FIFO pool, protected by a mutex.
- `matchmaking/matchmaker.go` already implements `TryMatch`: it extracts games,
  converts them to `[]Job`, constructs a `Batch`, and calls the shared scheduler.
  Its interface requires `ExtractMatches`, while the FIFO pool currently exposes
  `ExtractMatch`; this integration is unfinished.
- `taskscheduler/taskscheduler.go` runs jobs locally with a `sync.WaitGroup`.
  Its existing `Job`, `Batch`, `TaskScheduler`, and placeholder `Server` types
  are the foundation to extend. `Submit` starts goroutines immediately; pending
  queueing, durability, and load-based remote assignment are not implemented.
- `models/game.go` and `streaming/streaming.go` provide basic game objects and
  WebSocket move handling using `github.com/notnil/chess`.
- Accounts, sessions, tournament lifecycle, standings, and persistence have not
  been implemented. A player currently has a username and Elo field.
- `datastructures/queue_test.go` contains an empty test stub. Passing the current
  test command does not demonstrate the intended behavior.
- `cluster-config.yaml` already describes a local kind cluster with one control
  plane and two worker nodes. Helm charts are still to be built.

## Working rules

- Follow [the build plan](docs/BUILD_PLAN.md), one stage and one small build step
  at a time. Implement only the requested scope, and preserve unrelated edits.
- Use [the agreed REST endpoint spec](docs/API_ENDPOINTS.md) for HTTP API work.
  Implement handlers with Go's standard `net/http` package, without another
  framework.
- Keep the existing Go packages useful where practical. Introduce boundaries
  when a stage needs them, and explain why each new component is necessary.
- Ground each step in specific existing types, functions, and files. Distinguish
  implemented behavior, missing integration, and new functionality. Preserve
  `Player`/`GPlayer`, `Game`, the FIFO pool, `Matchmaker`, `Job`, `Batch`, and
  `TaskScheduler` as the starting points rather than planning parallel replacements.
- Start with account and tournament CRUD on PostgreSQL. Extend `Player` with
  persistent identity and add tournament records; wire those records into the
  existing pool/batch/scheduler flow afterward. Keep credentials separate from
  public player data and transport connections.
- Preserve the intended ownership: one pool per matchmaker, many matchmakers
  sharing a scheduler. Put tournament orchestration outside `models` to avoid
  an import cycle with matchmaking, which already imports models.
- Each stage must have observable acceptance checks and a runnable demo or test.
  Finish its checks before treating it as complete.
- Keep FIFO matchmaking initially. Arena scoring and participation rules need
  explicit definitions; proposed defaults in the build plan are not established
  product requirements.
- Use stable IDs for accounts, tournaments, games, batches, and jobs. Do not use
  usernames, pointers, or WebSocket connections as distributed identity.
- Separate domain/game state from transport connections. Workers must process
  serializable requests; a WebSocket connection cannot be sent through a queue.
- Maintain one authoritative owner for a game and one matchmaking authority per
  tournament. Adding replicas must not create duplicate games or concurrent
  writers. Remote ownership requires leases and fencing before failover.
- Plan for retries and duplicate delivery. A completed job/result must not apply
  twice, and a failed queue submission must not strand matched players.
- Keep game jobs and admin jobs from starving each other. Do not block account
  and tournament operations behind long-running games.
- Keep infrastructure, deployment instructions, and Helm configuration in this
  repo. Commit configuration and secret references, never real credentials.
- Keep WebSockets working during early stages. Specify the gRPC migration scope
  before changing clients; a browser client needs a compatible edge transport.

## Verification

The module currently declares Go 1.26.5. Use a compatible toolchain.

```sh
go test ./...
go test -race ./...
go vet ./...
```

Use focused tests during implementation and the full checks at stage completion.
Test clocks/deadlines with an injected clock instead of real-time sleeps. For
distributed stages, add integration checks for duplicate delivery, worker loss,
reconnection, and rolling restarts. Keep those checks separate from fast tests
and document their prerequisites.

The current manual server entry point is `go run game_server.go`. Two WebSocket
clients can connect to `ws://localhost:8080/joingame` for a manual game demo;
this is not automated verification.
