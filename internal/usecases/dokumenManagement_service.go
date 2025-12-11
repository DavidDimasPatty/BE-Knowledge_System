package usecases

import (
	dto "be-knowledge/internal/delivery/dto/dokumenManagement"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
)

type DokumenManagementService interface {
	GetAllDokumen() (data *dto.DokumenManagement_GetAllDokumen_Response, er error)
	AddDokumen(data *dto.DokumenManagement_AddDokumen_Request) error
	EditDokumenGet(id int) (data *entities.Dokumen, er error)
	EditDokumen(data *dto.DokumenManagement_EditDokumen_Request) error
	DeleteDokumen(id int) error
	DownloadDokumen(id int) (*entities.Dokumen, []byte, error)
}

type dokumenManagementService struct {
	repo repository.DokumenManagementRepository
}

func NewDokumenManagementService(repo repository.DokumenManagementRepository) DokumenManagementService {
	return &dokumenManagementService{repo}
}

func (s *dokumenManagementService) GetAllDokumen() (data *dto.DokumenManagement_GetAllDokumen_Response, er error) {
	dokumen, err := s.repo.GetAllDokumen()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return dokumen, nil
}

func (s *dokumenManagementService) AddDokumen(data *dto.DokumenManagement_AddDokumen_Request) error {
	parentPath := os.Getenv("STORAGE_PATH")
	if parentPath == "" {
		return errors.New("STORAGE_PATH is not set")
	}

	filePath := filepath.Join(parentPath, data.FileName)
	if err := ioutil.WriteFile(filePath, data.FileData, 0644); err != nil {
		return err
	}

	return s.repo.AddDokumen(data, filePath)
}

func (s *dokumenManagementService) EditDokumenGet(id int) (*entities.Dokumen, error) {
	dokumen, err := s.repo.EditDokumenGet(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return dokumen, nil
}

func (s *dokumenManagementService) EditDokumen(data *dto.DokumenManagement_EditDokumen_Request) error {
	return s.repo.EditDokumen(data)
}

func (s *dokumenManagementService) DeleteDokumen(id int) error {
	return s.repo.DeleteDokumen(id)
}

func (s *dokumenManagementService) DownloadDokumen(id int) (*entities.Dokumen, []byte, error) {
	return s.repo.DownloadDokumen(id)
}
