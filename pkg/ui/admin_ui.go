//go:build !test
// +build !test

package ui

import (
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	"strings"
)

func ShowAdminPanel(conn *pgx.Conn, adminService services.AdminService) {
	for {
		fmt.Println()
		fmt.Println(colorBlue("=================================="))
		fmt.Println(colorYellow("            Admin Menu            "))
		fmt.Println(colorBlue("=================================="))
		fmt.Println()
		fmt.Println("1. Manage Users")
		fmt.Println("2. View User Profiles")
		fmt.Println("3. Manage User Requests")
		fmt.Println("4. Logout")

		var choice int
		color.New(color.FgYellow).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			ManageUsers(conn, adminService)
		case 2:
			ViewUserProfiles(conn, adminService)
		case 3:
			ManageUserRequests(conn, adminService)
		case 4:
			color.New(color.FgCyan).Println("Logging out...")
			return
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}

func ManageUsers(conn *pgx.Conn, adminService services.AdminService) {
	fmt.Println()
	color.New(color.FgGreen).Println("Managing users")
	fmt.Println("1. Change a user status to admin")
	fmt.Println("2. Delete a user")

	var choice int
	color.New(color.FgYellow).Print("Enter your choice: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		var username string
		color.New(color.FgYellow).Print("Enter the username to change role: ")
		fmt.Scan(&username)
		err := adminService.ChangeUserStatus(username)
		if err != nil {
			color.Red("Error changing user status: %v", err)
		}
	case 2:
		var username string
		color.New(color.FgYellow).Print("Enter the username to delete: ")
		fmt.Scan(&username)
		err := adminService.DeleteUser(username)
		if err != nil {
			color.Red("Error deleting user: %v", err)
		}
	default:
		color.Red("Invalid choice")
	}
	fmt.Println()
}

// ViewUserProfiles displays a list of user profiles in a tabular format
func ViewUserProfiles(conn *pgx.Conn, adminService services.AdminService) {
	fmt.Println()
	users, err := adminService.ViewUserProfiles()
	if err != nil {
		color.New(color.FgRed).Println("Error fetching user profiles:", err)
		return
	}

	if len(users) == 0 {
		color.New(color.FgCyan).Println("No users found.")
		return
	}

	// Print table header
	printTableHeader()

	// Print user profiles in table format with correct serial numbers
	displayedNumber := 1
	for _, user := range users {
		if user.Role == "user" {
			printUserProfile(displayedNumber, user)
			displayedNumber++
		}
	}

	fmt.Println()
}

// printTableHeader prints the table header with bold formatting
func printTableHeader() {
	color.New(color.FgGreen, color.Bold).Printf("%-5s %-20s %-30s %-15s %-10s\n", "S.No", "Username", "Email", "Mobile", "Role")
	color.New(color.FgWhite).Println(strings.Repeat("-", 80))
}

// printUserProfile prints a formatted user profile
func printUserProfile(index int, user *models.User) {
	color.New(color.FgYellow).Printf("%-5d %-20s %-30s %-15d %-10s\n", index, user.Username, user.Email, user.Mobile, user.Role)
}

func ManageUserRequests(conn *pgx.Conn, adminService services.AdminService) {
	unavailableRequests, err := adminService.ManageUserRequests()
	if err != nil {
		color.New(color.FgRed).Println("Error fetching unavailable crypto requests:", err)
		return
	}

	if len(unavailableRequests) == 0 {
		color.New(color.FgCyan).Println("No pending unavailable crypto requests.")
		return
	}

	// Print table header with border
	printRequestTableHeader()

	// Print each request in table format
	for i, request := range unavailableRequests {
		printRequestTableRow(i+1, request)
	}

	// Get the crypto symbol from the admin
	var cryptoSymbol string
	color.New(color.FgYellow).Print("Enter the crypto symbol to manage: ")
	fmt.Scan(&cryptoSymbol)

	// Filter requests matching the given crypto symbol
	var matchingRequests []*models.UnavailableCryptoRequest
	for _, request := range unavailableRequests {
		if request.CryptoSymbol == cryptoSymbol {
			matchingRequests = append(matchingRequests, request)
		}
	}

	if len(matchingRequests) == 0 {
		color.New(color.FgRed).Println("No requests found for the given crypto symbol:", cryptoSymbol)
		return
	}

	// Display matching requests
	color.New(color.FgGreen).Printf("Found %d request(s) for the crypto symbol '%s'.\n", len(matchingRequests), cryptoSymbol)
	for i, req := range matchingRequests {
		color.New(color.FgCyan).Printf("%d. User: %s, Request: %s, Status: %s, Timestamp: %s\n", i+1, req.UserName, req.RequestMessage, req.Status, req.Timestamp)
	}

	// Prompt for action
	var action string
	color.New(color.FgYellow).Print("Enter 'approve' to approve the request(s) or 'reject' to reject them: ")
	fmt.Scan(&action)

	// Update all matching requests at once
	var newStatus string
	if action == "approve" {
		newStatus = "Approved"
	} else if action == "reject" {
		newStatus = "Rejected"
	} else {
		color.New(color.FgRed).Println("Invalid action.")
		return
	}

	// Call the update function with the slice of matching requests
	err = adminService.UpdateRequestStatus(matchingRequests, newStatus)
	if err != nil {
		color.New(color.FgRed).Println("Error updating request status:", err)
		return
	}

	color.New(color.FgGreen).Printf("Status of all requests for '%s' updated to '%s'.\n", cryptoSymbol, newStatus)
}

// printRequestTableHeader prints the table header with borders
func printRequestTableHeader() {
	color.New(color.FgGreen, color.Bold).Printf("+-----+-----------------+----------------------+------------------------------+------------+\n")
	color.New(color.FgGreen, color.Bold).Printf("| %3s | %-15s | %-20s | %-28s | %-10s |\n", "S.No", "Symbol", "User", "Message", "Status")
	color.New(color.FgGreen, color.Bold).Printf("+-----+-----------------+----------------------+------------------------------+------------+\n")
}

// printRequestTableRow prints a formatted row of the request table
func printRequestTableRow(index int, request *models.UnavailableCryptoRequest) {
	color.New(color.FgYellow).Printf("| %3d | %-15s | %-20s | %-28s | %-10s |\n", index, request.CryptoSymbol, request.UserName, request.RequestMessage, request.Status)
	color.New(color.FgWhite).Println("+-----+-----------------+----------------------+------------------------------+------------+")
}
