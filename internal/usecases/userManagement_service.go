package usecases

import (
	dto "be-knowledge/internal/delivery/dto/userManagement"
	"errors"

	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
)

type UserManagementService interface {
	GetAllUser() (data *dto.UserManagement_GetAllUser_Response, er error)
	AddUser(data dto.UserManagement_AddUser_Request) error
	EditUserGet(id int) (data *entities.User, er error)
	EditUser(data dto.UserManagement_EditUser_Request) error
	DeleteUser(id int) error
}

type userManagementService struct {
	repo repository.UserManagementRepository
}

func NewUserManagementService(repo repository.UserManagementRepository) UserManagementService {
	return &userManagementService{repo}
}

func (s *userManagementService) GetAllUser() (*dto.UserManagement_GetAllUser_Response, error) {
	user, err := s.repo.GetAllUser()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (s *userManagementService) AddUser(data dto.UserManagement_AddUser_Request) error {
	return s.repo.AddUser(data)
}

func (s *userManagementService) EditUserGet(id int) (*entities.User, error) {
	user, err := s.repo.EditUserGet(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (s *userManagementService) EditUser(data dto.UserManagement_EditUser_Request) error {
	return s.repo.EditUser(data)
}

func (s *userManagementService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
