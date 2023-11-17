package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	// connect to nats server
	nc, err := nats.Connect("nats://192.168.0.26:4222")
	if err != nil {
		panic(err)
	}

	// create jetstream context from nats connection
	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "tmp", History: 2})
	if err != nil {
		panic(err)
	}

	if _, err := kv.PutString(ctx, "foo", "hello there"); err != nil {
		panic(err)
	}

	if _, err := kv.PutString(ctx, "foo", "hello you"); err != nil {
		panic(err)
	}

	e, err := kv.Get(ctx, "foo2")
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			panic("lol")
		}
		panic(err)
	}

	fmt.Println(e.Bucket(), e.Created(), e.Key(), e.Operation(), e.Revision(), string(e.Value()))
	// fmt.Println(v.)

	entries, err := kv.History(ctx, "foo")
	if err != nil {
		panic(err)

	}

	for _, e := range entries {
		fmt.Println(e.Bucket(), e.Created(), e.Key(), e.Operation(), e.Revision(), string(e.Value()))
	}

}

// func main() {
// 	// connect to nats server
// 	nc, err := nats.Connect("nats://192.168.0.26:4222")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// create jetstream context from nats connection
// 	js, err := jetstream.New(nc)
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
// 		Name:     "foo",
// 		Subjects: []string{"foo.*"},
// 	})

// 	if _, err := js.Publish(ctx, "foo.bar", []byte("hello")); err != nil {
// 		panic(err)
// 	}

// 	// get existing stream handle
// 	// stream, err := js.Stream(ctx, "foo")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// retrieve consumer handle from a stream
// 	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
// 		Name:      "cons",
// 		AckPolicy: jetstream.AckExplicitPolicy,
// 	})
// 	// cons, err := stream.Consumer(ctx, "cons")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// consume messages from the consumer in callback
// 	cc, err := cons.Consume(func(msg jetstream.Msg) {
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("Received jetstream message: ", string(msg.Data()))
// 		msg.Ack()
// 	})
// 	defer cc.Stop()

// 	c := make(chan int)

// 	<-c
// }
