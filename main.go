package main

import (
	"log"
	"os"
	"test-relayer/kafka_test"
	"test-relayer/relay_conf"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/joho/godotenv"
)

func main() {
	if args := os.Args[1]; args == "relay" {
		relay := relay_conf.Relay{}

		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading environment variables: %v", err)
		}
		relay.PostgresURL = "postgres://" + os.Getenv("PSQL_USR") + ":" + os.Getenv("PSQL_PWD") + "@" + "localhost:" +
			os.Getenv("PSQL_PORT") + "/" + os.Getenv("PSQL_DB") + "?sslmode=disable"

		// relay.storage = &postgresql.PostgresBackend{DatabaseURL: relay.PostgresURL}
		relay.SetStorage(&postgresql.PostgresBackend{DatabaseURL: relay.PostgresURL})

		server, err := relayer.NewServer(&relay)
		if err != nil {
			log.Fatalf("failed to create server: %v", err)
		}
		if err := server.Start("localhost", 7447); err != nil {
			log.Fatalf("server terminated: %v", err)
		}

	} else if args == "kafka" {
		kafka_test.Test_producer_consumer()
	}

}
