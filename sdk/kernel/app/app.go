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

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"hei-gin/docs"
	_ "hei-gin/docs"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/infra/scheduler"
	"hei-gin/sdk/kernel/plugin"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/observability"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/middleware"
)

func Run() {
	auth.RegisterPlugin()
	captcha.RegisterPlugin()
	utils.RegisterPlugin()
	scheduler.RegisterPlugin()

	if err := config.FindAndLoad(); err != nil {
		log.Fatalf("[APP] Failed to load config: %v", err)
	}
	if err := config.C.ValidateRuntime(true); err != nil {
		log.Fatalf("[APP] Invalid config: %v", err)
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
	registry.Freeze()
	db.Freeze()
	logAssemblySummary()

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Metrics())
	r.Use(gin.Logger())
	r.Use(middleware.Trace())
	r.Use(middleware.CORS())
	r.Use(middleware.AuthCheck())

	registry.ApplyMiddlewares(r)
	SetupRouters(r)
	if err := plugin.StartAll(); err != nil {
		log.Fatalf("[APP] Plugin start failed: %v", err)
	}

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
	r.GET("/health/live", LiveHandler)
	r.GET("/health/ready", ReadyHandler)
	if config.C.App.Debug {
		r.GET("/debug/registry", RegistryHandler)
	}
	r.GET("/metrics", gin.WrapH(observability.Handler()))
	setupSwagger(r)
	registry.ExecuteRoutes(r)
}

func logAssemblySummary() {
	pluginSnapshot := plugin.Snapshot()
	routeSnapshot := registry.SnapshotState()
	migrationSnapshot := db.Snapshot()

	log.Printf(
		"[APP] Assembly frozen: plugins=%d routes=%d middlewares=%d models=%d seeds=%d",
		len(pluginSnapshot),
		len(routeSnapshot.Routes),
		len(routeSnapshot.Middlewares),
		len(migrationSnapshot.Models),
		len(migrationSnapshot.Seeds),
	)
}

func setupSwagger(r *gin.Engine) {
	if !config.C.Swagger.Enabled {
		return
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	useBasicAuth := config.C.Swagger.Username != "" && config.C.Swagger.Password != ""

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

	swaggerHandler := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	)

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
