# Montiery

**Live Demo**: [https://assign-nine-delta.vercel.app](https://assign-nine-delta.vercel.app)

Montiery is a small uptime monitoring MVP. It focuses on the core product surface only: authentication, user-owned URL monitors, scheduled health checks, API keys, rate limiting, and email notifications.

The current implementation is split into a Go backend, a React frontend, PostgreSQL, and Docker Compose for local development and runtime.

## What it includes

- Register and login flows with JWT session handling
- Monitor CRUD for tracked URLs
- Health history and latest status for each monitor
- Availability and SLA summaries
- API key management for authenticated clients
- Email-only notification support
- Request rate limiting on auth and protected API routes

## Stack

- Backend: Go, Gin, GORM, PostgreSQL
- Frontend: React, Vite, TypeScript
- Runtime: Docker Compose

## Project layout

```text
backend/   Go API, services, repositories, scheduler, database, models, DTOs
frontend/  React UI, pages, shared components, demo data, utility helpers
design/    Open Design reference material used for the UI direction
```

## Quick Start (1-Line Setup)

```bash
docker compose up --build
```

This single command will:
- Start PostgreSQL database
- Build and run the Go backend API
- Build and run the React frontend

Services will be available at:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- PostgreSQL: localhost:5432

## Testing the Monitor (Verification Steps)

To verify that the uptime monitor correctly detects both "up" and "down" states:

1. **Start the application**: Run `docker compose up --build`

2. **Register an account**: 
   - Navigate to http://localhost:5173
   - Click "Register" and create a new account
   - Login with your credentials

3. **Add a working URL** (should show as "UP"):
   - Click "Add Monitor" or the "+" button
   - Enter: `https://example.com`
   - Set check interval (e.g., 60 seconds)
   - Click "Save"
   - Wait 1-2 minutes and verify the monitor shows:
     - Status: ✅ UP
     - Response time: ~100-500ms
     - Latest check timestamp

4. **Add a broken URL** (should show as "DOWN"):
   - Click "Add Monitor" again
   - Enter: `https://this-url-definitely-does-not-exist-12345.com`
   - Set check interval (e.g., 60 seconds)
   - Click "Save"
   - Wait 1-2 minutes and verify the monitor shows:
     - Status: ❌ DOWN
     - Error or timeout message
     - Latest check timestamp

5. **Verify the data**:
   - Check the dashboard to see both monitors
   - View availability percentages
   - Click on each monitor to see detailed history
   - Confirm response times are logged for the working URL
   - Confirm error states are logged for the broken URL

**Expected Results:**
- `example.com` should consistently show UP with response times
- The non-existent URL should show DOWN with error messages
- All checks should be logged with timestamps in the history

## Local development

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Backend:

```bash
cd backend
go run ./cmd/api
```

## Verification

Frontend build:

```bash
cd frontend
npm run build
```

Backend tests:

```bash
cd backend
go test ./...
```

Docker status:

```bash
docker compose ps
```

## API overview

Public:

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`

Protected:

- `GET /auth/profile`
- `GET /monitors`
- `POST /monitors`
- `PUT /monitors/:id`
- `DELETE /monitors/:id`
- `GET /monitors/:id/history`
- `GET /monitors/:id/latest`
- `GET /apikeys`
- `POST /apikeys`
- `DELETE /apikeys/:id`

Requests can authenticate with either:

- `Authorization: Bearer <jwt>`
- `X-API-Key: <key>`

## Configuration

Backend environment variables:

- `PORT`
- `DATABASE_URL`
- `JWT_SECRET`
- `DEFAULT_REQUEST_TIMEOUT_SECONDS`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USER`
- `SMTP_PASSWORD`
- `SMTP_FROM`

SMTP is optional. When it is not configured, notification events are skipped rather than sent.

## Deployment

The app is container-friendly, so the simplest production path is to build the backend and frontend images and run them behind a reverse proxy.

### Recommended split

1. Deploy the database first.
1. Deploy the backend second and point it at the managed database.
1. Deploy the frontend last and point it at the backend API URL.

### AWS

Recommended AWS setup:

1. Store the frontend and backend images in Amazon ECR.
2. Run the backend in ECS Fargate or on EC2 with the `DATABASE_URL`, `JWT_SECRET`, and SMTP variables set in task or instance environment.
3. Use Amazon RDS for PostgreSQL instead of the local Compose database.
4. Put the frontend behind CloudFront, S3, or a second ECS service depending on whether you want static hosting or a containerized UI.
5. Terminate TLS with an Application Load Balancer and point your domain DNS to it.

### Backend deployment steps

1. Build the backend image from `backend/`.
1. Push the image to your container registry.
1. Create an ECS Fargate service or EC2 container service for the backend.
1. Set backend environment variables in the service task definition:
   - `DATABASE_URL`
   - `JWT_SECRET`
   - `PORT`
   - `SMTP_HOST`
   - `SMTP_PORT`
   - `SMTP_USER`
   - `SMTP_PASSWORD`
   - `SMTP_FROM`
1. Point `DATABASE_URL` at the managed PostgreSQL instance.
1. Open the backend service only to the frontend origin and required internal traffic.

### Database deployment steps

1. Create an Amazon RDS for PostgreSQL instance.
1. Put it in the same region and network boundary as the backend service.
1. Save the database connection string and use it as `DATABASE_URL`.
1. Make sure the backend can connect over SSL if your provider requires it.
1. Back up the database before first production traffic.

### Frontend deployment steps

1. Deploy the `frontend/` directory as a static site.
1. Set `VITE_API_BASE_URL` to the public backend URL.
1. Add the frontend domain to the backend CORS allowlist.
1. Keep `frontend/vercel.json` in the repo so client-side routing falls back to `index.html`.

### Other options

- Render or Railway: deploy the backend and managed Postgres as separate services, then deploy the frontend as a static site or container.
- Fly.io: run the backend as an app, attach a Postgres instance, and deploy the frontend as a separate static app or container.
- A single VPS: use Docker Compose, a reverse proxy such as Caddy or Nginx, and your own TLS certificates.
- Vercel: deploy the `frontend/` directory as a static React app, set `VITE_API_BASE_URL` to the deployed backend URL, and add the frontend origin to the backend CORS allowlist.

### Production checklist

- Use a real `DATABASE_URL` pointing to managed Postgres.
- Rotate `JWT_SECRET` before first public deployment.
- Configure SMTP if notifications should be delivered.
- Lock CORS to the final frontend origin instead of localhost.
- For Vercel, set `frontend/vercel.json` as the deployment config and keep client-side routing on the SPA fallback.

### Vercel notes

Vercel is a good fit for the `frontend/` app because it is a static React build. Set the production environment variable in the Vercel project settings, then deploy the `frontend/` folder as the project root.

Official references:

- [Vercel environment variables](https://vercel.com/docs/environment-variables)
- [Vercel deployments](https://vercel.com/docs/deployments)
- [Amazon RDS for PostgreSQL](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html)

## Design source

The React UI follows the Open Design wireframes from:

design folder 

## Notes

- The app uses a demo preview mode in the frontend for inspection without backend state changes.
