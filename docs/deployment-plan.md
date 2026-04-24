# Deployment plan

## Suggested hosting
- Frontend: Vercel or Netlify
- Go backend: Render, Fly.io, or Railway
- Python service: Render or Railway
- Database: Neon, Supabase Postgres, or managed Render Postgres

## Production improvements before deploy
- Replace demo secrets with environment variables
- Add structured logging
- Add rate limiting on auth endpoints
- Add CORS configuration
- Add database migrations via dedicated migration tooling
- Add request validation and better error envelopes
- Add monitoring and uptime checks

## Recommended roadmap
1. Deploy frontend first
2. Deploy backend and Python service with shared environment config
3. Attach managed Postgres
4. Add seeded demo account for recruiters
5. Add public screenshots and a short architecture diagram to the repo
