package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/internal/server"
	"github.com/DLzer/gojira/pkg/discord_cmd"
	"github.com/DLzer/gojira/pkg/logger"
	"github.com/DLzer/gojira/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	figure.NewFigure("GoJIRA", "isometric1", true).Print()

	stop := make(chan os.Signal, 1)

	// Loading Configuration
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// Starting Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	var discordSession *discordgo.Session
	if cfg.Discord.Enable {
		dg, err := discordgo.New("Bot " + cfg.Discord.BotToken)
		if err != nil {
			appLogger.Fatal(err)
			return
		}
		discordSession = dg
		appLogger.Infof("Discord Session Started")
	}

	if discordSession != nil {
		discordSession.Open()
		defer discordSession.Close()
	}

	discordCommandHandler := discord_cmd.NewDiscordCommandsHandler(discordSession, cfg.Discord.GuildID, cfg.Chatgpt.Token)
	discordCommandHandler.Init()
	discordCommandHandler.EnableCommands()
	appLogger.Infof("Discord Commands Enabled")

	// Starting Redis
	// redisDB := redis.NewRedisConnection(cfg.Redis.RedisAddr, cfg.Redis.RedisUsername, cfg.Redis.RedisPassword)

	// Start the server
	s := server.NewServer(cfg, nil, discordSession, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}

	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	discordCommandHandler.RemoveCommands()
}
