package main

import (
	"context"
	"time"

	"github.com/tylergdorn/prophet_exercise/pkg/server"
	"github.com/tylergdorn/prophet_exercise/pkg/storage"
	"github.com/tylergdorn/prophet_exercise/pkg/tornodes"
)

func main() {
	ctx := context.Background()
	db, err := storage.New("./test.db")
	if err != nil {
		panic(err)
	}

	fetcher := tornodes.TorNodeFetcher{
		Writer: db,
		Delay:  time.Minute * 30,
	}
	go fetcher.Start(ctx)

	err = server.StartServer(db, db, ":8080")
	if err != nil {
		panic(err)
	}
}
