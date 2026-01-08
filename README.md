# Mail Security Server (Go + React)

Academic project: an SMTP server with a web backend/UI for email security checks and monitoring.
Focus: SPF/DKIM/DMARC validation, suspicious message detection, and attachment scanning (optional integration).

> ⚠️ This repository is intended for **defensive/security testing** in controlled environments.
> No phishing functionality or credential collection is included in the public version.

## Tech Stack
- Backend: Go
- Frontend: React (Create React App)
- DB: PostgreSQL
- Docker

## Features
- SMTP server running on port 587 (STARTTLS / TLS config)
- Email validation pipeline (SPF / DKIM / DMARC modules)
- Backend API for logs/results and admin UI
- Optional attachment scanning / external scanner integration (keys are not stored in repo)

## Project Structure
- `backend/` — Go services, SMTP server, validators
- `frontend/` — React UI
- `init.sql` — DB init script
- `Dockerfile` — DB/container setup (or services)

## Quick Start (Ubuntu)
### 1) Database
Build and run Postgres container:
docker build -t mail-postgres .
docker run -d --name mail-db -p 5432:5432 mail-postgres
2) Backend

Create config:

cp config.example.json config.json


Run:

cd backend
go mod download
sudo go run cmd/main.go

3) Frontend
cd frontend
npm install
npm start
Open UI: http://localhost:3000
Troubleshooting
Check if port 587 is busy:
sudo lsof -i :587
sudo kill -9 <PID>
