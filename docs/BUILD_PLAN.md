# Build plan: extend the existing server

Start with PostgreSQL-backed account and tournament CRUD. Then connect those
records to the matchmaking, batching, and game code already in this repository.
Distributed workers and deployment come after the local user flow works.

This plan targets the current `main` branch. All stages are pending. Implement
one numbered step at a time; the file names suggested for new code are proposals,
not instructions to restructure the existing packages.

## What is already here

| Existing code | What it already does | What remains to connect or extend |
| --- | --- | --- |
| `models.Player` in `models/player.go` | Holds username and Elo | Persistent ID, account storage, authenticated identity |
| `models.GPlayer` in the same file | Connects a player, socket, and game | Populate `Player` from login; enforce `Game` ownership |
| `models.Game`, status/result enums, `NewGame`, `Run` | Describes a paired game and runs it | IDs, tournament association, ownership, completion bookkeeping |
| `FirstComeFirstServeMPool` in `matchmaking/fcfs_pool.go` | Mutex-protected FIFO join/leave and single-game extraction | Implement the interface's batch extraction; membership eligibility |
| `Matchmaker.TryMatch` in `matchmaking/matchmaker.go` | Converts extracted games into `[]Job`, builds a `Batch`, calls `Submit` | Constructor, handler wiring, visible submission failure |
| `Job`, `Batch`, `TaskScheduler.Submit/Close` in `taskscheduler/taskscheduler.go` | Executes local jobs concurrently and waits for them | Actual pending queue, bounded execution, eventual worker assignment |
| `Server` and `TaskScheduler.servers` in the same file | Placeholders for remote execution | Worker identity, load, capacity, liveness |
| `streaming.EstablishGameConnection` | Uses the chess library to validate and forward WebSocket moves; returns an outcome | Reliable termination, connection lifecycle, returning results into `Game` |
| `game_server.go` | HTTP server and a direct two-socket game demo | Account/tournament endpoints; later route games through the existing matchmaker |
| `datastructures.Queue` | FIFO storage, peek, length, deletion | Useful behavior tests; retain it for the local pool/queue |
| `cluster-config.yaml` | Existing three-node kind configuration | Application deployment and Helm chart |

The intended game execution path already exists in pieces:

```text
HTTP/WebSocket handler
  → GPlayer → FirstComeFirstServeMPool
  → Matchmaker.TryMatch → Batch{Jobs: []Job}
  → TaskScheduler.Submit → Game.Run
  → streaming.EstablishGameConnection
```

Today the handler bypasses that path and directly calls `go game.Run()`. Also,
`MatchmakingPool` requires `ExtractMatches`, but the FIFO implementation exposes
only `ExtractMatch`. Those are later integration tasks, not reasons to rebuild
matchmaking or batching before CRUD.

## Stage 1 — Connect PostgreSQL and persist accounts

**Build on:** `models.Player` and the existing Go module/server.
**Add:** a small PostgreSQL storage package and SQL migrations.

1. Document local PostgreSQL startup and connection configuration, and add a
   connection/ping/close path. Keep credentials in environment configuration.
   Use Go's database API with a PostgreSQL driver; choose the driver when coding.
2. Add an account ID to `Player`. Define account and credential storage without
   putting password hashes in client responses or `GPlayer`. Preserve the existing
   username/Elo fields; decide whether credentials need a separate account type.
3. Add an accounts migration and storage operations to create, read, update, and
   delete/deactivate accounts. Enforce unique usernames in PostgreSQL. Define
   deletion behavior before later tournament/game records reference accounts.

**Checks:** apply migrations to an empty test database; create/read/update/remove
an account; duplicate usernames fail even for concurrent inserts; restart the
server/storage connection and read the same account. Use a real PostgreSQL test
instance rather than relying entirely on mocked SQL.

**Checkpoint:** account records survive restart. No matchmaking changes yet.

## Stage 2 — Add signup, login, and account endpoints

**Build on:** `game_server.go` handlers and stage 1 account storage.
**Add:** authentication/session helpers; keep them separate from chess logic.

