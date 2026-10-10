# REST API endpoint spec

Agreed endpoint list for the chess tournament platform. Implement using Go's
standard `net/http` package and `http.ServeMux`, without another framework.
These are planned endpoints, not a list of implemented handlers.

## Users and sessions

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/users` | Register |
| `GET` | `/users` | List players, with pagination |
| `GET` | `/users/me` | Get signed-in account |
| `GET` | `/users/{userID}` | Get a specific player's public profile |
| `PATCH` | `/users/me` | Edit account |
| `POST` | `/sessions` | Login |
| `DELETE` | `/sessions/current` | Logout |

Player listings and public profiles return public fields such as ID, username,
and Elo. Keep credentials and session information private.

## Tournaments

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/tournaments` | List tournaments |
| `GET` | `/tournaments/{tournamentID}` | Get tournament |
| `POST` | `/tournaments` | Create tournament |
| `PATCH` | `/tournaments/{tournamentID}` | Edit tournament |
| `DELETE` | `/tournaments/{tournamentID}` | Delete tournament |
| `PUT` | `/tournaments/{tournamentID}/participants/me` | Join |
| `DELETE` | `/tournaments/{tournamentID}/participants/me` | Leave before start |
| `PATCH` | `/tournaments/{tournamentID}/participants/me` | Pause, resume, or withdraw |
| `GET` | `/tournaments/{tournamentID}/participants` | List participants |
| `GET` | `/tournaments/{tournamentID}/standings` | Get scores/rankings |

## Games

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/games` | List games; filter by tournament/player |
| `GET` | `/games/{gameID}` | Get position, players, status, result |
| `GET` | `/games/{gameID}/moves` | Get move history |
| `POST` | `/games/{gameID}/moves` | Submit move, if using HTTP |
| `POST` | `/games/{gameID}/resign` | Resign |

Move submission remains over WebSockets initially; the HTTP move endpoint is
optional. A proposed game-specific WebSocket endpoint is
`GET /games/{gameID}/stream`, which upgrades the connection rather than returning
a REST response.

Games are created internally by the matchmaker and scheduler, so there is no
public `POST /games`. Account and participation operations derive the acting
player from their session.
