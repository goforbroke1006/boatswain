package cmd

// go:generate go install github.com/golang/mock/mockgen@v1.6.0
// go:generate mockgen -source=./../../../domain/block.go -package=mocks -destination=mocks/block.mock.go

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"

	"github.com/goforbroke1006/boatswain/domain"
	"github.com/goforbroke1006/boatswain/internal/common"
	"github.com/goforbroke1006/boatswain/internal/storage"
)

func TestNodeReconciliation(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration tests")
	}

	compose, err := tc.NewDockerCompose("testdata/docker-compose.yaml")
	assert.NoError(t, err, "NewDockerComposeAPI()")
	t.Cleanup(func() {
		downErr := compose.Down(context.Background(),
			tc.RemoveOrphans(true), tc.RemoveImagesLocal)
		assert.NoError(t, downErr, "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	currentState := []*domain.Block{domain.Genesis}
	currentState = append(currentState, domain.NewBlock(2, currentState[len(currentState)-1].Hash, time.Now(), nil))
	currentState = append(currentState, domain.NewBlock(3, currentState[len(currentState)-1].Hash, time.Now(), nil))
	currentState = append(currentState, domain.NewBlock(4, currentState[len(currentState)-1].Hash, time.Now(), nil))
	currentState = append(currentState, domain.NewBlock(5, currentState[len(currentState)-1].Hash, time.Now(), nil))
	currentState = append(currentState, domain.NewBlock(6, currentState[len(currentState)-1].Hash, time.Now(), nil))

	// TODO: init and fill ./testdata/chat-blocks-node001.db
	// TODO: init and fill ./testdata/chat-blocks-node002.db
	// TODO: init and clear ./testdata/chat-blocks-node003.db
	db1, db1Err := common.OpenDBConn("./testdata/chat-blocks-node001.db")
	assert.NoError(t, db1Err)
	mig1Err := common.ApplyMigrationFile(db1, "./../../../db/schema.sql")
	assert.NoError(t, mig1Err)

	db2, db2Err := common.OpenDBConn("./testdata/chat-blocks-node002.db")
	assert.NoError(t, db2Err)
	mig2Err := common.ApplyMigrationFile(db2, "./../../../db/schema.sql")
	assert.NoError(t, mig2Err)

	db3, db3Err := common.OpenDBConn("./testdata/chat-blocks-node003.db")
	assert.NoError(t, db3Err)
	mig3Err := common.ApplyMigrationFile(db3, "./../../../db/schema.sql")
	assert.NoError(t, mig3Err)

	storage1 := storage.NewBlockStorage(db1)
	storage2 := storage.NewBlockStorage(db2)
	storage3 := storage.NewBlockStorage(db3)

	_ = storage1.Store(ctx, currentState...)                      // fill
	_ = storage2.Store(ctx, currentState...)                      // fill
	_, _ = db3.ExecContext(ctx, `DELETE FROM blocks WHERE TRUE;`) // clear

	// check DBs ready
	{
		count1, _ := storage1.GetCount(ctx)
		count2, _ := storage2.GetCount(ctx)
		count3, _ := storage3.GetCount(ctx)

		assert.Equal(t, len(currentState), int(count1))
		assert.Equal(t, len(currentState), int(count2))
		assert.Equal(t, 0, int(count3))
	}

	upErr := compose.Up(ctx, tc.Wait(true))
	assert.NoError(t, upErr, "compose.Up()")

	// wait for node-003 becomes READYz
	waitReady(t, "http://localhost:48083/readyz", 10, 2*time.Second)

	// ensure ./testdata/chat-blocks-node003.db has all blocks
	{
		count1, _ := storage1.GetCount(ctx)
		count2, _ := storage2.GetCount(ctx)
		count3, _ := storage3.GetCount(ctx)

		assert.Equal(t, len(currentState), int(count1))
		assert.Equal(t, len(currentState), int(count2))
		assert.Equal(t, len(currentState), int(count3))
	}
}

func waitReady(t *testing.T, addr string, retry uint, interval time.Duration) {
	for idx := uint(0); idx < retry; idx++ {
		resp, err := http.Get(addr)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode == http.StatusOK {
			t.Log(addr, resp.StatusCode)
			return
		}

		time.Sleep(interval)
	}
}
