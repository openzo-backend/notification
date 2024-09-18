package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/notification/internal/models"
	"github.com/tanush-128/openzo_backend/notification/internal/service"
)

type Handler struct {
	NotificationService service.LocalNotificationService
}

func NewHandler(NotificationService *service.LocalNotificationService) *Handler {
	return &Handler{NotificationService: *NotificationService}
}
func (h *Handler) CreateNotification(ctx *gin.Context) {
	var notification models.Notification
	ctx.BindJSON(&notification)
	err := h.NotificationService.CreateNotification(ctx, notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Notification sent successfully"})

}

func (h *Handler) CreateLocalNotification(ctx *gin.Context) {
	var notification models.LocalNotification
	ctx.BindJSON(&notification)
	createdNotification, err := h.NotificationService.CreateLocalNotification(ctx, notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdNotification)

}

func (h *Handler) CreateLocalNotificationUsingTopic(ctx *gin.Context) {
	var notification models.LocalNotification
	ctx.BindJSON(&notification)
	createdNotification, err := h.NotificationService.CreateLocalNotificationUsingTopic(ctx, notification)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdNotification)

}

func (h *Handler) GetNotifications(ctx *gin.Context) {
	pincode := ctx.Param("pincode")

	Notifications, err := h.NotificationService.GetNotifications(ctx, pincode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Notifications)

}

func (h *Handler) GetNotificationByID(ctx *gin.Context) {
	id := ctx.Param("id")

	Notification, err := h.NotificationService.GetLocalNotificationByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Notification)
}

func (h *Handler) GetNotificationsByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("id")

	Notifications, err := h.NotificationService.GetLocalNotificationsByStoreID(ctx, storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Notifications)

}

func (h *Handler) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.NotificationService.DeleteLocalNotification(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
