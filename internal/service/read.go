package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/internal/models"
)

func (s *notificationService) GetNotifications(ctx *gin.Context, pincode string) ([]models.LocalNotification, error) {
	Notifications, err := s.LocalNotificationRepository.GetLocalNotifications(pincode)
	if err != nil {
		return []models.LocalNotification{}, err
	}

	return Notifications, nil

}

func (s *notificationService) GetLocalNotificationByID(ctx *gin.Context, id string) (models.LocalNotification, error) {
	LocalNotification, err := s.LocalNotificationRepository.GetLocalNotificationByID(id)
	if err != nil {
		return models.LocalNotification{}, err
	}

	return LocalNotification, nil
}

func (s *notificationService) GetLocalNotificationsByStoreID(ctx *gin.Context, storeID string) ([]models.LocalNotification, error) {
	LocalNotifications, err := s.LocalNotificationRepository.GetLocalNotificationsByStoreID(storeID)
	if err != nil {
		return []models.LocalNotification{}, err
	}

	return LocalNotifications, nil
}