1. Add signup and login handlers. Hash passwords with an established library;
   verify hashes at login. Signup creates the existing player identity through
   the account store, rather than inventing another matchmaking player model.
2. Add expiring sessions, logout, and an authentication helper for handlers.
   Persist sessions or explicitly document restart invalidation. Add current-user
   read/update/delete operations; users cannot edit Elo or another account.
3. Authenticate `/joingame` before upgrading and construct
   `GPlayer{Player: authenticatedPlayer, Moves: conn}`. Keep the current direct
   game demo working until stage 5 connects the matchmaker.

**Checks:** signup → login → read own account; wrong passwords and expired
sessions fail; logout/deletion revokes access; duplicate signup is rejected;
one account cannot modify another; unauthenticated WebSocket upgrades fail.

**Checkpoint:** a signed-in user has a persisted `Player`, not just a socket.

## Stage 3 — Add tournament CRUD

**Build on:** the account IDs, authentication, and existing HTTP server.
**Add:** `models/tournament.go`, tournament storage/migration, and handlers.

1. Define a tournament ID, creator account ID, name, start/end times, and status.
   This fills in the tournament concept already anticipated by the commented
   `TournamentID` field in `models.Game`. Keep the model as data.
2. Add create/list/read/update/delete storage operations and endpoints. Authenticate
   writes; only the creator can update/delete their tournament. Validate time
   ranges and missing IDs. Make database transactions responsible for consistency.
3. Define deletion/edit rules. Proposed initial rule: scheduled tournaments can
   be edited/deleted; started tournaments cannot, preserving game history.
   Translate these failures into clear API responses.

**Checks:** signed-in user creates, fetches, edits, and deletes a scheduled
arena; another user cannot edit/delete it; invalid schedules fail; restart
retains tournaments; deleted IDs return not-found; active tournaments reject
edits/deletion under the chosen rule.

**Checkpoint:** PostgreSQL-backed account and tournament CRUD works through HTTP.
This is the first product milestone. Test it before touching batch execution.

## Stage 4 — Persist joining, leaving, pausing, and withdrawal

**Build on:** tournament CRUD and `FirstComeFirstServeMPool.Join/Leave`.
**Add:** persisted tournament memberships and a small tournament coordinator.

1. Add membership records keyed uniquely by tournament ID and account ID. Add
   authenticated join/leave/pause/resume/withdraw endpoints. Keep participation
   state distinct from `GPlayer.Game` so pausing does not erase an ongoing game.
2. Let the coordinator own one existing FIFO pool and, once wired in stage 5,
   one existing matchmaker per tournament. Keep this orchestration outside
   `models`: importing matchmaking there would create an import cycle through
   matchmaking's existing dependency on models.
3. Apply the participation rules below and rebuild eligibility from stored
   memberships on restart/reconnect. Registration can exist without a socket;
   only connected, eligible players should enter executable matchmaking.
   Remove queued participants through the existing pool's `Leave` operation.

**Checks:** duplicate joins create one membership; leave works before start;
pause/resume and withdrawal retain history; a user cannot alter another user's
membership; restarting retains participation state; deleting a scheduled arena
removes its memberships; account deletion cannot strand active participation.

**Checkpoint:** the participation API works, even before it launches games.

## Stage 5 — Wire the existing pool → batch → scheduler path

**Build on:** `ExtractMatch`, `MatchmakingPool`, `TryMatch`, `Batch`, `Submit`,
`NewGame`, and the tournament coordinator. Extend these, rather than replacing them.

1. Implement `ExtractMatches` on the FIFO pool using its existing extraction
   logic, with a batch limit and odd-player retention. Avoid nested locking when
   sharing logic with `ExtractMatch`. Use `NewGame` instead of the current direct
   struct literal, validate nil/self-pairing, and assign both `GPlayer.Game`
   references under a defined ownership lock. Reject duplicate pool joins.
2. Add a `NewMatchmaker` constructor for its private pool/scheduler fields.
   Share one local `TaskScheduler` between the tournament matchmakers, exactly
   as the current ownership comments intend. Make `TryMatch` report errors and
   recover reservations if `Submit` fails; do not silently exhaust retries.
