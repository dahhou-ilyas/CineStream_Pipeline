package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	fmt.Println("NATS producer started")
}

func Producer(ctx context.Context) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("failed to connect to NATS server", err)
	}

	defer nc.Close()

	subject := "filmsChan"

	i := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from producer")
			return
		default:
			i += 1
			message := fmt.Sprintf("message %v", i)

			err := nc.Publish(subject, []byte(message))

			if err != nil {
				log.Fatal("failed to publish message", err)
			} else {
				log.Println("published message", subject, message)
			}
		}

	}
}

func consumer(ctx context.Context) {
	nc, err := nats.Connect("localhost:4222")

	if err != nil {
		log.Fatal("failed to connect to NATS server", err)
	}

	defer nc.Close()

	fmt.Println("Connected to NATS server on port 4222")

	subject := "filmsChan"

	messages := make(chan *nats.Msg, 1000)

	subscription, err := nc.ChanSubscribe(subject, messages)

	if err != nil {
		log.Fatal("failed to subscribe to subject", err)
	}

	defer func() {
		subscription.Unsubscribe()
		close(messages)
	}()

	log.Println("Subscribed to", subject)

}
