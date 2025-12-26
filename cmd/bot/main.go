package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/2830897438/fofa-bot-go/internal/bot"
	"github.com/2830897438/fofa-bot-go/internal/config"
)

const (
	ConfigFile = "config.json"
)

func main() {
	log.Println("🚀 FOFA Bot (Go版本) 启动中...")

	// Load configuration
	cfg, err := config.Load(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("配置文件不存在，创建默认配置...")
			cfg = createDefaultConfig()
			if err := cfg.Save(ConfigFile); err != nil {
				log.Fatalf("保存配置文件失败: %v", err)
			}
			log.Println("✅ 已创建默认配置文件 config.json")
			log.Println("请编辑 config.json 文件，填入您的 Bot Token 和 FOFA API Key")
			return
		}
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// Validate configuration
	if cfg.BotToken == "" || cfg.BotToken == "YOUR_BOT_TOKEN_HERE" {
		log.Fatal("❌ 错误: 请在 config.json 中设置有效的 bot_token")
	}

	if len(cfg.APIs) == 0 || cfg.APIs[0] == "YOUR_FOFA_API_KEY_HERE" {
		log.Fatal("❌ 错误: 请在 config.json 中设置至少一个有效的 FOFA API Key")
	}

	// Create bot instance
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("创建机器人失败: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("\n🛑 收到停止信号，正在关闭...")
		os.Exit(0)
	}()

	// Start bot
	log.Println("✅ 机器人已启动，等待消息...")
	if err := b.Start(); err != nil {
		log.Fatalf("启动机器人失败: %v", err)
	}
}

func createDefaultConfig() *config.Config {
	return &config.Config{
		BotToken:   "YOUR_BOT_TOKEN_HERE",
		APIs:       []string{"YOUR_FOFA_API_KEY_HERE"},
		Admins:     []int64{},
		Proxy:      "",
		FullMode:   false,
		PublicMode: false,
		Presets:    []config.Preset{},
		UpdateURL:  "",
	}
}
