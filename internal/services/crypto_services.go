package services

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
	//"github.com/jackc/pgx/v5"
)

type CryptoService interface {
	DisplayTopCryptocurrencies(count int) ([]interface{}, error)
	SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, *models.Cryptocurrency, error)
	SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error)
}

type CryptoServiceImpl struct {
	cryptoRepo repositories.CryptoRepository
}

func NewCryptoService(cryptoRepo repositories.CryptoRepository) *CryptoServiceImpl {
	return &CryptoServiceImpl{
		cryptoRepo: cryptoRepo,
	}
}

func (s *CryptoServiceImpl) DisplayTopCryptocurrencies(count int) ([]interface{}, error) {
	return s.cryptoRepo.DisplayTopCryptocurrencies(count)
}

func (s *CryptoServiceImpl) SearchCryptocurrency(user *models.User, cryptoSymbol string) (float64, string, string, *models.Cryptocurrency, error) {
	return s.cryptoRepo.SearchCryptocurrency(user, cryptoSymbol)
}

func (s *CryptoServiceImpl) SetPriceAlert(user *models.User, symbol string, targetPrice float64) (float64, error) {
	return s.cryptoRepo.SetPriceAlert(user, symbol, targetPrice)
}
