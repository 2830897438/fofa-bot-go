# FOFA Bot - Go Version

A powerful Telegram bot for FOFA (Cyberspace Assets Search Engine) written in Go.

## Features

- 🔍 FOFA asset search with full query syntax support
- 📦 Host information lookup
- 📊 Aggregated statistics
- 💾 Smart caching system
- 🕰️ Query history
- ⚙️ Flexible configuration

## Quick Start

1. Get your Telegram Bot Token from [@BotFather](https://t.me/BotFather)
2. Get your FOFA API Key from [FOFA](https://fofa.info)
3. Create `config.json`:

```json
{
  "bot_token": "YOUR_BOT_TOKEN",
  "apis": ["YOUR_FOFA_API_KEY"],
  "admins": [],
  "proxy": "",
  "full_mode": false,
  "public_mode": false,
  "presets": [],
  "update_url": ""
}
```

4. Run the bot:

```bash
go run cmd/bot/main.go
```

## Commands

- `/start` - Start the bot
- `/help` - Show help message
- `/search <query>` - Search FOFA
- `/host <ip|domain>` - Query host information
- `/history` - View query history
- `/settings` - View settings

## Build

```bash
go build -o fofa-bot cmd/bot/main.go
```

## License

For educational and research purposes only.

## Credits

Based on the original Python version: [fofa_bot](https://github.com/gagmm/fofa_bot)
