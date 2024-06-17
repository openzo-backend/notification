package service

import (
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/internal/models"
	"github.com/tanush-128/openzo_backend/notification/internal/repository"
	"github.com/tanush-128/openzo_backend/notification/internal/utils"
)

type LocalNotificationService interface {

	//CRUD
	CreateLocalNotification(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error)
	GetNotifications(ctx *gin.Context, pincode string) ([]models.LocalNotification, error)
	GetLocalNotificationByID(ctx *gin.Context, id string) (models.LocalNotification, error)
	GetLocalNotificationsByStoreID(ctx *gin.Context, storeID string) ([]models.LocalNotification, error)
	DeleteLocalNotification(ctx *gin.Context, id string) error
}

type notificationService struct {
	LocalNotificationRepository repository.LocalNotificationRepository
}

func NewLocalNotificationService(LocalNotificationRepository repository.LocalNotificationRepository) LocalNotificationService {
	return &notificationService{LocalNotificationRepository: LocalNotificationRepository}
}


func (s *notificationService) CreateLocalNotification(ctx *gin.Context, req models.LocalNotification) (models.LocalNotification, error) {

	notification_tokens, err := s.LocalNotificationRepository.GetTokensByPincode(req.Pincode)
	if err != nil {
		return models.LocalNotification{}, err
	}
	go utils.SendNotificationBulk(&messaging.MulticastMessage{
		Tokens: notification_tokens,
		Notification: &messaging.Notification{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageURL,
		},
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
