package Handlers

import (
	"cryptotracker/REST-API/errors"
	"cryptotracker/REST-API/response"
	"cryptotracker/internal/services"
	"cryptotracker/pkg/logger"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserServices
}

func NewUserHandler(userService services.UserServices) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
	// Get the username from the request context (set by the middleware)
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		logger.Logger.Error("Unauthorized access attempt")
		err := errs.NewUnauthorizedAccessError()
		err.ToJSON(w)
		return
	}

	// Fetch user profile from the service layer
	userProfile, err := h.userService.GetUserProfile(username)
	if err != nil {
		logger.Logger.Error("Error retrieving user profile", err)
		appErr := errs.NewFailedToRetrieveUserProfileError()
		appErr.ToJSON(w)
		return
	}

	userData := map[string]string{
		"username": userProfile.Username,
		"email":    userProfile.Email,
		"mobile":   strconv.Itoa(userProfile.Mobile),
	}

	// Send the response with the user profile
	logger.Logger.Info("User profile retrieved successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "User profile retrieved successfully", userData, "")
}
