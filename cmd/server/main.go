package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	lg "log"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	db "github.com/retinotopic/GoChat/server/db/postgres"
	"github.com/retinotopic/GoChat/server/logger/loggers/zerolog"
	"github.com/retinotopic/GoChat/server/router"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeylimiter"
)

func main() {
	log := zerolog.NewZerologLogger(os.Stdout)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DATABASE")
	pgSSL := os.Getenv("PG_SSLMODE")

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s ",
		pgUser, pgPassword, pgHost, pgPort, pgDB, pgSSL,
	)
	vkoptions := valkey.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")},
		RetryDelay: func(attempts int, cmd valkey.Completed, err error) time.Duration {
			lg.Println(attempts)
			if attempts >= 1000 {
				return -1
			}
		    return 1 * time.Second
		},
	}

	dbgrd := false
	if os.Getenv("REDIS_DEBUG") == "true" {
		dbgrd = true
	}
	limiter, err := valkeylimiter.NewRateLimiter(
		valkeylimiter.RateLimiterOption{
			ClientOption: vkoptions,
			Window: time.Second * 1,
			Limit: 1 })

	if err != nil {
		panic(err)
	}
	pgclient, err := db.NewPgClient(ctx, dsn, limiter, dbgrd)
	if err != nil {
		log.Fatal("db new pool:", err)
	}
	// for debugging
	FetchUser := func(w http.ResponseWriter, r *http.Request) (string, error) {
		c, err := r.Cookie("username")
		if err != nil {
			return "", err
		}
		return c.Value, nil
	}
	dbs := stdlib.OpenDBFromPool(pgclient.Pool)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("goose set dialect:", err)
	}
	if err := goose.Up(dbs, os.Getenv("MIGRATIONS_DIR")); err != nil {
		log.Fatal("goose up:", err)
	}
	if err := dbs.Close(); err != nil {
		log.Fatal("close db conn for migrations:", err)
	}
	srv := router.NewRouter("0.0.0.0:"+os.Getenv("APP_PORT"), FetchUser, pgclient, log)
	err = srv.Run(ctx)
	if err != nil {
		log.Fatal("server run:", err)
	}
}
