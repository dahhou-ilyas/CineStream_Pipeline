package natsProducers

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go-films-pipline/cleaner"
	"log"
)

func Producer(ctx context.Context, films cleaner.MovieEnriched) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("failed to connect to NATS server", err)
	}

	defer nc.Close()

	subject := "filmsChan"

	select {
	case <-ctx.Done():
		log.Println("exiting from producer")
		return
	default:

		byteSlice, err := json.Marshal(films)
		if err != nil {
			log.Fatal("failed to marshal films to json", err)
		}

		errs := nc.Publish(subject, byteSlice)

		if errs != nil {
			log.Fatal("failed to publish message", err)
		} else {
			log.Println("published message", subject, byteSlice)
		}
	}

}
