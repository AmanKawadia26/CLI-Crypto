package main

import (
	"cryptotracker/REST-API/Handlers"
	"cryptotracker/REST-API/middleware"
	"cryptotracker/internal/repositories"
	"cryptotracker/internal/services"
	"cryptotracker/pkg/config"
	"cryptotracker/pkg/globals"
	"cryptotracker/pkg/logger"
	"cryptotracker/pkg/ui"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	conn := globals.GetPgConn()
	defer globals.ClosePgConn()

	config.LoadConfig()

	adminRepo := repositories.NewPostgresAdminRepository(conn)
	adminService := services.NewAdminService(adminRepo)

	userRepo := repositories.NewPostgresUserRepository(conn)
	userService := services.NewUserService(userRepo)

	cryptoRepo := repositories.NewPostgresCryptoRepository(conn)
	cryptoService := services.NewCryptoService(cryptoRepo)

	notificationRepo := repositories.NewPostgresNotificationRepository(conn)
	notificationService := services.NewNotificationService(notificationRepo)

	authRepo := repositories.NewPostgresAuthRepository(conn)
	authService := services.NewAuthService(authRepo, notificationService)

	authHandler := Handlers.NewAuthHandler(authService)
	userHandler := Handlers.NewUserHandler(userService)
	adminHandler := Handlers.NewAdminHandler(adminService)
	cryptoHandler := Handlers.NewCryptoHandler(cryptoService)
	notificationHandler := Handlers.NewNotificationHandler(notificationService)

	// Use the logger from logger.go
	logger.Logger.Info("Starting the server...")

	r := mux.NewRouter()

	// Start the HTTP server
	go func() {

		jwtMiddleware := middleware.JWTTokenService{}

		adminMiddleware := jwtMiddleware.AdminMiddleware
		userMiddleware := jwtMiddleware.UserMiddleware

		r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
		r.HandleFunc("/signup", authHandler.SignupHandler).Methods("POST")
		r.Handle("/logout", userMiddleware(http.HandlerFunc(authHandler.LogoutHandler))).Methods("POST")
		r.Handle("/users/me", userMiddleware(http.HandlerFunc(userHandler.UserProfile))).Methods("GET")
		r.Handle("/notifications", userMiddleware(http.HandlerFunc(notificationHandler.CheckNotificationHandler))).Methods("GET")
		r.Handle("/cryptos", userMiddleware(http.HandlerFunc(cryptoHandler.DisplayTopCryptos))).Methods("GET")
		r.Handle("/cryptos/{cryptoname}", userMiddleware(http.HandlerFunc(cryptoHandler.DisplayCryptoByName))).Methods("GET")
		r.Handle("/cryptos/alert", userMiddleware(http.HandlerFunc(cryptoHandler.SetPriceAlert))).Methods("POST")
		r.Handle("/cryptos", adminMiddleware(http.HandlerFunc(cryptoHandler.DisplayTopCryptos))).Methods("GET")
		r.Handle("/cryptos/{cryptoname}", adminMiddleware(http.HandlerFunc(cryptoHandler.DisplayCryptoByName))).Methods("GET")
		r.Handle("/cryptos/alert", adminMiddleware(http.HandlerFunc(cryptoHandler.SetPriceAlert))).Methods("POST")
		r.Handle("/admin/profiles", adminMiddleware(http.HandlerFunc(adminHandler.Profiles))).Methods("GET")
		//r.Handle("/admin/profiles/{username}", adminMiddleware(http.HandlerFunc(adminHandler.SpecificUserProfile))).Methods("GET")
		r.Handle("/admin/delete/{username}", adminMiddleware(http.HandlerFunc(adminHandler.DeleteUser))).Methods("DELETE")
		r.Handle("/admin/delegate/{username}", adminMiddleware(http.HandlerFunc(adminHandler.DelegateUser))).Methods("PATCH")
		r.Handle("/admin/requests", adminMiddleware(http.HandlerFunc(adminHandler.UnavailableCryptoRequests))).Methods("GET")
		r.Handle("/admin/requests/{username}", adminMiddleware(http.HandlerFunc(adminHandler.SpecificUserUnavailableCryptoRequests))).Methods("GET")
		r.Handle("/admin/requests/{crypto}", adminMiddleware(http.HandlerFunc(adminHandler.ActOnUnavailableCryptoRequestsBySymbol))).Methods("PUT")

		logger.Logger.Info("Server listening on port :5555")

		http.Handle("/", r)
		if err := http.ListenAndServe(":5556", nil); err != nil {
			logger.Logger.Fatal("Failed to start server", err)
		}
	}()

	ui.DisplayWelcomeBanner()

	ui_var := ui.NewUI(userService, adminService, cryptoService, notificationService)

	user, Role := ui_var.AuthenticateUser(conn)

	if Role == "admin" {
		ui.ShowAdminPanel(conn, adminService)
		return
	}

	ui.MainMenu(conn, user, userService, cryptoService)
}
