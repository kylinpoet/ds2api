# DS2API Deployment Guide (Go)

Language: [中文](DEPLOY.md) | [English](DEPLOY.en.md)

This guide is aligned with the current Go codebase.

## Deployment Modes

- Local run: `go run ./cmd/ds2api`
- Docker: `docker-compose up -d`
- Vercel: serverless entry at `api/index.go`
- Linux service mode: systemd

## 0. Prerequisites

- Go 1.25+
- Node.js 20+ (only if you need to build WebUI locally)
- `config.json` or `DS2API_CONFIG_JSON`

## 1. Local Run

```bash
git clone https://github.com/CJackHwang/ds2api.git
cd ds2api

cp config.example.json config.json
# edit config.json

go run ./cmd/ds2api
```

Default port is `5001` (override with `PORT`).

Build WebUI if `/admin` reports missing assets:

```bash
./scripts/build-webui.sh
```

## 2. Docker Deployment

```bash
cp .env.example .env
# edit .env

docker-compose up -d
docker-compose logs -f
```

Rebuild after updates:

```bash
docker-compose up -d --build
```

Notes:

- `Dockerfile` uses multi-stage build (WebUI + Go binary)
- Container entry command is `/usr/local/bin/ds2api`

## 3. Vercel Deployment

- Serverless entry: `api/index.go`
- Rewrites and cache headers: `vercel.json`
- Legacy `builds` has been removed to avoid the `unused-build-settings` warning

Minimum environment variables:

- `DS2API_ADMIN_KEY`
- `DS2API_CONFIG_JSON` (raw JSON or Base64)

Optional:

- `VERCEL_TOKEN`
- `VERCEL_PROJECT_ID`
- `VERCEL_TEAM_ID`
- `DS2API_ACCOUNT_MAX_INFLIGHT` (per-account inflight limit, default `2`)
- `DS2API_ACCOUNT_CONCURRENCY` (alias of the same setting)

Recommended concurrency is computed dynamically as `account_count * per_account_inflight_limit` (default is `account_count * 2`).

After deploy, verify:

- `/healthz`
- `/v1/models`
- `/admin`

## 3.1 GitHub Release Automation

This repo includes `.github/workflows/release-artifacts.yml`:

- Triggers only on Release `published`
- Does not run on `push`
- Builds Linux/macOS/Windows archives and uploads them to Release Assets
- Generates `sha256sums.txt` for integrity checks

## 4. Reverse Proxy (Nginx)

Disable buffering for SSE:

```nginx
location / {
    proxy_pass http://127.0.0.1:5001;
    proxy_http_version 1.1;
    proxy_set_header Connection "";
    proxy_buffering off;
    proxy_cache off;
    chunked_transfer_encoding on;
    tcp_nodelay on;
}
```

## 5. systemd Example (Linux)

```ini
[Unit]
Description=DS2API (Go)
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/ds2api
Environment=PORT=5001
Environment=DS2API_CONFIG_PATH=/opt/ds2api/config.json
Environment=DS2API_ADMIN_KEY=admin
ExecStart=/opt/ds2api/ds2api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Common commands:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ds2api
sudo systemctl start ds2api
sudo systemctl status ds2api
```

## 6. Post-Deploy Checks

```bash
curl -s http://127.0.0.1:5001/healthz
curl -s http://127.0.0.1:5001/readyz
curl -s http://127.0.0.1:5001/v1/models
```

If admin UI is required:

```bash
curl -s http://127.0.0.1:5001/admin
```
