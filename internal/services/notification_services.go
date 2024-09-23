package services

import (
	"cryptotracker/internal/api"
	"cryptotracker/internal/repositories"
	//"cryptotracker/models"
	"encoding/json"
	"fmt"
	//"strings"
)

type Notification struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

type NotificationService interface {
	CheckNotification(username string) ([]Notification, error)
	CheckUnavailableCryptoRequestsService(username string) ([]Notification, error)
	CheckPriceAlertService(username string) ([]Notification, error)
}

type NotificationServiceImpl struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &NotificationServiceImpl{repo: repo}
}

func (s *NotificationServiceImpl) CheckNotification(username string) ([]Notification, error) {
	var notifications []Notification

	unavailableCryptoMsgs, err := s.CheckUnavailableCryptoRequestsService(username)
	if err != nil {
		return nil, err
	}
	notifications = append(notifications, unavailableCryptoMsgs...)

	priceAlertMsgs, err := s.CheckPriceAlertService(username)
	if err != nil {
		return nil, err
	}
	notifications = append(notifications, priceAlertMsgs...)

	return notifications, nil
}

func (s *NotificationServiceImpl) CheckUnavailableCryptoRequestsService(username string) ([]Notification, error) {
	rows, err := s.repo.CheckUnavailableCryptoRequestsRepo(username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	index := 1
	for rows.Next() {
		var cryptoSymbol string
		var status string
		err := rows.Scan(&cryptoSymbol, &status)
		if err != nil {
			return nil, err
		}

		message := fmt.Sprintf("Your request for %s has been %s.", cryptoSymbol, status)
		notifications = append(notifications, Notification{Index: index, Message: message})
		index++
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return notifications, nil
}

func (s *NotificationServiceImpl) CheckPriceAlertService(username string) ([]Notification, error) {
	notifications, err := s.repo.CheckPriceAlertsRepo(username)
	if err != nil {
		return nil, err
	}

	var fulfilledNotifications []Notification
	index := 1
	for _, notification := range notifications {
		params := map[string]string{
			"id":      fmt.Sprintf("%d", notification.CryptoID),
			"convert": "USD",
		}
		apiResponse := api.CoinMarketCapClient{}

		response := apiResponse.GetAPIResponse("/listings/latest", params)

		var result map[string]interface{}
		err := json.Unmarshal(response, &result)
		if err != nil {
			return nil, err
		}

		data, dataOk := result["data"].(map[string]interface{})
		if !dataOk || data[fmt.Sprintf("%d", notification.CryptoID)] == nil {
			return nil, fmt.Errorf("crypto data not found for ID: %d", notification.CryptoID)
		}

		cryptoData := data[fmt.Sprintf("%d", notification.CryptoID)].(map[string]interface{})
		priceData := cryptoData["quote"].(map[string]interface{})
		currentPrice := priceData["USD"].(map[string]interface{})["price"].(float64)

		if currentPrice >= notification.TargetPrice {
			err := s.repo.UpdatePriceNotificationStatusRepo(&notification, username, currentPrice)
			if err != nil {
				return nil, err
			}

			message := fmt.Sprintf("Your price alert for %s (ID: %d) at target price %d has been fulfilled at current price $%.2f.", notification.Crypto, notification.TargetPrice, notification.CryptoID, currentPrice)
			fulfilledNotifications = append(fulfilledNotifications, Notification{Index: index, Message: message})
			index++
		}
	}

	return fulfilledNotifications, nil
}
