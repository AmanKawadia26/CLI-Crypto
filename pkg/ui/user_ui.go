////go:build !test
//// +build !test
//
//package ui
//
//import (
//	//"cryptotracker/internal/repositories"
//	"cryptotracker/internal/services"
//	"cryptotracker/models"
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/olekukonko/tablewriter"
//	"go.mongodb.org/mongo-driver/mongo"
//	"log"
//	"os"
//	"strings"
//)
//
//// MainMenu displays the main menu for a regular user
//func MainMenu(client *mongo.Client, user *models.User, userService *services.UserServiceImpl, cryptoService *services.CryptoServiceImpl) {
//	for {
//		ClearScreen()
//		DisplayMainMenu()
//
//		var choice int
//		color.New(color.FgYellow).Print("Enter your choice: ")
//		fmt.Scan(&choice)
//
//		switch choice {
//		case 1:
//			DisplayTopCryptocurrencies(cryptoService)
//		case 2:
//			SearchCryptocurrency(client, user, cryptoService)
//		case 3:
//			SetPriceAlert(client, user, cryptoService)
//		case 4:
//			DisplayUserProfile(userService, user.Username)
//		case 5:
//			color.New(color.FgCyan).Println("Logging out...")
//			log.Println("Logging out...")
//			os.Exit(0)
//		default:
//			color.New(color.FgRed).Println("Invalid choice, please try again.")
//		}
//	}
//}
//
//// DisplayUserProfile handles the display of the user profile
//func DisplayUserProfile(userService *services.UserServiceImpl, username string) {
//	user, err := userService.GetUserProfile(username)
//	if err != nil {
//		color.New(color.FgRed).Println("Error fetching user profile:", err)
//		return
//	}
//
//	// Display user profile
//	fmt.Println()
//
//	// Print section title with a border
//	color.New(color.FgGreen, color.Bold).Println("==== User Profile ====")
//	fmt.Println()
//
//	// Define maximum width for formatting
//	width := 20
//
//	// Print user profile details
//	printDetail("Username", user.Username, width)
//	printDetail("Email", user.Email, width)
//	printDetail("Mobile", fmt.Sprintf("%d", user.Mobile), width)
//	printDetail("Role", user.Role, width)
//
//	fmt.Println()
//}
//
//// printDetail prints a formatted profile detail
//func printDetail(label, value string, width int) {
//	labelColor := color.New(color.FgCyan, color.Bold)
//	valueColor := color.New(color.FgWhite)
//
//	labelColor.Printf("%-*s: ", width, label)
//	valueColor.Println(value)
//}
//
//func DisplayTopCryptocurrencies(cryptoService *services.CryptoServiceImpl) {
//
//	data, err := cryptoService.DisplayTopCryptocurrencies()
//	if err != nil {
//		color.New(color.FgRed).Printf("Error displaying cryptocurrencies: %s\n", err)
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//
//	// Set table headers
//	table.SetHeader([]string{"No.", "Name", "Symbol", "Price"})
//
//	// Set table column alignment and padding
//	table.SetAlignment(tablewriter.ALIGN_LEFT)
//	table.SetBorder(true) // Enable borders
//	table.SetHeaderLine(true)
//	table.SetCenterSeparator("|")
//	table.SetColumnSeparator("|")
//	table.SetRowLine(true)
//	table.SetHeaderColor(
//		tablewriter.Colors{tablewriter.Bold}, // Header color and style
//		tablewriter.Colors{tablewriter.Bold},
//		tablewriter.Colors{tablewriter.Bold},
//		tablewriter.Colors{tablewriter.Bold},
//	)
//
//	fmt.Println()
//	color.New(color.FgGreen).Println("Top 10 Cryptocurrencies:")
//	fmt.Println()
//	for i, crypto := range data {
//		cryptoMap := crypto.(map[string]interface{})
//		name := cryptoMap["name"].(string)
//		symbol := cryptoMap["symbol"].(string)
//		price := cryptoMap["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)
//
//		// Add row to table
//		table.Append([]string{
//			fmt.Sprintf("%d", i+1),
//			name,
//			symbol,
//			fmt.Sprintf("$%.2f", price),
//		})
//	}
//	fmt.Println()
//
//	table.Render()
//}
//
//// SearchCryptocurrency prompts user for a cryptocurrency symbol or name and searches for it
//func SearchCryptocurrency(client *mongo.Client, user *models.User, cryptoService *services.CryptoServiceImpl) {
//	var input string
//	color.New(color.FgCyan).Print("Enter the symbol or name of the cryptocurrency: ")
//	fmt.Scan(&input)
//
//	// Normalize the input to lowercase for case-insensitive comparison
//	input = strings.ToLower(input)
//
//	// Use the crypto service to search for the cryptocurrency
//	price, cryptoName, cryptoSymbol, err := cryptoService.SearchCryptocurrency(client, user, input)
//	if err.Error() == fmt.Sprintf("Request to add the cryptocurrency has been submitted.") {
//		color.New(color.FgGreen).Print("Request to add the cryptocurrency has been submitted.")
//	} else if price != 0 {
//		fmt.Println()
//		color.New(color.FgGreen).Printf("%s (%s): $%.2f\n", cryptoName, cryptoSymbol, price)
//		fmt.Println()
//		fmt.Println()
//
//		DisplayCryptoGraph(cryptoName, price)
//	} else {
//
//		color.New(color.FgRed).Print("Error searching cryptocurreny : ", err)
//		//fmt.Println()
//		//
//		//DisplayCryptoGraph(cryptoName, price)
//	}
//}
//
//// SetPriceAlert prompts user to set a price alert
//func SetPriceAlert(client *mongo.Client, user *models.User, cryptoService *services.CryptoServiceImpl) {
//	var symbol string
//	var targetPrice float64
//
//	color.New(color.FgCyan).Print("Enter the cryptocurrency symbol: ")
//	fmt.Scan(&symbol)
//	color.New(color.FgCyan).Print("Enter the target price: ")
//	fmt.Scan(&targetPrice)
//
//	currentPrice, err := cryptoService.SetPriceAlert(client, user, symbol, targetPrice)
//	if err.Error() == fmt.Sprintf("%s is still below your target price. Current price: $%.2f. Notification created.\n", symbol, currentPrice) {
//		color.New(color.FgGreen).Printf("%s is still below your target price. Current price: $%.2f. Notification created.\n", symbol, currentPrice)
//	} else if err.Error() == fmt.Sprintf("Alert: %s has reached your target price of $%.2f. Current price: $%.2f\n", symbol, targetPrice, currentPrice) {
//		color.New(color.FgGreen).Printf("Alert: %s has reached your target price of $%.2f. \nCurrent price: $%.2f\n", symbol, targetPrice, currentPrice)
//	} else {
//		color.New(color.FgRed).Printf("Error setting price alert: %v\n", err)
//	}
//}

