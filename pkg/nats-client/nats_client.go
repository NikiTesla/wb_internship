package natsclient

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
}

func NatsConnTest(url string) error {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}
	defer nc.Drain()

	nc.Publish("greet.joe", []byte("hello"))

	sub, err := nc.SubscribeSync("greet.*")
	if err != nil {
		return err
	}

	nc.Publish("greet.joe", []byte("hello joe"))
	nc.Publish("greet.pam", []byte("hello"))

	msg, err := sub.NextMsg(10 * time.Millisecond)
	if err != nil {
		return err
	}
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	msg, err = sub.NextMsg(10 * time.Millisecond)
	if err != nil {
		return err
	}
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	nc.Publish("greet.bob", []byte("hello"))

	msg, err = sub.NextMsg(10 * time.Millisecond)
	if err != nil {
		return err
	}
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	return err
}
