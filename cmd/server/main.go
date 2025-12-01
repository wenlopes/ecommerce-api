package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/category"
	category_gorm "github.com/mytheresa/go-hiring-challenge/app/category/gorm"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/log/terminal"
	product_gorm "github.com/mytheresa/go-hiring-challenge/app/product/gorm"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize handlers
	logger := terminal.NewTerminalLog()

	prodRepo := product_gorm.NewProductsRepository(db)
	catalog := catalog.NewCatalogHandler(prodRepo, logger)

	cateRepo := category_gorm.NewCategoryRepository(db)
	category := category.NewCategoryHandler(cateRepo, logger)

	// Set up the HTTP server
	srv := &http.Server{
		Addr: fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: api.NewMuxRouter(api.Handlers{
			Catalog:  catalog,
			Category: category,
		}),
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
