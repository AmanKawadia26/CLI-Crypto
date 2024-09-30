package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidUsername             = errors.New("invalid username")
	ErrFailedToFetchProfiles       = errors.New("failed to fetch user profiles")
	ErrFailedToDeleteUser          = errors.New("failed to delete user")
	ErrFailedToDelegateUser        = errors.New("failed to delegate user")
	ErrFailedToFetchRequests       = errors.New("failed to fetch requests")
	ErrFailedToUpdateRequests      = errors.New("failed to update requests")
	ErrInvalidCryptoSymbol         = errors.New("invalid crypto symbol")
	ErrInvalidRequestStatus        = errors.New("invalid request status")
	ErrInvalidRequestPayload       = errors.New("invalid request payload")
	ErrInvalidCredentials          = errors.New("invalid username or password")
	ErrTokenGenerationFailed       = errors.New("failed to generate token")
	ErrTokenRequired               = errors.New("authorization token required")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInvalidPassword             = errors.New("invalid password")
	ErrInvalidEmail                = errors.New("invalid email")
	ErrInvalidMobile               = errors.New("invalid mobile number")
	ErrSignupFailed                = errors.New("signup failed")
	ErrRetrievingCryptos           = errors.New("error retrieving top cryptocurrencies")
	ErrUnauthorizedAccess          = errors.New("unauthorized access")
	ErrMissingCryptoSymbol         = errors.New("cryptocurrency symbol is required")
	ErrSearchingCrypto             = errors.New("error searching for cryptocurrency")
	ErrSettingPriceAlert           = errors.New("error setting price alert")
	ErrCurrentPriceHigher          = errors.New("current price is higher than target price")
	ErrFailedToCheckNotifications  = errors.New("failed to check notifications")
	ErrFailedToRetrieveUserProfile = errors.New("failed to retrieve user profile")
)

type AppError struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func NewAppError(statusCode int, code int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) ToJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}

// Custom error functions

//Auth Errors

func NewInvalidUsernameError() *AppError {
	return NewAppError(http.StatusBadRequest, 1001, "Invalid username")
}

func NewInvalidPasswordError() *AppError {
	return NewAppError(http.StatusBadRequest, 1002, "Invalid password: must be at least 8 characters, include an uppercase letter, a number, and a special character")
}

func NewInvalidEmailError() *AppError {
	return NewAppError(http.StatusBadRequest, 1003, "Invalid email: must be a valid email address")
}

func NewInvalidMobileError() *AppError {
	return NewAppError(http.StatusBadRequest, 1004, "Invalid mobile number: must be 10 digits")
}

func NewSignupFailedError() *AppError {
	return NewAppError(http.StatusInternalServerError, 1005, "Signup failed")
}

func NewInvalidCredentialsError() *AppError {
	return NewAppError(http.StatusUnauthorized, 1006, "Invalid username or password")
}

//Crypto Response Errors

func NewMissingCryptoSymbolError() *AppError {
	return NewAppError(http.StatusBadRequest, 2001, "Cryptocurrency symbol is required")
}

func NewSearchingCryptoError() *AppError {
	return NewAppError(http.StatusInternalServerError, 2002, "Error searching for cryptocurrency")
}

func NewInvalidCryptoSymbolError() *AppError {
	return NewAppError(http.StatusBadRequest, 2003, "Invalid crypto symbol")
}

func NewRetrievingCryptosError() *AppError {
	return NewAppError(http.StatusInternalServerError, 2004, "Error retrieving top cryptocurrencies")
}

//DB Errors

func NewUserNotFoundError() *AppError {
	return NewAppError(http.StatusNotFound, 3001, "User not found")
}

func NewFailedToFetchProfilesError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3002, "Failed to fetch user profiles")
}

func NewFailedToDeleteUserError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3003, "Failed to delete user")
}

func NewFailedToDelegateUserError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3004, "Failed to delegate user")
}

func NewFailedToFetchRequestsError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3005, "Failed to fetch requests")
}

func NewFailedToUpdateRequestsError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3006, "Failed to update requests")
}

func NewFailedToCheckNotificationsError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3007, "Failed to check notifications")
}

func NewFailedToRetrieveUserProfileError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3008, "Failed to retrieve user profile")
}

func NewSettingPriceAlertError() *AppError {
	return NewAppError(http.StatusInternalServerError, 3009, "Error setting price alert")
}

//Middleware Errors

func NewTokenGenerationFailedError() *AppError {
	return NewAppError(http.StatusInternalServerError, 4001, "Failed to generate token")
}

func NewTokenRequiredError() *AppError {
	return NewAppError(http.StatusUnauthorized, 4002, "Authorization token required")
}

func NewInvalidTokenError() *AppError {
	return NewAppError(http.StatusUnauthorized, 4003, "Invalid token")
}

func NewUnauthorizedAccessError() *AppError {
	return NewAppError(http.StatusUnauthorized, 4004, "Unauthorized access")
}

func NewBlackListedTokenError() *AppError {
	return NewAppError(http.StatusUnauthorized, 4005, "Blacklisted token")
}

//Api Request Errors

func NewInvalidRequestStatusError() *AppError {
	return NewAppError(http.StatusBadRequest, 5001, "Invalid request status")
}

func NewInvalidRequestPayloadError() *AppError {
	return NewAppError(http.StatusBadRequest, 5002, "Invalid request payload")
}

//Other Errors

func NewCurrentPriceHigherError() *AppError {
	return NewAppError(http.StatusOK, 6001, "Current price is higher than the target price")
}
