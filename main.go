package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/pubsub"
)

func main() {
	fmt.Println("Publishing...")

	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)

	projectID := "messagedelivery-1605371057361"
	topicID := "MessageDeliveryTopic"
	msg := "Please work"
	err := publish(w, projectID, topicID, msg)

	if err != nil {
		log.Fatal(err)
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
