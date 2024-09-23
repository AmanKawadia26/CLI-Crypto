package services

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
	"github.com/jackc/pgx/v4"
	//"github.com/jackc/pgx/v5"
)

type CryptoService interface {
	DisplayTopCryptocurrencies() ([]interface{}, error)
	SearchCryptocurrency(conn *pgx.Conn, user *models.User, cryptoSymbol string) (float64, string, string, error)
	SetPriceAlert(conn *pgx.Conn, user *models.User, symbol string, targetPrice float64) (float64, error)
}

type CryptoServiceImpl struct {
	cryptoRepo repositories.CryptoRepository
}

func NewCryptoService(cryptoRepo repositories.CryptoRepository) *CryptoServiceImpl {
	return &CryptoServiceImpl{
		cryptoRepo: cryptoRepo,
	}
}

func (s *CryptoServiceImpl) DisplayTopCryptocurrencies() ([]interface{}, error) {
	return s.cryptoRepo.DisplayTopCryptocurrencies()
}

func (s *CryptoServiceImpl) SearchCryptocurrency(conn *pgx.Conn, user *models.User, cryptoSymbol string) (float64, string, string, error) {
	return s.cryptoRepo.SearchCryptocurrency(conn, user, cryptoSymbol)
}

func (s *CryptoServiceImpl) SetPriceAlert(conn *pgx.Conn, user *models.User, symbol string, targetPrice float64) (float64, error) {
	return s.cryptoRepo.SetPriceAlert(conn, user, symbol, targetPrice)
}