//go:build !test
// +build !test

package ui

import (
	"cryptotracker/internal/services"
	"cryptotracker/models"
	"fmt"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strings"
	//"github.com/jackc/pgx/v5"
)

// MainMenu displays the main menu for a regular user
func MainMenu(conn *pgx.Conn, user *models.User, userService *services.UserServiceImpl, cryptoService *services.CryptoServiceImpl) {
	for {
		ClearScreen()
		DisplayMainMenu()

		var choice int
		color.New(color.FgYellow).Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			DisplayTopCryptocurrencies(cryptoService)
		case 2:
			SearchCryptocurrency(conn, user, cryptoService)
		case 3:
			SetPriceAlert(conn, user, cryptoService)
		case 4:
			DisplayUserProfile(userService, user.Username)
		case 5:
			color.New(color.FgCyan).Println("Logging out...")
			log.Println("Logging out...")
			os.Exit(0)
		default:
			color.New(color.FgRed).Println("Invalid choice, please try again.")
		}
	}
}

// DisplayUserProfile handles the display of the user profile
func DisplayUserProfile(userService *services.UserServiceImpl, username string) {
	user, err := userService.GetUserProfile(username)
	if err != nil {
		color.New(color.FgRed).Println("Error fetching user profile:", err)
		return
	}

	// Display user profile
	fmt.Println()

	// Print section title with a border
	color.New(color.FgGreen, color.Bold).Println("==== User Profile ====")
	fmt.Println()

	// Define maximum width for formatting
	width := 20

	// Print user profile details
	printDetail("Username", user.Username, width)
	printDetail("Email", user.Email, width)
	printDetail("Mobile", fmt.Sprintf("%d", user.Mobile), width)
	printDetail("Role", user.Role, width)

	fmt.Println()
}

