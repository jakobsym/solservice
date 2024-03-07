package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func New() *App {
	app := &App{
		db: SetupDB(),
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// pass db connections to handlers here?
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func SetupDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    os.Getenv("DBNET"),
		Addr:   os.Getenv("DBADDR"),
		DBName: os.Getenv("DBNAME"),
	}

	DB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Db connection established.")

	return DB
}
