package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllowList(t *testing.T) {
	db, err := New(":memory:")
	ctx := context.Background()
	assert.NoError(t, err)
	now := time.Now()
	allowListOne := AllowListIP{IP: "1.1.1.1", TimeAdded: now, Note: "good"}
	allowListTwo := AllowListIP{IP: "2.2.2.2", TimeAdded: now, Note: "good"}
	assert.NoError(t, db.InsertAllowIP(ctx, allowListOne))
	assert.NoError(t, db.InsertAllowIP(ctx, allowListTwo))
	list, err := db.GetAllowListIPs(ctx)
	assert.NoError(t, err)
	for _, item := range list {
		if item.IP != "1.1.1.1" && item.IP != "2.2.2.2" {
			t.Errorf("got %s, not 1.1.1.1 or 2.2.2.2", item.IP)
		}
		assert.Equal(t, now.Format(time.DateTime), item.TimeAdded.Format(time.DateTime))
		assert.Equal(t, item.Note, "good")
	}
	db.DeleteAllowIP(ctx, "1.1.1.1")
	list, err = db.GetAllowListIPs(ctx)
	assert.NoError(t, err)
	for _, item := range list {
		if item.IP != "2.2.2.2" {
			t.Errorf("got %s, not 2.2.2.2", item.IP)
		}
		assert.Equal(t, now.Format(time.DateTime), item.TimeAdded.Format(time.DateTime))
		assert.Equal(t, item.Note, "good")
	}

}

func TestAllowListNodes(t *testing.T) {
	db, err := New(":memory:")
	ctx := context.Background()
	assert.NoError(t, err)
	now := time.Now()
	err = db.InsertAllowIP(ctx, AllowListIP{IP: "1.1.1.1", TimeAdded: now})
	assert.NoError(t, err)
	err = db.InsertBatch(ctx, []TorExitNode{{IP: "1.1.1.1"}, {IP: "2.2.2.2"}})
	assert.NoError(t, err)
	nodes, err := db.GetTorNodes(ctx)
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].IP, "2.2.2.2")
}
