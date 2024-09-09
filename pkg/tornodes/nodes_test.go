package tornodes

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tylergdorn/prophet_exercise/pkg/storage"
)

func makeTestServer(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(data)
		if err != nil {
			panic(err)
		}

	}
}

func TestFetchList(t *testing.T) {
	testmux := http.NewServeMux()
	testmux.HandleFunc("/tornodes", makeTestServer("./testdata/tornodes"))
	testmux.HandleFunc("/torproject", makeTestServer("./testdata/torproject"))

	testserv := httptest.NewServer(testmux)
	defer testserv.Close()

	// dan.uk.me/tornodes input
	res, err := FetchList(testserv.URL+"/tornodes", "tornodes")
	if err != nil {
		t.Errorf("fetching: %v", err)
	}
	assert.Equal(t, res[0].FetchedFrom, "tornodes")
	assert.Len(t, res, 1838)
	assert.Contains(t, mapOutIPs(res), "102.130.113.9")

	// torproject input
	res, err = FetchList(testserv.URL+"/torproject", "torproject")
	if err != nil {
		t.Errorf("fetching: %v", err)
	}
	assert.Equal(t, res[0].FetchedFrom, "torproject")
	assert.Len(t, res, 1158)
	assert.Contains(t, mapOutIPs(res), "102.130.113.9")

}

type FakeWriter struct {
	received []storage.TorExitNode
}

func (fw *FakeWriter) InsertBatch(ctx context.Context, nodes []storage.TorExitNode) error {
	fw.received = append(fw.received, nodes...)
	return nil
}

func mapOutIPs(input ExitNodeList) (res []string) {
	for _, item := range input {
		res = append(res, item.IP)
	}
	return res
}

func TestFetcher(t *testing.T) {
	testmux := http.NewServeMux()
	testmux.HandleFunc("/tornodes", makeTestServer("./testdata/tornodes_minimal"))
	testmux.HandleFunc("/torproject", makeTestServer("./testdata/torproject_minimal"))

	testserv := httptest.NewServer(testmux)
	defer testserv.Close()

	fw := FakeWriter{}
	fetcher := TorNodeFetcher{
		Writer: &fw,
		Delay:  time.Second * 2,
		Sources: []TorNodeDataSource{
			{URL: testserv.URL + "/tornodes", Name: "tornodes"},
			{URL: testserv.URL + "/torproject", Name: "torproject"},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	go fetcher.Start(ctx)
	time.Sleep(time.Second * 3)
	cancel()
	assert.ElementsMatch(t, []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5", "6.6.6.6"}, mapOutIPs(fw.received))
}
