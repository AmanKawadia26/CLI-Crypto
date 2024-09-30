package Handlers

import (
	"cryptotracker/REST-API/errors"
	"cryptotracker/REST-API/response"
	"cryptotracker/internal/services"
	"cryptotracker/pkg/logger"
	"net/http"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	notificationService services.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler instance
func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// CheckNotificationHandler handles the request to check notifications
func (h *NotificationHandler) CheckNotificationHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		logger.Logger.Error("Unauthorized access")
		err := errs.NewUnauthorizedAccessError()
		err.ToJSON(w)
		return
	}

	notifications, err := h.notificationService.CheckNotification(username)
	if err != nil {
		logger.Logger.Error("Failed to check notifications", err)
		appErr := errs.NewFailedToCheckNotificationsError()
		appErr.ToJSON(w)
		return
	}

	// Format notifications into a list of objects
	notificationList := make([]map[string]interface{}, len(notifications))
	for i, notification := range notifications {
		notificationList[i] = map[string]interface{}{
			"index":   i + 1, // Indexing starts at 1 for better readability
			"message": notification.Message,
		}
	}

	logger.Logger.Info("Notifications retrieved")
	response.SendJSONResponse(w, http.StatusOK, "success", "Notifications retrieved", notificationList, "")
}