3. Route tournament-ready players from the handler/coordinator into their pool,
   then invoke `TryMatch` on a documented trigger/tick. Add game/tournament IDs.
   Retain `Batch.Jobs` and `Job.Run`; validate `Count` against the slice length.
   Replace direct pairing in the handler once this path passes its tests.

**Checks:** add a compile-time assertion that the FIFO pool implements
`MatchmakingPool`; five players extract two games and leave one waiting;
`TryMatch` submits those games in one batch; fake jobs verify `Submit` invokes
each once; a rejected submission restores eligible players; different arenas
never mix; concurrent joins cannot pair the same player twice. Fill the empty
queue test with FIFO/deletion/empty-operation cases supporting these behaviors.

**Checkpoint:** the initialized matchmaking and batch code now runs real games.

## Stage 6 — Close the arena game loop

**Build on:** `Game.Run`, existing status/result enums, `GPlayer.Game`, and
`streaming.EstablishGameConnection`'s existing return value.

1. Capture the outcome currently discarded by `Game.Run`; move status from
   paired to playing to finished, translate the chess outcome into `GameResult`,
   and notify the tournament coordinator once. Fix streaming termination after
   either player's terminal move and invalid-input/read-error paths. Keep the
   existing chess library for move legality.
2. Release both game references, persist the result, update arena scores once,
   and requeue ready participants through the same FIFO pool. The streaming
   function currently closes both sockets per game: define a reconnect-to-ready
   flow initially or retain tournament connections with one clear socket owner.
3. Add arena start/deadline handling with an injected clock. Stop new games at
   expiry and honor pause/withdraw while a game is finishing. Add a game clock
   and a defined disconnect outcome so an arena cannot wait forever for a game.

**Checks:** game termination by either color; illegal moves and disconnects;
duplicate completion cannot double-score or double-requeue; a completed game
allows another pairing; pause/withdraw stops subsequent games; the arena
starts/stops at its configured times. Test actual WebSocket behavior as well
as model transitions.

**Checkpoint:** run signup/login → create/join → repeated games → pause/resume
→ withdraw → deadline → standings through public endpoints. Also demonstrate
leave-before-start. A scripted client is sufficient; a new UI is optional.

## Stage 7 — Extend the scheduler with queued game and admin work

**Build on:** `Job.Run`, `Batch`, `TaskScheduler.Submit/Close`, and
`datastructures.Queue`. `Submit` currently starts goroutines immediately; it is
not yet a pending job queue.

1. Add a bounded, synchronized pending-batch queue and local worker loops to
   `TaskScheduler`. Keep `Submit(*Batch)` as the entry point; define acceptance,
   full-queue errors, completion, and shutdown. Start with one batch per worker.
2. Wrap the account/tournament service operations from stages 2–4 as admin jobs
   implementing the same local `Job` contract. Handlers submit work and receive
   a bounded result/error response; success requires a committed operation.
   Reuse the service logic rather than duplicate SQL in workers.
3. Give game/admin work separate capacity so long games do not block login or
   joining. Keep `Game.Run` as a game job; add error/result reporting only where
   the current interface needs it. Retries must not repeat committed effects.

**Checks:** accepted batches wait when capacity is full; jobs execute once;
full/closed queues return errors; shutdown drains or explicitly cancels work;
admin jobs return operation results; a busy game batch does not starve login.

**Checkpoint:** existing batches execute via a real local scheduler queue.

## Stage 8 — Make those batches durable and assign remote workers

**Build on:** the queued scheduler, `Server`, `TaskScheduler.servers`, and the
same game/admin handlers. Complete this stage's steps separately.

1. Add stable job/batch IDs and serializable payloads beside the local `Job`
   interface. A worker decodes a payload into the existing execution logic;
   do not attempt to serialize Go pointers, methods, or WebSocket connections.
   Persist pending work, claims, and results in PostgreSQL initially, with atomic
   state-to-job handoff, retries, and idempotency keys.
2. Extend `Server` with worker identity, supported work, active/reserved load,
   capacity, and heartbeat. Replace local-only dispatch with a worker adapter;
   retain the local adapter for tests. Assign a batch to the healthy compatible
   worker with the lowest load that has room, reserving its capacity immediately.
