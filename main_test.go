package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

type User struct {
	ID    string
	Name  string
	Email string
	Msg   string
}

func TestKafka(t *testing.T) {
	kafkaBroker := []string{"localhost:9094", "localhost:9097", "localhost:9100"}
	topic := "my-topic"
	groupID := "my-topic-001"

	t.Run("check leader kafka", func(t *testing.T) {
		defer func() {
			r := recover()
			if r != nil {
				log.Fatal(r)
			}
		}()
		conn, err := kafka.Dial("tcp", kafkaBroker[2])
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
		controller, err := conn.Controller()
		if err != nil {
			panic(err.Error())
		}
		var connLeader *kafka.Conn
		connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			panic(err.Error())
		}
		defer connLeader.Close()
		log.Print(connLeader.Broker().Host)
		log.Print(connLeader.Broker().ID)
		log.Print(connLeader.Broker().Port)
		log.Print(connLeader.Broker().Rack)
	})

	t.Run("create topic", func(t *testing.T) {
		defer func() {
			r := recover()
			if r != nil {
				log.Fatal(r)
			}
		}()
		conn, err := kafka.Dial("tcp", kafkaBroker[2])
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
		controller, err := conn.Controller()
		if err != nil {
			panic(err.Error())
		}
		var connLeader *kafka.Conn
		connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			panic(err.Error())
		}
		defer connLeader.Close()
		log.Print(connLeader.Broker().Host)
		log.Print(connLeader.Broker().ID)
		log.Print(connLeader.Broker().Port)
		log.Print(connLeader.Broker().Rack)

		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     -1,
				ReplicationFactor: -1,
			},
		}

		err = connLeader.CreateTopics(topicConfigs...)
		if err != nil {
			panic(err.Error())
		}
	})

	t.Run("list topic", func(t *testing.T) {
		conn, err := kafka.Dial("tcp", kafkaBroker[1])
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
		partitions, err := conn.ReadPartitions()
		if err != nil {
			panic(err.Error())
		}
		m := map[string]struct{}{}
		for _, p := range partitions {
			m[p.Topic] = struct{}{}
		}
		for k := range m {
			fmt.Println(k)
		}
	})

	t.Run("create message", func(t *testing.T) {
		partition := 0

		conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker[2], topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		user := &User{
			ID:    uuid.New().String(),
			Name:  "test user",
			Email: "test@example.com",
			Msg:   "hello world 2",
		}
		userByte, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		msg := kafka.Message{
			Key:   []byte(uuid.New().String()),
			Value: userByte,
			Headers: []protocol.Header{
				{
					Key:   "ID",
					Value: []byte("12345"),
				},
			},
		}
		_, err = conn.WriteMessages(msg)
		if err != nil {
			log.Fatalln(err)
		}
		log.Print("success send message : ", string(msg.Value))
	})

	t.Run("read message", func(t *testing.T) {
		partition := 0

		conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker[1], topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}
		defer conn.Close()

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   kafkaBroker,
			Topic:     topic,
			Partition: partition,
			GroupID:   groupID,
		})
		defer reader.Close()
		fmt.Println("start consuming ... !!")
		msg := kafka.Message{
			Headers: []protocol.Header{},
		}
		for {
			msg, err = reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			// Process the message
			log.Print("topic : " + msg.Topic)
			log.Print(fmt.Sprintf("partition : %v", msg.Partition))
			log.Print(fmt.Sprintf("offset : %v", msg.Offset))
			log.Print(fmt.Sprintf("High Water Mark : %v", msg.HighWaterMark))
			log.Print("key : " + string(msg.Key))
			log.Print("value : " + string(msg.Value))
			log.Print("time : " + msg.Time.String())
			log.Print(fmt.Sprintf("header : %v", msg.Headers))

			// Commit the offset to mark the message as processed
			err = reader.CommitMessages(context.Background(), msg)
			if err != nil {
				log.Fatal("failed to commit offset:", err)
			}
		}
	})
}
