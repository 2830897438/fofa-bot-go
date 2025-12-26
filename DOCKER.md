# Docker 部署指南

本文档介绍如何使用 Docker 部署 FOFA Bot。

## 📦 快速开始

### 方式一：使用 Docker Compose（推荐）

1. **创建配置文件**

```bash
# 下载示例配置
wget https://raw.githubusercontent.com/2830897438/fofa-bot-go/main/config.example.json -O config.json

# 编辑配置文件，填入你的 Bot Token 和 FOFA API Key
vim config.json
```

2. **下载 docker-compose.yml**

```bash
wget https://raw.githubusercontent.com/2830897438/fofa-bot-go/main/docker-compose.yml
```

3. **启动容器**

```bash
docker-compose up -d
```

4. **查看日志**

```bash
docker-compose logs -f
```

### 方式二：使用 Docker 命令

1. **拉取镜像**

```bash
# 从 GitHub Container Registry 拉取
docker pull ghcr.io/2830897438/fofa-bot-go:latest
```

2. **创建配置文件**

创建 `config.json` 文件，内容如下：

```json
{
  "bot_token": "YOUR_BOT_TOKEN_HERE",
  "apis": ["YOUR_FOFA_API_KEY_HERE"],
  "admins": [],
  "proxy": "",
  "full_mode": false,
  "public_mode": false,
  "presets": [],
  "update_url": ""
}
```

3. **运行容器**

```bash
docker run -d \
  --name fofa-bot \
  --restart unless-stopped \
  -v $(pwd)/config.json:/app/config.json:ro \
  -v $(pwd)/fofa_cache:/app/fofa_cache \
  -v $(pwd)/history.json:/app/history.json \
  -e TZ=Asia/Shanghai \
  ghcr.io/2830897438/fofa-bot-go:latest
```

4. **查看日志**

```bash
docker logs -f fofa-bot
```

### 方式三：从源码构建

1. **克隆仓库**

```bash
git clone https://github.com/2830897438/fofa-bot-go.git
cd fofa-bot-go
```

2. **构建镜像**

```bash
docker build -t fofa-bot-go:local .
```

3. **运行容器**

```bash
docker run -d \
  --name fofa-bot \
  --restart unless-stopped \
  -v $(pwd)/config.json:/app/config.json:ro \
  -v $(pwd)/fofa_cache:/app/fofa_cache \
  -v $(pwd)/history.json:/app/history.json \
  fofa-bot-go:local
```

## 🔧 高级配置

### 使用代理

如果需要通过代理访问 Telegram，可以设置环境变量：

```bash
docker run -d \
  --name fofa-bot \
  -v $(pwd)/config.json:/app/config.json:ro \
  -e HTTP_PROXY=http://proxy.example.com:8080 \
  -e HTTPS_PROXY=http://proxy.example.com:8080 \
  ghcr.io/2830897438/fofa-bot-go:latest
```

或在 `docker-compose.yml` 中添加：

```yaml
environment:
  - HTTP_PROXY=http://proxy.example.com:8080
  - HTTPS_PROXY=http://proxy.example.com:8080
```

### 使用主机网络

```bash
docker run -d \
  --name fofa-bot \
  --network host \
  -v $(pwd)/config.json:/app/config.json:ro \
  ghcr.io/2830897438/fofa-bot-go:latest
```

### 持久化数据

容器使用以下目录存储数据：

- `/app/config.json` - 配置文件（只读）
- `/app/fofa_cache/` - 查询结果缓存
- `/app/history.json` - 查询历史

确保挂载这些目录以持久化数据。

## 📊 容器管理

### 查看运行状态

```bash
# 使用 docker-compose
docker-compose ps

# 使用 docker
docker ps | grep fofa-bot
```

### 停止容器

```bash
# 使用 docker-compose
docker-compose stop

# 使用 docker
docker stop fofa-bot
```

### 重启容器

```bash
# 使用 docker-compose
docker-compose restart

# 使用 docker
docker restart fofa-bot
```

### 删除容器

```bash
# 使用 docker-compose
docker-compose down

# 使用 docker
docker stop fofa-bot && docker rm fofa-bot
```

### 更新镜像

```bash
# 拉取最新镜像
docker pull ghcr.io/2830897438/fofa-bot-go:latest

# 重新创建容器
docker-compose up -d --force-recreate
```

## 🐛 故障排查

### 查看日志

```bash
# 查看实时日志
docker logs -f fofa-bot

# 查看最近 100 行日志
docker logs --tail 100 fofa-bot
```

### 进入容器

```bash
docker exec -it fofa-bot sh
```

### 检查配置文件

```bash
docker exec fofa-bot cat /app/config.json
```

### 检查网络连接

```bash
# 测试容器网络
docker exec fofa-bot ping -c 4 api.telegram.org
```

## 📦 可用镜像标签

- `latest` - 最新稳定版本（main 分支）
- `v1.0.0` - 具体版本号
- `v1.0` - 次版本号
- `v1` - 主版本号

## 🔒 安全建议

1. **配置文件权限**：确保 `config.json` 文件权限设置正确（建议 600）
2. **只读挂载**：配置文件使用只读挂载（`:ro`）
3. **非 root 用户**：容器内使用非 root 用户运行（UID 1000）
4. **网络隔离**：不需要时避免使用 host 网络模式

## 📝 示例配置

完整的 `docker-compose.yml` 示例：

```yaml
version: '3.8'

services:
  fofa-bot:
    image: ghcr.io/2830897438/fofa-bot-go:latest
    container_name: fofa-bot
    restart: unless-stopped

    volumes:
      - ./config.json:/app/config.json:ro
      - ./fofa_cache:/app/fofa_cache
      - ./history.json:/app/history.json

    environment:
      - TZ=Asia/Shanghai

    healthcheck:
      test: ["CMD", "pgrep", "-f", "fofa-bot"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🆘 获取帮助

如有问题，请提交 [Issue](https://github.com/2830897438/fofa-bot-go/issues)。