3. Keep client WebSockets at the server/gateway; route moves to the worker owning
   the game and events back. Separate chess state from streaming I/O at this
   boundary. Persist recoverable game state and clocks, add ownership leases
   and stale-owner rejection, and retry only unfinished games in a batch.

**Checks:** queued work survives restart; duplicate admin/result jobs apply once;
two workers receive batches according to available capacity; concurrent dispatch
cannot overfill a worker; dead workers receive no assignments; killing a worker
allows game recovery/reconnect without two authoritative writers. Keep one
scheduler/matchmaker authority initially; replicas require explicit ownership.

**Checkpoint:** one server, one scheduler, and two local worker processes.

## Stage 9 — Put this deployment into Kubernetes and Helm

**Build on:** the existing `cluster-config.yaml` and stage 8 process roles.

1. Add container builds and health/readiness/drain behavior for server, scheduler,
   and workers. Supply PostgreSQL connection and worker configuration externally.
2. Add `deploy/helm/chess-platform/` with chart templates and local values for
   Deployments, Services, resource settings, secret references, and persistent
   PostgreSQL (or a documented external instance). Reuse the existing kind YAML;
   move it only if its documented commands are updated too.
3. Add `deploy/README.md` with clean-checkout build/load/install/upgrade/cleanup
   commands. Keep coordinator replicas at one until stage 8 ownership permits
   safe scaling; scale the workers to demonstrate distributed execution.

**Checks:** `helm lint`, `helm template`, install in kind, run the stage 6 user
flow, scale workers, restart pods, and perform a rolling upgrade. Accounts,
tournaments, queued work, and results must survive restarts without duplication.

**Checkpoint:** Kubernetes YAML, Helm charts, and deployment instructions live
in this repo and reproduce the working platform.

## Stage 10 — Migrate the worker transport to gRPC

**Build on:** stage 8's worker adapter and job/game command boundaries.

1. Define protobuf contracts for batch assignment, admin results, game commands/
   events, and worker status. Commit schemas and generation instructions.
2. Implement gRPC behind the adapter with deadlines, service authentication,
   cancellation, idempotent retries, and stream reconnection. Keep scheduler,
   pool, and tournament logic independent of the transport change.
3. Rerun the same user-flow and worker-failure checks through gRPC and Helm.
   Decide browser transport separately: existing WebSockets can remain at the
   edge, while a full replacement needs a compatible browser client/bridge.

**Checks:** RPC/stream interruptions, duplicate requests, mixed compatible
versions during upgrades, and recovery preserve the behavior already tested.

**Checkpoint:** remote execution uses gRPC without rewriting arena logic.

## Participation rules to settle before stages 3–6

The user has specified arena matchmaking, leaving before start, and pausing or
withdrawing from remaining games. The following details are proposed defaults:

- Joining before start registers membership; active, connected participants enter
  the runnable pool. Late joining is allowed while the arena is active.
- Pause is reversible; withdrawal is final for that arena. Both prevent future
  games while an already-started game continues. Cancel pending pairings safely.
- A player can belong to multiple arenas but play only one game at a time.
- At expiry, no new game starts. Started games finish, then standings finalize.
- Initial scoring is win = 2, draw = 1, loss = 0, with shared ranks for ties.
- Scheduled arenas can be edited/deleted; started ones retain their history.

Choose the game time control and disconnect policy in stage 6. Record any rule
changes here before implementing them.

## Verification and next step

Use focused package tests for each step, then `go test ./...`,
`go test -race ./...`, and `go vet ./...` at local stage completion. Database
checks need a disposable PostgreSQL instance and documented migration/reset
commands; cluster checks need kind and Helm. Keep test credentials out of Git.

The baseline `go test ./...` passes, but the existing queue test is empty.
Passing today is not evidence that the full pipeline works. Add meaningful tests
alongside each implementation change, rather than treating test scaffolding as
a separate product milestone.

**Next step: stage 1, step 1 — connect to a local PostgreSQL instance and establish
configuration, migrations, and a database connectivity check.**
