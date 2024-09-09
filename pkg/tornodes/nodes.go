package tornodes

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/tylergdorn/prophet_exercise/pkg/storage"
)

type ExitNodeList []storage.TorExitNode

type TorNodeWriter interface {
	InsertBatch(ctx context.Context, nodes []storage.TorExitNode) error
}

type TorNodeDataSource struct {
	URL  string
	Name string
}

type TorNodeFetcher struct {
	Writer  TorNodeWriter
	Delay   time.Duration
	Sources []TorNodeDataSource
}

func (t *TorNodeFetcher) Start(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			break
		}
		slog.DebugContext(ctx, fmt.Sprintf("sleeping %s until fetching", t.Delay))
		time.Sleep(t.Delay)
		// add error handling here
		slog.DebugContext(ctx, "fetching lists")
		for _, source := range t.Sources {
			slog.DebugContext(ctx, fmt.Sprintf("fetching %s list", source.Name))
			list, err := FetchList(source.URL, source.Name)
			if err != nil {
				slog.ErrorContext(ctx, "fetching", "error", err)
			}
			if err = t.Writer.InsertBatch(ctx, list); err != nil {
				slog.ErrorContext(ctx, "inserting", "error", err)
			}
		}
	}
}

func FetchList(sourceURL string, sourceName string) (ExitNodeList, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return nil, err
	}
	listBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	ips := strings.Split(string(listBytes), "\n")
	nodes := ExitNodeList{}
	for _, ip := range ips {
		parsedIp := net.ParseIP(ip)
		if parsedIp != nil {
			nodes = append(nodes, storage.TorExitNode{
				IP:          parsedIp.String(),
				LastSeen:    now,
				FetchedFrom: sourceName,
			})
		}
	}
	return nodes, nil

}
