package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/notification/internal/models"

	"gorm.io/gorm"
)

type LocalNotificationRepository interface {
	CreateLocalNotification(LocalNotification models.LocalNotification) (models.LocalNotification, error)
	GetLocalNotificationByID(id string) (models.LocalNotification, error)
	GetLocalNotifications(pincode string) ([]models.LocalNotification, error)
	GetLocalNotificationsByStoreID(id string) ([]models.LocalNotification, error)

	GetTokensByPincode(pincode string) ([]string, error)

	DeleteLocalNotification(id string) error
	// Add more methods for other LocalNotification operations (GetLocalNotificationByEmail, UpdateLocalNotification, etc.)

}

type notificationRepository struct {
	db *gorm.DB
}

func NewLocalNotificationRepository(db *gorm.DB) LocalNotificationRepository {

	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateLocalNotification(LocalNotification models.LocalNotification) (models.LocalNotification, error) {
	LocalNotification.ID = uuid.New().String()
	tx := r.db.Create(&LocalNotification)

	if tx.Error != nil {
		return models.LocalNotification{}, tx.Error
	}

	return LocalNotification, nil
}

func (r *notificationRepository) GetLocalNotifications(pincode string) ([]models.LocalNotification, error) {
	var LocalNotifications []models.LocalNotification
	tx := r.db.Where("pincode = ?", pincode).Find(&LocalNotifications)
	if tx.Error != nil {
		return []models.LocalNotification{}, tx.Error
	}

	return LocalNotifications, nil
}

func (r *notificationRepository) GetLocalNotificationByID(id string) (models.LocalNotification, error) {
	var LocalNotification models.LocalNotification
	tx := r.db.Preload("Images").Where("id = ?", id).First(&LocalNotification)
	if tx.Error != nil {
		return models.LocalNotification{}, tx.Error
	}

	return LocalNotification, nil
}

func (r *notificationRepository) GetLocalNotificationsByStoreID(id string) ([]models.LocalNotification, error) {
	var LocalNotifications []models.LocalNotification
	tx := r.db.Preload("Images").Where("store_id = ?", id).Find(&LocalNotifications)
	if tx.Error != nil {
		return []models.LocalNotification{}, tx.Error
	}

	return LocalNotifications, nil
}

func (r *notificationRepository) DeleteLocalNotification(id string) error {
	var LocalNotification models.LocalNotification
	tx := r.db.Where("id = ?", id).Delete(&LocalNotification)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *notificationRepository) GetTokensByPincode(pincode string) ([]string, error) {
	var notification_tokens []string
	// UserData had pincode 	and token
	
	r.db.Table("users").Where("pincode = ?", pincode).Pluck("notification_token", &notification_tokens)

	return notification_tokens, nil
}

// Implement other repository methods (GetLocalNotificationByID, GetLocalNotificationByEmail, UpdateLocalNotification, etc.) with proper error handling
