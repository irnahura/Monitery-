# AI_LOG.md

Built this with AI doing most of the typing, me doing most of the deciding. Kept prompts below mostly verbatim — cleaned up typos, nothing else.

## Tech Stack

- Go + Gin + GORM (backend)
- PostgreSQL
- React + Vite + TypeScript (frontend)
- Docker Compose
- ChatGPT (GPT-5.5) — planning/architecture side
- OpenAI Codex — actual implementation

Started from Peekaping as a reference point rather than from scratch — closest thing I found to what this needed without dragging in enterprise baggage.

---

## Getting from spec to something buildable

```
Read the assignment specification and build a complete uptime monitoring MVP with:
- Go backend with Gin framework, GORM, PostgreSQL
- React frontend with Vite and TypeScript
- Docker Compose orchestration
- JWT authentication
- Monitor CRUD operations
- Scheduled health checks
- Email notifications
- API key management
Follow the Open Design wireframes in the design/ folder
```

Before this, I'd asked it separately to just strip the assignment down to an MVP and tell me what's actually load-bearing vs. nice-to-have. Backend, frontend, scheduler, persistence, docker, docs — everything else waited until it was actually needed.

Also had it go compare a few open-source uptime monitors first (readability and package structure over stars) — that's where Peekaping came from.

---

## Backend

```
Create a Go backend API with:
- Gin web framework
- GORM for PostgreSQL ORM
- JWT authentication with login/register endpoints
- Monitor CRUD endpoints (create, read, update, delete)
- Health check scheduler that pings URLs every minute
- Store HTTP status codes, response times, and timestamps
- API key authentication support
- Rate limiting middleware
Structure it in internal/ packages: api, auth, monitor, scheduler, database, models, repository
```

Got JWT auth, monitor CRUD, the scheduler, health history, availability/SLA calc, email notifications, and API keys out of this and the follow-ups. Handler → service → repository → db, established early so nothing later turned into a cross-package mess.

Most of the back-and-forth after the initial generation wasn't new features, it was pulling coupling apart before it calcified.

---

## Frontend

```
Build a React + Vite + TypeScript frontend that:
- Has login/register screens
- Dashboard showing all monitors with their status (up/down) and response times
- Add/edit monitor forms
- Service detail page showing health history
- Settings page for API key management
- Use the design wireframes from design/ folder as reference
- Make API calls to the Go backend at http://localhost:8080
- Handle JWT token storage and authentication
```

Told it to treat the wireframes as the actual design system, not a vibe to riff on. Stayed deliberately shallow on visuals — effort went into state management, routing, and API integration instead.

---

## Docker

```
Create a docker-compose.yml that orchestrates:
- PostgreSQL 17 Alpine container with health checks
- Go backend container that depends on postgres
- React frontend container that depends on backend
- Proper environment variables for DATABASE_URL, JWT_SECRET, etc.
- Volume for postgres data persistence
- Port mappings: 5432 for postgres, 8080 for backend, 5173 for frontend
```

One shot, worked mostly as-is.

---

## Where it actually broke

**CORS.** First backend had none. Frontend on 5173 calling 8080, blocked immediately, predictable.
```
The frontend is getting CORS errors when calling the backend API. 
Add CORS middleware to the Gin router that allows requests from http://localhost:5173
```
`gin-contrib/cors`, fixed in one pass.

**Connection string.** Generated `postgres://user:pass@host:port/db` — no `sslmode`. Docker Postgres rejected it outright. Not really a hallucination so much as an assumption that didn't hold locally. Fixed manually:
```
postgres://peekaping:peekaping@postgres:5432/peekaping?sslmode=disable
```

**Vite proxy.** `/auth`, `/monitors`, `/apikeys` all 404ing. Vite config had no proxy set up, frontend was hitting itself instead of the backend.
```
The frontend is getting 404s when calling API endpoints. 
Configure Vite's dev server proxy to forward /auth, /monitors, and /apikeys requests to http://backend:8080
```
```ts
server: {
  port: 5173,
  proxy: {
    "/auth": "http://backend:8080",
    "/monitors": "http://backend:8080",
    "/apikeys": "http://backend:8080"
  }
}
```

**Token expiry.** Nothing checked it. Backend happily accepted stale tokens, frontend had no idea what to do with a 401 either. Added expiry checks in the auth middleware and a redirect-to-login on the frontend when a token's dead.

**Scheduler scope creep.** At one point it proposed worker pools, retry queues, multiple notification providers, distributed scheduling — all reasonable, none of it for a few dozen URLs. Told it to optimize for something readable instead of something that scales to a load this project doesn't have. Ended up as a plain polling loop with explicit state transitions, which is honestly all it needed to be.

**SQLite → Postgres.** Started on SQLite since it needed zero setup. Swapped once it was clear the deployment story and the multi-container requirement wanted an actual service. Repository contracts didn't move, only the infra underneath them.

**Frontend state blob.** Everything — routing, fetching, render logic — was piling into one component early on. Split into pages / components / service layer / shared types / hooks once it got annoying to read.

---

## Roughly

- ~6–8 hours end to end, would've been a couple days doing this by hand
- Most of the codebase came out of prompts first, then got picked apart — CORS, connection strings, and the proxy config are the parts I'd flag as "AI got this wrong in a way you'd only catch by actually running it," everything else was more architecture judgment calls than bugs