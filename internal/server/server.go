package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	redisClient *redis.Client
	dg          *discordgo.Session
	logger      logger.Logger
}

// NewServer constructor
func NewServer(cfg *config.Config, redisClient *redis.Client, dg *discordgo.Session, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, redisClient: redisClient, dg: dg, logger: logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	// Start the HTTP server
	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	// If Debug=true start Debug HTTP server
	if s.cfg.Server.Debug {
		go func() {
			s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
			if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()
	}

	// Map route handlers
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	// Detect quit and safely shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
