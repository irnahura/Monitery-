# AI Collaboration Log

## 🤖 The AI Tech Stack

**Primary AI Assistant**: Kiro (Claude Sonnet 4.5)
**Development Environment**: Kiro AI-powered development environment
**Model**: Claude Sonnet 4.5 by Anthropic

## 📝 The Prompts that Shipped It

### Initial Project Setup Prompt
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

### Backend Framework Generation
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

### Frontend UI Generation
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

### Docker Compose Setup
```
Create a docker-compose.yml that orchestrates:
- PostgreSQL 17 Alpine container with health checks
- Go backend container that depends on postgres
- React frontend container that depends on backend
- Proper environment variables for DATABASE_URL, JWT_SECRET, etc.
- Volume for postgres data persistence
- Port mappings: 5432 for postgres, 8080 for backend, 5173 for frontend
```

## 🔧 The Course Corrections

### Issue 1: CORS Configuration
**Problem**: Initially, the AI generated a backend with no CORS middleware, causing the frontend to fail with CORS policy errors when making API requests from http://localhost:5173 to http://localhost:8080.

**AI's Mistake**: The Gin router was set up without CORS headers, blocking cross-origin requests.

**Resolution Prompt**: 
```
The frontend is getting CORS errors when calling the backend API. 
Add CORS middleware to the Gin router that allows requests from http://localhost:5173
```

**Fix**: AI added `github.com/gin-contrib/cors` middleware with proper AllowOrigins, AllowMethods, and AllowCredentials configuration.

### Issue 2: Database Connection String Format
**Problem**: The AI initially used an incorrect PostgreSQL connection string format that didn't include the sslmode parameter, causing connection failures in the Docker environment.

**AI's Hallucination**: Generated `DATABASE_URL` as `postgres://user:pass@host:port/db` without the `?sslmode=disable` suffix needed for local development.

**Resolution**: Manually corrected the connection string in docker-compose.yml to:
```
postgres://peekaping:peekaping@postgres:5432/peekaping?sslmode=disable
```

### Issue 3: Frontend Proxy Configuration
**Problem**: The Vite dev server wasn't properly proxying API requests to the backend, resulting in 404 errors for all `/auth`, `/monitors`, and `/apikeys` endpoints.

**AI's Initial Code**: Generated a basic Vite config without proxy settings, requiring the frontend to make direct requests to `http://localhost:8080` which created CORS complications.

**Resolution Prompt**:
```
The frontend is getting 404s when calling API endpoints. 
Configure Vite's dev server proxy to forward /auth, /monitors, and /apikeys requests to http://backend:8080
```

**Fix**: AI updated vite.config.ts with proper proxy configuration:
```typescript
server: {
  port: 5173,
  proxy: {
    "/auth": "http://backend:8080",
    "/monitors": "http://backend:8080",
    "/apikeys": "http://backend:8080"
  }
}
```

### Issue 4: JWT Token Expiration Handling
**Problem**: The AI generated JWT authentication logic but didn't implement proper token expiration validation on the backend, and the frontend didn't handle expired token scenarios gracefully.

**Course Correction**: Added explicit token expiration checks in the auth middleware and frontend redirect logic to the login page when tokens are invalid or expired.

## 📊 Development Process Summary

1. **Architecture Design** (AI-assisted): Started with the assignment requirements, AI proposed the Go + React + PostgreSQL stack
2. **Backend Scaffolding** (95% AI-generated): Gin router, GORM models, JWT auth, scheduler, all generated via prompts
3. **Frontend Components** (90% AI-generated): React pages, routing, API integration, all based on design wireframes
4. **Docker Orchestration** (100% AI-generated): Complete docker-compose.yml with multi-stage builds
5. **Debugging & Fixes** (Human + AI collaboration): CORS, proxy configs, connection strings refined through iterative prompting
6. **Deployment** (AI-assisted): Vercel deployment configuration and execution

**Total Development Time**: ~6-8 hours (would have taken 20+ hours without AI assistance)

**AI Contribution**: ~85-90% of the codebase generated, human provided architecture decisions, debugging, and refinement prompts
