# FOFA Bot (Go版本)

这是一个用 Go 语言重构的 FOFA Telegram 机器人，提供强大的网络空间资产搜索功能。

## ✨ 主要功能

- 🔍 **FOFA 资产搜索** - 支持完整的 FOFA 查询语法
- 📦 **主机信息查询** - 获取单个 IP 或域名的详细信息
- 📊 **聚合统计** - 对查询结果进行全局统计分析
- 💾 **智能缓存** - 自动缓存查询结果，节省 API 调用
- 🕰️ **查询历史** - 记录最近的查询历史
- ⚙️ **灵活配置** - 支持多 API Key、管理员权限等

## 📋 准备工作

1. **Telegram Bot Token**
   - 在 Telegram 中搜索 `@BotFather`
   - 发送 `/newbot` 创建新机器人
   - 获取 Bot Token

2. **FOFA API Key**
   - 注册 [FOFA](https://fofa.info) 账号
   - 在个人中心获取 API Key

## 🚀 快速开始

### 方式一：使用预编译二进制文件

1. 从 [Releases](https://github.com/yourusername/fofa-bot-go/releases) 下载对应平台的二进制文件

2. 创建配置文件 `config.json`：
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

3. 运行机器人：
```bash
chmod +x fofa-bot
./fofa-bot
```

### 方式二：从源码编译

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/fofa-bot-go.git
cd fofa-bot-go
```

2. 安装依赖：
```bash
go mod download
```

3. 编译：
```bash
go build -o fofa-bot cmd/bot/main.go
```

4. 配置并运行（同方式一）

## 📖 使用指南

### 基本命令

- `/start` - 启动机器人
- `/help` - 查看帮助信息
- `/search <query>` - FOFA 搜索
  - 示例：`/search domain="example.com"`
- `/host <ip|domain>` - 查询主机信息
  - 示例：`/host 1.1.1.1`
- `/stats <query>` - 聚合统计（开发中）
- `/history` - 查看查询历史
- `/settings` - 查看当前设置

### 权限说明

- 首次使用 `/start` 的用户会自动成为第一个管理员
- 大部分功能需要管理员权限
- 可在 `config.json` 中手动添加管理员 ID

## ⚙️ 配置说明

```json
{
  "bot_token": "Telegram Bot Token",
  "apis": ["FOFA API Key 列表"],
  "admins": [管理员 Telegram ID 列表],
  "proxy": "HTTP 代理地址（可选）",
  "full_mode": false,  // 是否使用完整模式查询
  "public_mode": false,  // 是否允许非管理员使用
  "presets": [],  // 预设查询（开发中）
  "update_url": ""  // 更新地址（开发中）
}
```

## 🔧 高级功能

### 缓存机制

- 查询结果自动缓存 24 小时
- 缓存文件存储在 `fofa_cache/` 目录
- 历史记录保存在 `history.json`

### 多 API Key 支持

在 `config.json` 中配置多个 API Key：
```json
{
  "apis": [
    "key1",
    "key2",
    "key3"
  ]
}
```

机器人会自动使用第一个可用的 Key。

## 📦 项目结构

```
fofa-bot-go/
├── cmd/
│   └── bot/
│       └── main.go          # 主程序入口
├── internal/
│   ├── bot/
│   │   └── bot.go           # Telegram bot 逻辑
│   ├── fofa/
│   │   └── client.go        # FOFA API 客户端
│   ├── config/
│   │   └── config.go        # 配置管理
│   └── cache/
│       └── cache.go         # 缓存管理
├── go.mod
├── go.sum
└── README.md
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目基于原 Python 版本重构，仅供学习和研究使用。

## 🙏 致谢

- 原项目：[fofa_bot](https://github.com/gagmm/fofa_bot)
- [FOFA](https://fofa.info) - 网络空间资产搜索引擎
- [Telegram Bot API](https://core.telegram.org/bots/api)

## ⚠️ 免责声明

本项目仅供学习和研究使用，请勿用于任何非法用途。使用本项目所造成的一切后果由使用者自行承担。

## 📞 支持

如有问题或建议，请提交 [Issue](https://github.com/yourusername/fofa-bot-go/issues)。
