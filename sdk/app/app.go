package app

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"

	"hei-gin/docs"
	_ "hei-gin/docs"
	"hei-gin/sdk/config"
	"hei-gin/sdk/db"
	"hei-gin/sdk/middleware"
	"hei-gin/sdk/plugin"
	"hei-gin/sdk/registry"

	_ "hei-gin/sdk/auth"
	_ "hei-gin/sdk/captcha"
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
	r.Use(middleware.Recovery())  // Catch all panics → JSON
	r.Use(gin.Logger())           // Request logging
	r.Use(middleware.Trace())     // Trace ID injection
	r.Use(middleware.CORS())      // CORS headers
	r.Use(middleware.AuthCheck()) // Authentication

	registry.ApplyMiddlewares(r)

	SetupRouters(r)

	plugin.StartAll()

	addr := fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port)
	idleTimeout := time.Duration(config.C.App.TimeoutKeepAlive) * time.Second
	if idleTimeout <= 0 {
		idleTimeout = 15 * time.Second
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    1 << 20,
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

	setupSwagger(r)

	registry.ExecuteRoutes(r)
}

// setupSwagger registers Swagger UI routes when enabled via config.
// Swagger spec is served from the generated docs package.
func setupSwagger(r *gin.Engine) {
	if !config.C.Swagger.Enabled {
		return
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Decide whether to protect routes with Basic Auth
	useBasicAuth := config.C.Swagger.Username != "" && config.C.Swagger.Password != ""

	// ── Common OpenAPI spec endpoints (for tools like FoxAPI, Postman, etc.) ──
	specData := func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/doc.json")
	}
	if useBasicAuth {
		r.GET("/openapi.json", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), specData)
		r.GET("/v3/api-docs", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), specData)
		r.GET("/swagger.json", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), specData)
		r.GET("/api/v1/swagger/doc.json", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), specData)
	} else {
		r.GET("/openapi.json", specData)
		r.GET("/v3/api-docs", specData)
		r.GET("/swagger.json", specData)
		r.GET("/api/v1/swagger/doc.json", specData)
	}

	// ── Swagger UI ──
	swaggerHandler := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	)

	// Redirect bare paths to Swagger UI
	redirectToSwagger := func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	}
	if useBasicAuth {
		r.GET("/swagger", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), redirectToSwagger)
		r.GET("/docs", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), redirectToSwagger)
		r.GET("/redoc", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), redirectToSwagger)
		r.GET("/swagger/*any", basicAuth(config.C.Swagger.Username, config.C.Swagger.Password), swaggerHandler)
	} else {
		r.GET("/swagger", redirectToSwagger)
		r.GET("/docs", redirectToSwagger)
		r.GET("/redoc", redirectToSwagger)
		r.GET("/swagger/*any", swaggerHandler)
	}

	log.Println("[Swagger] Swagger UI enabled at /swagger/index.html")
}
// basicAuth returns a Gin middleware that enforces HTTP Basic Authentication
// with the given username and password using constant-time comparison.
func basicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			c.Header("WWW-Authenticate", `Basic realm="Swagger UI"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
