package main

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/internal/services"
	"cryptotracker/pkg/config"
	"cryptotracker/pkg/globals"
	"cryptotracker/pkg/ui"
)

func main() {

	// Initialize PostgreSQL Client
	conn := globals.GetPgConn()
	defer globals.ClosePgConn()

	// Load the configuration
	config.LoadConfig()

	// Display welcome banner
	ui.DisplayWelcomeBanner()

	// Start login/signup process
	user, Role := ui.AuthenticateUser(conn)

	adminRepo := repositories.NewPostgresAdminRepository(conn)
	adminService := services.NewAdminService(adminRepo)

	// If user is admin, show admin panel
	if Role == "admin" {
		ui.ShowAdminPanel(conn, adminService)
		return
	}

	userRepo := repositories.NewPostgresUserRepository(conn)
	userService := services.NewUserService(userRepo)

	cryptoRepo := repositories.NewPostgresCryptoRepository(conn)
	cryptoService := services.NewCryptoService(cryptoRepo)

	// Display main user menu
	ui.MainMenu(conn, user, userService, cryptoService)
}
