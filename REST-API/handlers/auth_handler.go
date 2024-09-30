package Handlers

import (
	errs "cryptotracker/REST-API/errors"
	"cryptotracker/REST-API/middleware"
	"cryptotracker/REST-API/response"
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"cryptotracker/pkg/logger" // Import the logger package
	"cryptotracker/pkg/utils"
	"cryptotracker/pkg/validation"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Mobile   int    `json:"mobile"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		logger.Logger.Error("Failed to decode login request", err)
		errs.NewInvalidRequestPayloadError().ToJSON(w)
		return
	}

	logger.Logger.Info("Login request received for user: ", loginReq.Username)

	user, _, err := h.authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		logger.Logger.Error("Login failed for user: ", loginReq.Username, err)
		errs.NewInvalidCredentialsError().ToJSON(w)
		return
	}

	middlewareVar := middleware.JWTTokenService{}

	token, err := middlewareVar.GenerateToken(loginReq.Username, user.Role)
	if err != nil {
		logger.Logger.Error("Error generating token for user: ", loginReq.Username, err)
		errs.NewTokenGenerationFailedError().ToJSON(w)
		return
	}

	logger.Logger.Info("Login successful for user: ", user.Username)
	response.SendJSONResponse(w, http.StatusOK, "success", "Login successful", nil, token)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	middlewareVar := middleware.JWTTokenService{}
	tokenString := middlewareVar.ExtractToken(r)
	if tokenString == "" {
		logger.Logger.Warn("Authorization token required for logout", "method", r.Method, "path", r.URL.Path)
		errs.NewTokenRequiredError().ToJSON(w)
		return
	}

	claims, err := middlewareVar.ParseToken(tokenString)
	if err != nil {
		logger.Logger.Warn("Invalid token during logout", "error", err)
		errs.NewInvalidTokenError().ToJSON(w)
		return
	}

	logger.Logger.Info("Logout request received", "username", claims.Username, "method", r.Method, "path", r.URL.Path)

	middlewareVar.BlacklistToken(tokenString)

	response.SendJSONResponse(w, http.StatusOK, "success", "Logout successful", nil, "")
	logger.Logger.Info("Logout successful for user", "username", claims.Username)
}

func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signupReq SignupRequest

	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		logger.Logger.Error("Failed to decode signup request", err)
		errs.NewInvalidRequestPayloadError().ToJSON(w)
		return
	}

	logger.Logger.Info("Signup request received for user: ", signupReq.Username)

	if !validation.IsValidUsername(signupReq.Username) {
		logger.Logger.Warn("Invalid username during signup: ", signupReq.Username)
		errs.NewInvalidUsernameError().ToJSON(w)
		return
	}

	if !validation.IsValidPassword(signupReq.Password) {
		logger.Logger.Warn("Invalid password during signup for user: ", signupReq.Username)
		errs.NewInvalidPasswordError().ToJSON(w)
		return
	}

	if !validation.IsValidEmail(signupReq.Email) {
		logger.Logger.Warn("Invalid email during signup for user: ", signupReq.Username)
		errs.NewInvalidEmailError().ToJSON(w)
		return
	}

	if !validation.IsValidMobile(signupReq.Mobile) {
		logger.Logger.Warn("Invalid mobile number during signup for user: ", signupReq.Username)
		errs.NewInvalidMobileError().ToJSON(w)
		return
	}

	hashedPassword := utils.HashPassword(signupReq.Password)

	user := &models.User{
		Username: signupReq.Username,
		Password: hashedPassword,
		Email:    signupReq.Email,
		Mobile:   signupReq.Mobile,
		IsAdmin:  false,
		Role:     "user",
	}

	err = h.authService.Signup(user)
	if err != nil {
		logger.Logger.Error("Signup failed for user: ", signupReq.Username, err)
		errs.NewSignupFailedError().ToJSON(w)
		return
	}

	logger.Logger.Info("Signup successful for user: ", user.Username)
	response.SendJSONResponse(w, http.StatusOK, "success", "Signup successful", nil, "")
}
