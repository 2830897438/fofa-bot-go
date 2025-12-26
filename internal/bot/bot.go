package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/2830897438/fofa-bot-go/internal/cache"
	"github.com/2830897438/fofa-bot-go/internal/config"
	"github.com/2830897438/fofa-bot-go/internal/fofa"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot
type Bot struct {
	api          *tgbotapi.BotAPI
	config       *config.Config
	cacheManager *cache.Manager
}

// New creates a new bot instance
func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:          api,
		config:       cfg,
		cacheManager: cache.NewManager(),
	}, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	if err := b.cacheManager.Init(); err != nil {
		return err
	}

	// Set bot commands
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "🚀 启动机器人"},
		{Command: "help", Description: "❓ 命令手册"},
		{Command: "search", Description: "🔍 FOFA搜索"},
		{Command: "host", Description: "📦 主机查询"},
		{Command: "stats", Description: "📊 统计信息"},
		{Command: "history", Description: "🕰️ 查询历史"},
		{Command: "settings", Description: "⚙️ 设置"},
	}

	cmdConfig := tgbotapi.NewSetMyCommands(commands...)
	if _, err := b.api.Request(cmdConfig); err != nil {
		log.Printf("Failed to set commands: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if !message.IsCommand() {
		return
	}

	// Check admin for restricted commands
	restrictedCommands := []string{"search", "host", "stats", "settings"}
	command := message.Command()

	isRestricted := false
	for _, cmd := range restrictedCommands {
		if command == cmd {
			isRestricted = true
			break
		}
	}

	if isRestricted && !b.config.IsAdmin(message.From.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "⛔️ 抱歉，您没有权限执行此操作。")
		b.api.Send(msg)
		return
	}

	switch command {
	case "start":
		b.handleStart(message)
	case "help":
		b.handleHelp(message)
	case "search":
		b.handleSearch(message)
	case "host":
		b.handleHost(message)
	case "stats":
		b.handleStats(message)
	case "history":
		b.handleHistory(message)
	case "settings":
		b.handleSettings(message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "未知命令。使用 /help 查看可用命令。")
		b.api.Send(msg)
	}
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	// Auto-add first admin
	if len(b.config.Admins) == 0 {
		b.config.Admins = append(b.config.Admins, message.From.ID)
		if err := b.config.Save("config.json"); err != nil {
			log.Printf("Failed to save config: %v", err)
		}

		msg := tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("ℹ️ 已自动将您 (ID: %d) 添加为第一个管理员。", message.From.ID))
		b.api.Send(msg)
	}

	text := fmt.Sprintf(`👋 欢迎, %s！

这是一个功能强大的 FOFA Telegram 机器人。

🔍 主要功能：
• /search - FOFA 资产搜索
• /host - 主机详细信息查询
• /stats - 全局聚合统计
• /history - 查看查询历史

使用 /help 查看完整命令列表。`, message.From.FirstName)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	b.api.Send(msg)
}

