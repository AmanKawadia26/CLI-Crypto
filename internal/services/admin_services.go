package services

import (
	"cryptotracker/internal/repositories"
	"cryptotracker/models"
	"github.com/jackc/pgx/v4"
)

type AdminService interface {
	ChangeUserStatus(conn *pgx.Conn, username string) error
	DeleteUser(conn *pgx.Conn, username string) error
	ViewUserProfiles(conn *pgx.Conn) ([]*models.User, error)
	ManageUserRequests(conn *pgx.Conn) ([]*models.UnavailableCryptoRequest, error)
	UpdateRequestStatus(conn *pgx.Conn, request *models.UnavailableCryptoRequest, status string) error
}

type AdminServiceImpl struct {
	repo repositories.AdminRepository
}

func NewAdminService(repo repositories.AdminRepository) AdminService {
	return &AdminServiceImpl{repo: repo}
}

func (s *AdminServiceImpl) ChangeUserStatus(conn *pgx.Conn, username string) error {
	return s.repo.ChangeUserStatus(conn, username)
}

func (s *AdminServiceImpl) DeleteUser(conn *pgx.Conn, username string) error {
	return s.repo.DeleteUser(conn, username)
}

func (s *AdminServiceImpl) ViewUserProfiles(conn *pgx.Conn) ([]*models.User, error) {
	return s.repo.ViewUserProfiles(conn)
}

func (s *AdminServiceImpl) ManageUserRequests(conn *pgx.Conn) ([]*models.UnavailableCryptoRequest, error) {
	return s.repo.ManageUserRequests(conn)
}

func (s *AdminServiceImpl) UpdateRequestStatus(conn *pgx.Conn, request *models.UnavailableCryptoRequest, status string) error {
	request.Status = status
	return s.repo.SaveUnavailableCryptoRequest(conn, request)
}
