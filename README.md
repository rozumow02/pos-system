# POS System MVP

Local-network retail POS and inventory MVP for small electronics shops.

## Stack

- Frontend: Nuxt 4 + TypeScript
- Backend: Go + Fiber
- Database: PostgreSQL
- Deployment: Docker Compose + Nginx

## Architecture

- Nginx is the LAN entrypoint and serves the frontend while proxying `/api` to the Go backend.
- The backend is stateless, runs migrations on startup, and stores all business data in PostgreSQL.
- The frontend keeps the product catalog in memory for fast POS search and refreshes after product edits or completed sales.

## Database schema

- `products`: catalog, price, stock, status
- `orders`: sale headers
- `order_items`: sale line items
- `inventory_movements`: stock adjustments and sales ledger

## Project structure

- `backend/cmd/api`: API bootstrap
- `backend/internal`: config, HTTP, handlers, services, repositories, platform
- `backend/migrations`: schema and seed migrations
- `frontend/pages`: dashboard, products, POS, reports
- `frontend/components`: reusable UI building blocks
- `frontend/composables`: API/data helpers and POS cart state
- `deploy`: nginx and compose copies for deployment

## Run

1. Copy `.env.example` to `.env` if you want custom credentials or ports.
2. Run `docker compose up --build`.
3. Open `http://localhost` or your machine LAN IP.

## Frontend Tooling

- Install frontend dependencies with `npm.cmd install` from [frontend/package.json](C:/Users/User/projects/pos-system/frontend/package.json) when using PowerShell on Windows.
- Run `npm.cmd run typecheck` for Vue and TypeScript checks.
- Run `npm.cmd run lint` for ESLint checks.
- Run `npm.cmd run format` to apply the ESLint-based formatter.
- Run `npm.cmd run check` to execute typecheck, lint, and production build together.

## Key flows

- Add or edit products from `/products`
- Search and sell from `/pos`
- Review today's analytics on `/dashboard`
- Inspect top products and inventory status on `/reports`

## Notes

- Authentication is intentionally out of scope for this MVP.
- Stock never goes below zero.
- The backend calculates order totals from stored product prices; the client total is only for UX.
