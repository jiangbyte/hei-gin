package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/config"
	"hei-gin/sdk/db"
	"hei-gin/sdk/middleware"
	"hei-gin/sdk/plugin"
	"hei-gin/sdk/registry"

	_ "hei-gin/sdk/auth"
	_ "hei-gin/sdk/captcha"
	_ "hei-gin/sdk/scheduler"
	_ "hei-gin/sdk/utils"
)

func Run() {
	if err := config.FindAndLoad(); err != nil {
		log.Fatalf("[APP] Failed to load config: %v", err)
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("[APP] Failed to init database: %v", err)
	}

	if err := db.InitRedis(); err != nil {
		log.Fatalf("[APP] Failed to init Redis: %v", err)
	}


	if err := plugin.InitAll(); err != nil {
		log.Fatalf("[APP] Plugin init failed: %v", err)
	}

	// NOTE: Use gin.New() NOT gin.Default(). gin.Default() includes gin.Recovery()
	// which returns HTML on panic, breaking our JSON API contract.
	// Our custom middleware.Recovery() handles panics with proper JSON output + stack logging.
	r := gin.New()

	// Middleware order matters — outermost first (Recovery must be outermost to catch all)
	r.Use(middleware.Recovery()) // Catch all panics → JSON
	r.Use(gin.Logger())          // Request logging
	r.Use(middleware.Trace())    // Trace ID injection
	r.Use(middleware.CORS())     // CORS headers
	r.Use(middleware.AuthCheck()) // Authentication

	registry.ApplyMiddlewares(r)

	SetupRouters(r)

			plugin.StartAll()

	addr := fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("[APP] Server started on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[APP] Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[APP] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[APP] Server forced to shutdown: %v", err)
	}

	plugin.StopAll()

	db.Close()
	db.CloseRedis()
	log.Println("[APP] Server exited")
}

func SetupRouters(r *gin.Engine) {
	r.GET("/", HealthHandler)
	registry.ExecuteRoutes(r)
}
