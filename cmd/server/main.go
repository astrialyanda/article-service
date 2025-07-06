package main

import (
	"context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "article-service/internal/config"
    "article-service/internal/handler"
    "article-service/internal/repository"
    "article-service/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    
    // Initialize database
    db, err := repository.NewDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize repository, service, and handler
    articleRepo := repository.NewArticleRepository(db)
    articleService := service.NewArticleService(articleRepo)
    articleHandler := handler.NewArticleHandler(articleService)

    // Setup routes
    router := gin.Default()
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    // Article routes
    v1 := router.Group("/api/v1")
    {
        v1.POST("/articles", articleHandler.CreateArticle)
        v1.GET("/articles", articleHandler.GetArticles)
    }

    // Start server with graceful shutdown
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }

    go func() {
        log.Printf("Server starting on port %s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}