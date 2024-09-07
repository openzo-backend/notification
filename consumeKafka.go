package main

import (
	"encoding/json"
	"log"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/tanush-128/openzo_backend/notification/internal/utils"
)

func consumeKafka() {
	conf := ReadConfig()
	topic := "notification"
	conf["group.id"] = "go-group-1"
	conf["auto.offset.reset"] = "earliest"

	var notification Notification

	for {
		consumer, err := kafka.NewConsumer(&conf)
		if err != nil {
			log.Printf("Error creating consumer: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = consumer.SubscribeTopics([]string{topic}, nil)
		if err != nil {
			log.Printf("Error subscribing to topic: %v. Retrying in 5 seconds...", err)
			consumer.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			e := consumer.Poll(1000)
			if e == nil {
				continue
			}

			switch ev := e.(type) {
			case *kafka.Message:
				err := json.Unmarshal(ev.Value, &notification)
				if err != nil {
					log.Printf("Error unmarshalling JSON: %v", err)
					continue
				}

				notifDataMap := map[string]string{}
				err = json.Unmarshal([]byte(notification.Data), &notifDataMap)
				if err != nil {
					log.Printf("Error unmarshalling notification data: %v", err)
					continue
				}
				log.Printf("Notification received: %+v", notification)

				err = utils.SendNotification(&messaging.Message{
					Notification: &messaging.Notification{
						Title: notification.Topic + " Notification",
						Body:  notification.Message,
					},
					Data:  notifDataMap,
					Token: notification.FCMToken,
				})

				if err != nil {
					log.Printf("Error sending notification: %v", err)
				}

			case kafka.Error:
				log.Printf("Kafka error: %v", ev)
				break
			}
		}

		log.Println("Consumer disconnected. Reconnecting in 5 seconds...")
		consumer.Close()
		time.Sleep(5 * time.Second)
	}
}
