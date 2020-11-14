package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
)

func main() {
	fmt.Println("Publishing...")
	w := os.Stdout
	projectID := "messagedelivery-1605371057361"
	topicID := "MessageDeliveryTopic"
	msg := "Foobar"

	err := publish(w, projectID, topicID, msg)

	if err != nil {
		log.Fatal(err)
	}

	subID := "ClientSubscription"
	fmt.Println("Subscribing...")
	errSub := pullMsgs(w, projectID, subID)

	if errSub != nil {
		log.Fatal(errSub)
	}
}

func publish(w io.Writer, projectID, topicID, msg string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}

func pullMsgs(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
		received++
		if received >= 10 {
			fmt.Println("closing")
			cancel()
		}
	})
	if err != nil {
		fmt.Println("6")
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
