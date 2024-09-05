package tornodes

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/tylergdorn/prophet_exercise/pkg/storage"
)

// func FetchDanMe() ([]byte, error) {
// 	resp, err := http.Get("https://www.dan.me.uk/torlist/?exit")
// 	if err != nil {
// 		return nil, err
// 	}
// 	byt, err := httputil.DumpResponse(resp, false)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(string(byt))

// 	file, err := os.OpenFile(fmt.Sprintf("./dan_md_tornodes-%s", time.Now().Format(time.DateOnly)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	byts, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = file.Write(byts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return byts, nil
// }

type ExitNodeList []storage.TorExitNode

type TorNodeWriter interface {
	InsertBatch(ctx context.Context, nodes []storage.TorExitNode)
}

type TorNodeFetcher struct {
	Writer TorNodeWriter
	Delay  time.Duration
}

func (t *TorNodeFetcher) Start(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			break
		}
		time.Sleep(t.Delay)
		// add error handling here
		err := t.backgroundProcess(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "running background process", "error", err)
		}
	}
}

func (t *TorNodeFetcher) backgroundProcess(ctx context.Context) error {
	list, err := FetchList("https://check.torproject.org/torbulkexitlist", "torproject")
	if err != nil {
		return err
	}
	t.Writer.InsertBatch(ctx, list)
	list, err = FetchList("https://www.dan.me.uk/torlist/?exit", "danmeuk")
	if err != nil {
		return err
	}
	t.Writer.InsertBatch(ctx, list)
	return nil
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

// func FetchTorExitNodes() (ExitNodeList, error) {
// 	resp, err := http.Get("https://check.torproject.org/torbulkexitlist")
// 	if err != nil {
// 		return nil, err
// 	}
// 	listBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	now := time.Now()
// 	ips := strings.Split(string(listBytes), "\n")
// 	nodes := ExitNodeList{}
// 	for _, ip := range ips {
// 		parsedIp := net.ParseIP(ip)
// 		if parsedIp != nil {
// 			nodes = append(nodes, storage.TorExitNode{
// 				IP:          parsedIp.String(),
// 				LastSeen:    now,
// 				FetchedFrom: "torproject",
// 			})
// 		}
// 	}
// 	return nodes, nil
// }
