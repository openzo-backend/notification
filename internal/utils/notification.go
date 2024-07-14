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

	absPath, _ := filepath.Abs("firebase-config.json")

	opt := option.WithCredentialsFile(absPath)

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

	return nil

}

func SendNotification(Message *messaging.Message) error {
	// ...
	// ...
	// ...
	absPath, _ := filepath.Abs("firebase-config.json")
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

	_, err = messaginClient.Send(context.Background(), Message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	// ...
	return nil

}

func RemoveDuplicates(elements []string) []string {

	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
