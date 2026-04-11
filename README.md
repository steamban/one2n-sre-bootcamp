# Student CRUD REST API

## Purpose

A REST API for student management built with TypeScript, Express, and Prisma. Uses PostgreSQL for data persistence.

## Prerequisites

- **Node.js 20+**
- **PostgreSQL 14+**

## Setup

```bash
cp example.env .env
npm install
npm run db:generate
npm run migrate:dev
```

## Commands

| Command                  | Description                       |
| ------------------------ | --------------------------------- |
| `npm run dev`            | Start dev server (tsx watch)      |
| `npm run build`          | Build to `dist/`                  |
| `npm start`              | Run production server             |
| `npm run test`           | Run tests                         |
| `npm run lint`           | Lint code                         |
| `npm run format`         | Format code                       |
| `npm run migrate:dev`    | Create and apply migrations (dev) |
| `npm run migrate:deploy` | Apply migrations (production)     |
| `npm run db:generate`    | Generate Prisma client            |

## API Routes

- `GET /healthcheck`
- `POST /api/v1/students`
- `GET /api/v1/students`
- `GET /api/v1/students/:id`
- `PATCH /api/v1/students/:id`
- `DELETE /api/v1/students/:id`

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
