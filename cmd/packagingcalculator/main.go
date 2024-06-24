package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"re/cmd/packagingcalculator/api"
	"re/pkg/db"
	"time"
)

func main() {
	godotenv.Load(".env")

	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	d, err := db.New(cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbName)
	if err != nil {
		panic(err)
	}

	//migrate, create db schema and tables
	err = db.MigrateDb(cfg.DbName, d.DB)
	if err != nil {
		panic(err)
	}

	m := db.NewMySQLPersistence(d)
	rc := api.NewProductController(m)
	e := api.BuildRoutes(rc)

	port := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         port,
		Handler:      e,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func() {
		log.Println("Starting server on port: ", port)
		if err = srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server Shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return

}

type config struct {
	DbUsername string
	DbPassword string
	DbHost     string
	DbName     string
	ServerPort string
}

func getConfig() (*config, error) {
	return &config{
		DbPassword: os.Getenv("MARIADB_ROOT_PASSWORD"),
		DbName:     os.Getenv("MARIADB_DATABASE"),
		DbHost:     os.Getenv("MARIADB_HOST"),
		DbUsername: "root",
		ServerPort: os.Getenv("SERVER_PORT"),
	}, nil
}
