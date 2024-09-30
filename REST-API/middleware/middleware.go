package middleware

import (
	"context"
	errs "cryptotracker/REST-API/errors"
	"cryptotracker/pkg/logger" // Import the logger package
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var secretKey = []byte("abcdefg")

var tokenBlacklist = make(map[string]bool)

type TokenService interface {
	ParseToken(tokenString string) (*Claims, error)
	GenerateToken(username, role string) (string, error)
	BlacklistToken(token string)
	IsTokenBlacklisted(token string) bool
	ExtractToken(r *http.Request) string
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type JWTTokenService struct{}

func (j *JWTTokenService) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := j.ExtractToken(r)
		if tokenString == "" {
			logger.Logger.Warn("Authorization token required", "method", r.Method, "path", r.URL.Path)
			//http.Error(w, "Authorization token required", http.StatusUnauthorized)
			errs.NewTokenRequiredError().ToJSON(w)
			return
		}

		// Check if the token is blacklisted
		if j.IsTokenBlacklisted(tokenString) {
			logger.Logger.Warn("Token is blacklisted", "method", r.Method, "path", r.URL.Path)
			//http.Error(w, "Token is invalid", http.StatusUnauthorized)
			errs.NewBlackListedTokenError().ToJSON(w)
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("Invalid token", "error", err)
			//http.Error(w, "Invalid token", http.StatusUnauthorized)
			errs.NewInvalidTokenError().ToJSON(w)
			return
		}

		// Check if the user has the "admin" role
		if claims.Role != "admin" {
			logger.Logger.Warn("Access denied. Admin role required", "username", claims.Username, "role", claims.Role)
			//http.Error(w, "Access denied. Admin role required", http.StatusForbidden)
			errs.NewUnauthorizedAccessError().ToJSON(w)
			return
		}

		logger.Logger.Info("Admin access granted", "username", claims.Username, "role", claims.Role)
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(r.Context(), "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWTTokenService) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := j.ExtractToken(r)
		if tokenString == "" {
			logger.Logger.Warn("Authorization token required", "method", r.Method, "path", r.URL.Path)
			//http.Error(w, "Authorization token required", http.StatusUnauthorized)
			errs.NewTokenRequiredError().ToJSON(w)
			return
		}

		// Check if the token is blacklisted
		if j.IsTokenBlacklisted(tokenString) {
			logger.Logger.Warn("Token is blacklisted", "method", r.Method, "path", r.URL.Path)
			//http.Error(w, "Token is invalid", http.StatusUnauthorized)
			errs.NewBlackListedTokenError().ToJSON(w)
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("Invalid token", "error", err)
			//http.Error(w, "Invalid token", http.StatusUnauthorized)
			errs.NewInvalidTokenError().ToJSON(w)
			return
		}

		// Check if the user has the "user" role
		if claims.Role != "user" {
			logger.Logger.Warn("Access denied. User role required", "username", claims.Username, "role", claims.Role)
			//http.Error(w, "Access denied. User role required", http.StatusForbidden)
			errs.NewUnauthorizedAccessError().ToJSON(w)
			return
		}

		logger.Logger.Info("User access granted", "username", claims.Username, "role", claims.Role)
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ParseToken implements the TokenService interface
func (j *JWTTokenService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

// GenerateToken implements the TokenService interface
func (j *JWTTokenService) GenerateToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logger.Logger.Error("Error generating token", err)
		return "", err
	}

	return tokenString, nil
}

// BlacklistToken implements the TokenService interface
func (j *JWTTokenService) BlacklistToken(token string) {
	tokenBlacklist[token] = true
}

// IsTokenBlacklisted implements the TokenService interface
func (j *JWTTokenService) IsTokenBlacklisted(token string) bool {
	_, exists := tokenBlacklist[token]
	return exists
}

// ExtractToken implements the TokenService interface
func (j *JWTTokenService) ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