// printDetail prints a formatted profile detail
func printDetail(label, value string, width int) {
	labelColor := color.New(color.FgCyan, color.Bold)
	valueColor := color.New(color.FgWhite)

	labelColor.Printf("%-*s: ", width, label)
	valueColor.Println(value)
}

func DisplayTopCryptocurrencies(cryptoService *services.CryptoServiceImpl) {

	data, err := cryptoService.DisplayTopCryptocurrencies()
	if err != nil {
		color.New(color.FgRed).Printf("Error displaying cryptocurrencies: %s\n", err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	// Set table headers
	table.SetHeader([]string{"No.", "Name", "Symbol", "Price"})

	// Set table column alignment and padding
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true) // Enable borders
	table.SetHeaderLine(true)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowLine(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold}, // Header color and style
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	fmt.Println()
	color.New(color.FgGreen).Println("Top 10 Cryptocurrencies:")
	fmt.Println()
	for i, crypto := range data {
		cryptoMap := crypto.(map[string]interface{})
		name := cryptoMap["name"].(string)
		symbol := cryptoMap["symbol"].(string)
		price := cryptoMap["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"].(float64)

		// Add row to table
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			name,
			symbol,
			fmt.Sprintf("$%.2f", price),
		})
	}
	fmt.Println()

	table.Render()
}

func SearchCryptocurrency(conn *pgx.Conn, user *models.User, cryptoService *services.CryptoServiceImpl) {

	if conn == nil {
		log.Fatalf("Error: PostgreSQL connection is nil")
	}

	if user == nil {
		log.Fatalf("Error: User object is nil")
	}

	var input string
	color.New(color.FgCyan).Print("Enter the symbol or name of the cryptocurrency: ")
	fmt.Scan(&input)

	// Normalize the input to lowercase for case-insensitive comparison
	input = strings.ToLower(input)

	// Use the crypto service to search for the cryptocurrency
	price, cryptoName, cryptoSymbol, err := cryptoService.SearchCryptocurrency(conn, user, input)
	if err != nil && err.Error() == "Request to add the cryptocurrency has been submitted." {
		color.New(color.FgGreen).Print("Request to add the cryptocurrency has been submitted.")
	} else if price != 0 {
		fmt.Println()
		color.New(color.FgGreen).Printf("%s (%s): $%.2f\n", cryptoName, cryptoSymbol, price)
		fmt.Println()
		fmt.Println()

		DisplayCryptoGraph(cryptoName, price)
	} else {
		color.New(color.FgRed).Print("Error searching cryptocurrency: ", err)
	}
}

// SetPriceAlert prompts user to set a price alert
func SetPriceAlert(conn *pgx.Conn, user *models.User, cryptoService *services.CryptoServiceImpl) {
	var symbol string
	var targetPrice float64

	color.New(color.FgCyan).Print("Enter the cryptocurrency symbol: ")
	fmt.Scan(&symbol)
	color.New(color.FgCyan).Print("Enter the target price: ")
	fmt.Scan(&targetPrice)

	currentPrice, err := cryptoService.SetPriceAlert(conn, user, symbol, targetPrice)
	if err.Error() == fmt.Sprintf("%s is still below your target price. Current price: $%.2f. Notification created.\n", symbol, currentPrice) {
		color.New(color.FgGreen).Printf("%s is still below your target price. Current price: $%.2f. Notification created.\n", symbol, currentPrice)
	} else if err.Error() == fmt.Sprintf("Alert: %s has reached your target price of $%.2f. Current price: $%.2f\n", symbol, targetPrice, currentPrice) {
		color.New(color.FgGreen).Printf("Alert: %s has reached your target price of $%.2f. \nCurrent price: $%.2f\n", symbol, targetPrice, currentPrice)
	} else {
		color.New(color.FgRed).Printf("Error setting price alert: %v\n", err)
	}
}
