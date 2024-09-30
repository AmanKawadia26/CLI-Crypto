package middleware_test

import (
	"bytes"
	"cryptotracker/REST-API/middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var jwtService = &middleware.JWTTokenService{}

func TestGenerateToken(t *testing.T) {
	username := "testuser"
	role := "user"

	token, err := jwtService.GenerateToken(username, role)
	assert.NoError(t, err, "Error generating token")

	// Parse the token and verify claims
	claims, err := jwtService.ParseToken(token)
	assert.NoError(t, err, "Error parsing token")
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

func TestParseToken(t *testing.T) {
	username := "testuser"
	role := "user"
	tokenString, _ := jwtService.GenerateToken(username, role)

	// Parse valid token
	claims, err := jwtService.ParseToken(tokenString)
	assert.NoError(t, err, "Error parsing valid token")
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)

	// Parse invalid token
	_, err = jwtService.ParseToken("invalidToken")
	assert.Error(t, err, "Expected error when parsing invalid token")
}

func TestBlacklistToken(t *testing.T) {
	token, _ := jwtService.GenerateToken("testuser", "user")
	jwtService.BlacklistToken(token)

	assert.True(t, jwtService.IsTokenBlacklisted(token), "Token should be blacklisted")
}

func TestIsTokenBlacklisted(t *testing.T) {
	token, _ := jwtService.GenerateToken("testuser1", "user")
	assert.False(t, jwtService.IsTokenBlacklisted(token), "Token should not be blacklisted yet")

	jwtService.BlacklistToken(token)
	assert.True(t, jwtService.IsTokenBlacklisted(token), "Token should be blacklisted")
}

func TestExtractToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer validToken")

	token := jwtService.ExtractToken(req)
	assert.Equal(t, "validToken", token, "Extracted token should match")
}

type MockJWTTokenService struct {
	middleware.JWTTokenService
	MockParseToken func(tokenString string) (*middleware.Claims, error)
}

func (m *MockJWTTokenService) ParseToken(tokenString string) (*middleware.Claims, error) {
	if m.MockParseToken != nil {
		return m.MockParseToken(tokenString)
	}
	return nil, nil
}

// Capture log output
var buf bytes.Buffer

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		expectedCode int
		expectedLog  string
	}{
		{
			name:         "No Token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Blacklisted Token",
			token:        "blacklistedToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Invalid Token",
			token:        "invalidToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Non-Admin Role",
			token:        "userToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Valid Admin Token",
			token:        "adminToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()
			mockService := &MockJWTTokenService{
				MockParseToken: func(tokenString string) (*middleware.Claims, error) {
					switch tokenString {
					case "invalidToken":
						return nil, jwt.ErrSignatureInvalid
					case "adminToken":
						return &middleware.Claims{Username: "adminuser", Role: "admin"}, nil
					case "userToken":
						return &middleware.Claims{Username: "testuser", Role: "user"}, nil
					default:
						return nil, nil
					}
				},
			}

			if tt.token == "blacklistedToken" {
				mockService.BlacklistToken(tt.token)
			}

			adminHandler := mockService.AdminMiddleware(handler)
			adminHandler.ServeHTTP(rr, req)

			// Check logs
			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.expectedLog)
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Reset buffer for next test case
			buf.Reset()
		})
	}
}

func TestUserMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		expectedCode int
		expectedLog  string
	}{
		{
			name:         "No Token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Blacklisted Token",
			token:        "blacklistedToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Invalid Token",
			token:        "invalidToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Non-User Role",
			token:        "adminToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "",
		},
		{
			name:         "Valid User Token",
			token:        "userToken",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "User access granted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()
			mockService := &MockJWTTokenService{
				MockParseToken: func(tokenString string) (*middleware.Claims, error) {
					switch tokenString {
					case "invalidToken":
						return nil, jwt.ErrSignatureInvalid
					case "adminToken":
						return &middleware.Claims{Username: "adminuser", Role: "admin"}, nil
					case "userToken":
						return &middleware.Claims{Username: "testuser", Role: "user"}, nil
					default:
						return nil, nil
					}
				},
			}

			if tt.token == "blacklistedToken" {
				mockService.BlacklistToken(tt.token)
			}

			userHandler := mockService.UserMiddleware(handler)
			userHandler.ServeHTTP(rr, req)

			// Check logs
			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.expectedLog)
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Reset buffer for next test case
			buf.Reset()
		})
	}
}

func TestAdminMiddlewareWithoutToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	adminHandler := jwtService.AdminMiddleware(handler)

	// Case: No token provided
	adminHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "No token should return unauthorized")
}

func TestUserMiddlewareWithoutToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	userHandler := jwtService.UserMiddleware(handler)

	// Case: No token provided
	userHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "No token should return unauthorized")
}
