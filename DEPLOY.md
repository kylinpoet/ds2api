# DS2API 部署指南（Go）

语言 / Language: [中文](DEPLOY.md) | [English](DEPLOY.en.md)

本指南基于当前 Go 代码库。

## 部署方式

- 本地运行：`go run ./cmd/ds2api`
- Docker：`docker-compose up -d`
- Vercel：`api/index.go` serverless 入口
- Linux 服务化：systemd

## 0. 前置要求

- Go 1.25+
- Node.js 20+（仅在需要本地构建 WebUI 时）
- `config.json` 或 `DS2API_CONFIG_JSON`

## 1. 本地运行

```bash
git clone https://github.com/CJackHwang/ds2api.git
cd ds2api

cp config.example.json config.json
# 编辑 config.json

go run ./cmd/ds2api
```

默认监听 `5001`，可通过 `PORT` 覆盖。

构建 WebUI（可选，仅当 `/admin` 缺少静态文件时）：

```bash
./scripts/build-webui.sh
```

## 2. Docker 部署

```bash
cp .env.example .env
# 编辑 .env

docker-compose up -d
docker-compose logs -f
```

更新镜像：

```bash
docker-compose up -d --build
```

说明：

- `Dockerfile` 使用多阶段构建（WebUI + Go 二进制）
- 容器内默认启动命令：`/usr/local/bin/ds2api`

## 3. Vercel 部署

- serverless 入口：`api/index.go`
- 路由与缓存头：`vercel.json`
- 已移除 legacy `builds` 字段，避免 `unused-build-settings` 警告

至少配置环境变量：

- `DS2API_ADMIN_KEY`
- `DS2API_CONFIG_JSON`（JSON 或 Base64）

可选：

- `VERCEL_TOKEN`
- `VERCEL_PROJECT_ID`
- `VERCEL_TEAM_ID`
- `DS2API_ACCOUNT_MAX_INFLIGHT`（每账号并发上限，默认 `2`）
- `DS2API_ACCOUNT_CONCURRENCY`（同上别名）

并发建议值会动态按 `账号数量 × 每账号并发上限` 计算（默认即 `账号数量 × 2`）。

部署后建议先访问：

- `/healthz`
- `/v1/models`
- `/admin`

## 3.1 GitHub Release 自动构建

仓库包含 `.github/workflows/release-artifacts.yml`：

- 仅在 Release `published` 时触发
- 不在 `push` 时触发
- 自动构建 Linux/macOS/Windows 二进制包并上传到 Release Assets
- 生成 `sha256sums.txt` 供校验

## 4. 反向代理（Nginx）

如果在 Nginx 后挂载，建议关闭缓冲以保证 SSE：

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

## 5. systemd 示例（Linux）

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

常用命令：

```bash
sudo systemctl daemon-reload
sudo systemctl enable ds2api
sudo systemctl start ds2api
sudo systemctl status ds2api
```

## 6. 部署后检查

```bash
curl -s http://127.0.0.1:5001/healthz
curl -s http://127.0.0.1:5001/readyz
curl -s http://127.0.0.1:5001/v1/models
```

如果你依赖管理台接口，再检查：

```bash
curl -s http://127.0.0.1:5001/admin
```
