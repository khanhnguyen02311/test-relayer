package relay_conf

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nbd-wtf/go-nostr"
)

type Relay struct {
	PostgresURL string
	storage     *postgresql.PostgresBackend
}

func (r *Relay) Name() string {
	return "TestRelay"
}

func (r *Relay) Storage(ctx context.Context) relayer.Storage {
	return r.storage
}

func (r *Relay) SetStorage(store *postgresql.PostgresBackend) {
	r.storage = store
}

func (r *Relay) Init() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading environment variables: %v", err)
	}
	r.PostgresURL = "postgres://" + os.Getenv("PSQL_USR") + ":" + os.Getenv("PSQL_PWD") + "@" + "localhost/" + os.Getenv("PSQL_DB") + "?sslmode=disable"

	go func() {
		db := r.Storage(context.TODO()).(*postgresql.PostgresBackend)

		for {
			time.Sleep(60 * time.Minute)
			db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, 0, -14).Unix()) // clear old messages after 14 days
		}
	}()

	return nil
}

func (r *Relay) AcceptEvent(ctx context.Context, evt *nostr.Event) bool {
	// reject events that have timestamps greater than 30 minutes in the future.
	if evt.CreatedAt > nostr.Now()+30*60 {
		return false
	}
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}
	return true
}

// func (r *Relay) BeforeSave(evt *nostr.Event) {
// }

// func (r *Relay) AfterSave(evt *nostr.Event) {
// 	json.NewEncoder(os.Stderr).Encode(evt)
// }
