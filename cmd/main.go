package main

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tylergdorn/prophet_exercise/pkg/server"
	"github.com/tylergdorn/prophet_exercise/pkg/storage"
	"github.com/tylergdorn/prophet_exercise/pkg/tornodes"
)

func main() {
	dbPath := flag.String("dbpath", "./tor.db", "path to use for sqlite db")
	port := flag.String("port", "8080", "port to listen on")
	prodMode := flag.Bool("production", false, "set to true to make server production mode")
	flag.Parse()

	ctx := context.Background()
	db, err := storage.New(*dbPath)
	if err != nil {
		panic(err)
	}
	slog.SetLogLoggerLevel(slog.LevelDebug)

	fetcher := tornodes.TorNodeFetcher{
		Writer: db,
		Delay:  time.Minute * 30,
		Sources: []tornodes.TorNodeDataSource{
			{URL: "https://check.torproject.org/torbulkexitlist", Name: "torproject"},
			{URL: "https://www.dan.me.uk/torlist/?exit", Name: "danmeuk"},
		},
	}
	go fetcher.Start(ctx)

	if *prodMode {
		gin.SetMode("release")
	}
	err = server.StartServer(db, db, ":"+*port)
	if err != nil {
		panic(err)
	}
}
