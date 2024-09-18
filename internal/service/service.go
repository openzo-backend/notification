package service

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/internal/models"
	"github.com/tanush-128/openzo_backend/notification/internal/repository"
	"github.com/tanush-128/openzo_backend/notification/internal/utils"
	"google.golang.org/api/option"
)

type LocalNotificationService interface {

	//CRUD
	CreateNotification(ctx *gin.Context, req models.Notification) error
	CreateLocalNotification(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error)
	CreateLocalNotificationUsingTopic(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error)
	GetNotifications(ctx *gin.Context, pincode string) ([]models.LocalNotification, error)
	GetLocalNotificationByID(ctx *gin.Context, id string) (models.LocalNotification, error)
	GetLocalNotificationsByStoreID(ctx *gin.Context, storeID string) ([]models.LocalNotification, error)
	DeleteLocalNotification(ctx *gin.Context, id string) error
	Subscribe( pincode string) error
}

type notificationService struct {
	LocalNotificationRepository repository.LocalNotificationRepository
}

func NewLocalNotificationService(LocalNotificationRepository repository.LocalNotificationRepository) LocalNotificationService {
	return &notificationService{LocalNotificationRepository: LocalNotificationRepository}
}
func (s *notificationService) CreateNotification(ctx *gin.Context, req models.Notification) error {
	// createdNotification, err := s.LocalNotificationRepository.CreateNotification(req)
	// if err != nil {
	// 	return models.LocalNotification{}, err // Propagate error
	// }
	err := utils.SendNotification(&messaging.Message{
		Notification: &messaging.Notification{
			Title: req.Topic,
			Body:  req.Message,
		},

		Token: req.FCMToken,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *notificationService) Subscribe( pincode string) error {

	notification_tokens, err := s.LocalNotificationRepository.GetTokensByPincode(pincode)
	if err != nil {
		return err
	}

	// rm duplicate tokens

	notification_tokens = utils.RemoveDuplicates(notification_tokens)

	log.Println("notification_tokens :", len(notification_tokens))

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

	messaginClient.SubscribeToTopic(context.Background(), notification_tokens, "pincode_"+pincode)

	return nil
}

func (s *notificationService) CreateLocalNotification(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error) {

	notification_tokens, err := s.LocalNotificationRepository.GetTokensByPincode(req.Pincode)
	if err != nil {
		return models.LocalNotification{}, err
	}

	// rm duplicate tokens

	notification_tokens = utils.RemoveDuplicates(notification_tokens)

	log.Println("notification_tokens :", len(notification_tokens))

	for _, token := range notification_tokens {
		go utils.SendNotification(&messaging.Message{
			Notification: &messaging.Notification{
				Title: req.Title,
				Body:  req.Body,

				ImageURL: req.ImageURL,
			},

			Token: token,
		})

	}

	// err = utils.SendNotificationBulk(&messaging.MulticastMessage{
	// 	Tokens: notification_tokens,

	// 	Notification: &messaging.Notification{

	// 		Title:    req.Title,
	// 		Body:     req.Body,
	// 		ImageURL: req.ImageURL,
	// 	},
	// })
	// if err != nil {
	// 	return models.LocalNotification{}, err
	// }

	createdLocalNotification, err := s.LocalNotificationRepository.CreateLocalNotification(req)
	if err != nil {
		return models.LocalNotification{}, err // Propagate error
	}

	return createdLocalNotification, nil
}

func (s *notificationService) CreateLocalNotificationUsingTopic(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error) {
	absPath, _ := filepath.Abs("firebase-config.json")
	opt := option.WithCredentialsFile(absPath)
	config := &firebase.Config{ProjectID: "openzo-rt"}
	// println()
	app, err := firebase.NewApp(context.Background(),
		config, opt)
	if err != nil {
		return models.LocalNotification{}, fmt.Errorf("error initializing app: %v", err)
	}
	messaginClient, err := app.Messaging(context.Background())
	if err != nil {
		return models.LocalNotification{}, fmt.Errorf("error getting Messaging client: %v", err)
	}

	messaginClient.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageURL,
		},
		Topic: "pincode_" + req.Pincode,
	})

	createdLocalNotification, err := s.LocalNotificationRepository.CreateLocalNotification(req)
	if err != nil {
		return models.LocalNotification{}, err // Propagate error
	}

	return createdLocalNotification, nil

}

func (s *notificationService) DeleteLocalNotification(ctx *gin.Context, id string) error {
	err := s.LocalNotificationRepository.DeleteLocalNotification(id)
	if err != nil {
		return err
	}

	return nil
}
