package services

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
)

type AdminService interface {
	ChangeUserStatus(username string) error
	DeleteUser(username string) error
	ViewUserProfiles() ([]*models.User, error)
	ManageUserRequests() ([]*models.UnavailableCryptoRequest, error)
	ManageSpecificCryptoRequests(cryptoSymbol string) ([]*models.UnavailableCryptoRequest, error)
	UpdateRequestStatus(requests []*models.UnavailableCryptoRequest, status string) error
}

type AdminServiceImpl struct {
	repo repositories.AdminRepository
}

func NewAdminService(repo repositories.AdminRepository) AdminService {
	return &AdminServiceImpl{repo: repo}
}

func (s *AdminServiceImpl) ChangeUserStatus(username string) error {
	return s.repo.ChangeUserStatus(username)
}

func (s *AdminServiceImpl) DeleteUser(username string) error {
	return s.repo.DeleteUser(username)
}

func (s *AdminServiceImpl) ViewUserProfiles() ([]*models.User, error) {
	return s.repo.ViewUserProfiles()
}

func (s *AdminServiceImpl) ManageUserRequests() ([]*models.UnavailableCryptoRequest, error) {
	return s.repo.ManageUserRequests()
}

func (s *AdminServiceImpl) ManageSpecificCryptoRequests(cryptoSymbol string) ([]*models.UnavailableCryptoRequest, error) {
	return s.repo.ManageSpecificCryptoRequests(cryptoSymbol)
}

func (s *AdminServiceImpl) UpdateRequestStatus(requests []*models.UnavailableCryptoRequest, status string) error {

	for _, request := range requests {
		request.Status = status
	}

	return s.repo.SaveUnavailableCryptoRequest(requests)
}
