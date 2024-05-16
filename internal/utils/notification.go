package utils

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"google.golang.org/api/option"
)

func SendNotificationBulk(Message *messaging.MulticastMessage) error {
	// ...
	// ...
	// ...
	absPath, _ := filepath.Abs("utils/openzo-rt-firebase-adminsdk-u6rwj-1989861d1f.json")

	opt := option.WithCredentialsFile(absPath)
	
	// config := &firebase.Config{ProjectID: "openzo-rt"}
	// println()
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	messaginClient, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v", err)
	}

	res, err := messaginClient.SendMulticast(context.Background(), Message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	fmt.Println(res)

	// ...
	return nil

}

func SendNotification(Message *messaging.Message) error {
	// ...
	// ...
	// ...
	absPath, _ := filepath.Abs("openzo-rt-firebase-adminsdk-u6rwj-1989861d1f.json")
	opt := option.WithCredentialsFile(absPath)
	config := &firebase.Config{ProjectID: "openzo-rt"}
	// println()
	app, err := firebase.NewApp(context.Background(),
		config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	messaginClient, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v", err)
	}

	res, err := messaginClient.Send(context.Background(), Message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	fmt.Println(res)

	// ...
	return nil

}
