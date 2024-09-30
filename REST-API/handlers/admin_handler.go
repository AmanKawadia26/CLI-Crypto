package Handlers

import (
	errs "cryptotracker/REST-API/errors"
	"cryptotracker/REST-API/response"
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"cryptotracker/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) Profiles(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("Profiles request received")

	username := r.URL.Query().Get("username")

	if username != "" {
		logger.Logger.Info("Fetching profile for user: ", username)

		users, err := h.adminService.ViewUserProfiles()
		if err != nil {
			logger.Logger.Error("Error fetching user profiles", err)
			errs.NewFailedToFetchProfilesError().ToJSON(w)
			return
		}

		var specificUser *models.User
		for _, user := range users {
			if user.Username == username {
				specificUser = user
				break
			}
		}

		if specificUser == nil {
			logger.Logger.Warn("User not found: ", username)
			errs.NewUserNotFoundError().ToJSON(w)
			return
		}

		logger.Logger.Info("User profile fetched successfully")
		response.SendJSONResponse(w, http.StatusOK, "success", "User profile fetched successfully", specificUser, "")
		return
	}

	users, err := h.adminService.ViewUserProfiles()
	if err != nil {
		logger.Logger.Error("Error fetching user profiles", err)
		errs.NewFailedToFetchProfilesError().ToJSON(w)
		return
	}

	logger.Logger.Info("User profiles fetched successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "User profiles fetched successfully", users, "")
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		logger.Logger.Warn("Username is missing in delete user request")
		errs.NewInvalidUsernameError().ToJSON(w)
		return
	}

	logger.Logger.Info("Delete user request received for: ", username)

	err := h.adminService.DeleteUser(username)
	if err != nil {
		logger.Logger.Error("Error deleting user: ", username, err)
		errs.NewFailedToDeleteUserError().ToJSON(w)
		return
	}

	logger.Logger.Info("User deleted successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "User deleted successfully", nil, "")
}

func (h *AdminHandler) DelegateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		logger.Logger.Warn("Username is missing in delegate user request")
		errs.NewInvalidUsernameError().ToJSON(w)
		return
	}

	logger.Logger.Info("Delegate user request received for: ", username)

	err := h.adminService.ChangeUserStatus(username)
	if err != nil {
		logger.Logger.Error("Error delegating user to admin: ", username, err)
		errs.NewFailedToDelegateUserError().ToJSON(w)
		return
	}

	logger.Logger.Info("User delegated to admin successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "User delegated to admin successfully", nil, "")
}

func (h *AdminHandler) UnavailableCryptoRequests(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("Fetching unavailable crypto requests")

	cryptoSymbol := r.URL.Query().Get("crypto")
	var err error
	var result map[string]interface{}

	if cryptoSymbol == "" {
		result, err = h.groupCryptoRequests()
	} else {
		result, err = h.getSpecificCryptoRequests(cryptoSymbol)
	}

	if err != nil {
		logger.Logger.Error("Error fetching unavailable crypto requests", err)
		errs.NewFailedToFetchRequestsError().ToJSON(w)
		return
	}

	logger.Logger.Info("Requests fetched successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Requests fetched successfully", result, "")
}

// ... (keep the helper functions getSpecificCryptoRequests and groupCryptoRequests)

func (h *AdminHandler) SpecificUserUnavailableCryptoRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		logger.Logger.Warn("Username is missing in request")
		errs.NewInvalidUsernameError().ToJSON(w)
		return
	}

	logger.Logger.Info("Fetching unavailable crypto requests for user: ", username)

	requests, err := h.adminService.ManageUserRequests()
	if err != nil {
		logger.Logger.Error("Error fetching specific user unavailable crypto requests", err)
		errs.NewFailedToFetchRequestsError().ToJSON(w)
		return
	}

	var userRequests []*models.UnavailableCryptoRequest
	for _, req := range requests {
		if req.UserName == username {
			userRequests = append(userRequests, req)
		}
	}

	if len(userRequests) == 0 {
		response.SendJSONResponse(w, http.StatusOK, "success", "No requests found for user", nil, "")
		return
	}

	logger.Logger.Info("Requests fetched successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Requests fetched successfully", userRequests, "")
}

func (h *AdminHandler) ActOnUnavailableCryptoRequestsBySymbol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	crypto := vars["crypto"]
	if crypto == "" {
		logger.Logger.Warn("Crypto Symbol is missing in request")
		errs.NewInvalidCryptoSymbolError().ToJSON(w)
		return
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	if status == "" {
		logger.Logger.Warn("Status is missing in request")
		errs.NewInvalidRequestStatusError().ToJSON(w)
		return
	}

	requests, err := h.adminService.ManageSpecificCryptoRequests(crypto)
	if err != nil {
		logger.Logger.Error("Error fetching user unavailable crypto requests", err)
		errs.NewFailedToFetchRequestsError().ToJSON(w)
		return
	}

	logger.Logger.Info("Acting on request for crypto: ", crypto)

	err = h.adminService.UpdateRequestStatus(requests, status)
	if err != nil {
		logger.Logger.Error("Error updating user unavailable crypto requests", err)
		errs.NewFailedToUpdateRequestsError().ToJSON(w)
		return
	}

	logger.Logger.Info("Acting on request of ", crypto, " done successfully")
	response.SendJSONResponse(w, http.StatusOK, "success", "Acted on Requests successfully", nil, "")
}

// Helper function to handle specific crypto requests
func (h *AdminHandler) getSpecificCryptoRequests(cryptoSymbol string) (map[string]interface{}, error) {
	requests, err := h.adminService.ManageSpecificCryptoRequests(cryptoSymbol)
	if err != nil {
		return nil, err
	}

	responseData := map[string]interface{}{
		"crypto_symbol": cryptoSymbol,
		"count":         len(requests),
		"requests":      requests,
	}

	return responseData, nil
}

// Helper function to group crypto requests by symbol
func (h *AdminHandler) groupCryptoRequests() (map[string]interface{}, error) {
	requests, err := h.adminService.ManageUserRequests()
	if err != nil {
		return nil, err
	}

	groupedRequests := make(map[string][]*models.UnavailableCryptoRequest)
	for _, request := range requests {
		cryptoSymbol := request.CryptoSymbol
		groupedRequests[cryptoSymbol] = append(groupedRequests[cryptoSymbol], request)
	}

	// Create the final response
	responseData := make(map[string]interface{})
	responseData["count"] = len(requests)
	cryptoRequests := make([]map[string]interface{}, 0)

	for cryptoSymbol, reqs := range groupedRequests {
		cryptoReq := map[string]interface{}{
			"crypto_symbol": cryptoSymbol,
			"count":         len(reqs),
			"requests":      reqs, // Subarray of requests for this crypto symbol
		}
		cryptoRequests = append(cryptoRequests, cryptoReq)
	}

	responseData["data"] = cryptoRequests
	return responseData, nil
}