func (b *Bot) handleHelp(message *tgbotapi.Message) {
	text := `📖 *FOFA 机器人命令手册*

*🔍 资产搜索*
/search <query> - FOFA搜索
示例: /search domain="example.com"

*📦 主机查询*
/host <ip|domain> - 查询主机信息
示例: /host 1.1.1.1

*📊 统计分析*
/stats <query> - 全局聚合统计
示例: /stats app="nginx"

*📚 其他功能*
/history - 查看查询历史
/settings - 设置管理（管理员）
/help - 显示此帮助信息

*提示：* 大部分功能需要管理员权限。`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

func (b *Bot) handleSearch(message *tgbotapi.Message) {
	args := strings.TrimSpace(message.CommandArguments())
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "用法: /search <FOFA查询语法>\n\n示例:\n/search domain=\"example.com\"")
		b.api.Send(msg)
		return
	}

	// Send processing message
	processingMsg := tgbotapi.NewMessage(message.Chat.ID, "⏳ 正在查询...")
	sentMsg, _ := b.api.Send(processingMsg)

	// Check cache first
	cached := b.cacheManager.FindCache(args)
	if cached != nil {
		elapsed := time.Since(cached.Timestamp)
		if elapsed < cache.CacheExpiration {
			// Send cached file
			doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(cached.FilePath))
			doc.Caption = fmt.Sprintf("✅ 从缓存返回结果\n查询: %s\n共 %d 条记录\n缓存时间: %s",
				args, cached.Count, cached.Timestamp.Format("2006-01-02 15:04:05"))
			b.api.Send(doc)

			deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, sentMsg.MessageID)
			b.api.Request(deleteMsg)
			return
		}
	}

	// Perform search
	if len(b.config.APIs) == 0 {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, "❌ 未配置 FOFA API Key")
		b.api.Send(editMsg)
		return
	}

	client := fofa.NewClient(b.config.APIs[0])
	result, err := client.Search(args, 1, 10000, "host", b.config.FullMode)
	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID,
			fmt.Sprintf("❌ 查询失败: %v", err))
		b.api.Send(editMsg)
		return
	}

	if result.Size == 0 {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, "🤷‍♀️ 未找到结果。")
		b.api.Send(editMsg)
		return
	}

	// Save results to file
	filename := fmt.Sprintf("fofa_%d.txt", time.Now().Unix())
	filepath := b.cacheManager.GetCachePath(filename)

	file, err := os.Create(filepath)
	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID,
			fmt.Sprintf("❌ 创建文件失败: %v", err))
		b.api.Send(editMsg)
		return
	}
	defer file.Close()

	for _, host := range result.Results {
		file.WriteString(host + "\n")
	}

	// Add to cache
	b.cacheManager.AddQuery(args, filepath, len(result.Results))

	// Send file
	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(filepath))
	doc.Caption = fmt.Sprintf("✅ 查询完成\n查询: %s\n共 %d 条记录（已下载 %d 条）",
		args, result.Size, len(result.Results))
	b.api.Send(doc)

	// Delete processing message
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, sentMsg.MessageID)
	b.api.Request(deleteMsg)
}

func (b *Bot) handleHost(message *tgbotapi.Message) {
	args := strings.TrimSpace(message.CommandArguments())
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "用法: /host <ip或域名>\n\n示例:\n/host 1.1.1.1")
		b.api.Send(msg)
		return
	}

	processingMsg := tgbotapi.NewMessage(message.Chat.ID, "⏳ 正在查询主机信息...")
	sentMsg, _ := b.api.Send(processingMsg)

	if len(b.config.APIs) == 0 {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, "❌ 未配置 FOFA API Key")
		b.api.Send(editMsg)
		return
	}

	// Determine if it's an IP or domain
	query := fmt.Sprintf("host=\"%s\"", args)

	client := fofa.NewClient(b.config.APIs[0])
	result, err := client.Search(query, 1, 100, "ip,port,protocol,title,server", b.config.FullMode)
	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID,
			fmt.Sprintf("❌ 查询失败: %v", err))
		b.api.Send(editMsg)
		return
	}

	if result.Size == 0 {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID,
			fmt.Sprintf("🤷‍♀️ 未找到关于 %s 的任何信息。", args))
		b.api.Send(editMsg)
		return
	}

	// Format results
	text := fmt.Sprintf("📌 *主机信息: %s*\n\n共发现 %d 个服务\n\n", args, result.Size)

	count := 0
	for _, host := range result.Results {
		if count >= 10 {
			text += "\n_（仅显示前10条）_"
			break
		}
		text += fmt.Sprintf("• %s\n", host)
		count++
	}

	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, text)
	editMsg.ParseMode = "Markdown"
	b.api.Send(editMsg)
}

func (b *Bot) handleStats(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "📊 统计功能开发中...")
	b.api.Send(msg)
}

func (b *Bot) handleHistory(message *tgbotapi.Message) {
	history, err := b.cacheManager.LoadHistory()
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("❌ 加载历史失败: %v", err))
		b.api.Send(msg)
		return
	}

	if len(history.Queries) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "查询历史为空。")
		b.api.Send(msg)
		return
	}

	text := "*🕰️ 最近查询历史*\n\n"
	for i, q := range history.Queries {
		if i >= 10 {
			break
		}
		text += fmt.Sprintf("%d. `%s`\n   _%s_ (%d条)\n\n",
			i+1, q.QueryText, q.Timestamp.Format("2006-01-02 15:04"), q.Count)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

func (b *Bot) handleSettings(message *tgbotapi.Message) {
	text := fmt.Sprintf(`⚙️ *当前设置*

*API Keys:* %d 个
*管理员:* %d 个
*完整模式:* %v
*公开模式:* %v

使用配置文件 config.json 修改设置。`,
		len(b.config.APIs),
		len(b.config.Admins),
		b.config.FullMode,
		b.config.PublicMode)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}
