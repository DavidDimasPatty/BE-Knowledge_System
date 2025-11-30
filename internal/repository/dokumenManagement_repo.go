package repository

import (
	dto "be-knowledge/internal/delivery/dto/dokumenManagement"
	"be-knowledge/internal/entities"
	"errors"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type DokumenManagementRepository interface {
	GetAllDokumen() (data *dto.DokumenManagement_GetAllDokumen_Response, er error)
	AddDokumen(data *dto.DokumenManagement_AddDokumen_Request) error
	EditDokumenGet(id int) (data *entities.Dokumen, er error)
	EditDokumen(data *dto.DokumenManagement_EditDokumen_Request) error
	DeleteDokumen(id int) error
	DownloadDokumen(id int) (*entities.Dokumen, []byte, error)
}

type dokumenManagementRepository struct {
	db *sqlx.DB
}

func NewDokumenManagementRepository(db *sqlx.DB) DokumenManagementRepository {
	return &dokumenManagementRepository{db}
}

func (r *dokumenManagementRepository) GetAllDokumen() (*dto.DokumenManagement_GetAllDokumen_Response, error) {
	res := dto.DokumenManagement_GetAllDokumen_Response{}

	dokumens := []entities.Dokumen{}

	query := "SELECT * FROM dokumen"

	err := r.db.Select(&dokumens, query)
	if err != nil {
		return nil, err
	}

	res.Data = dokumens

	return &res, nil
}

func (r *dokumenManagementRepository) AddDokumen(data *dto.DokumenManagement_AddDokumen_Request) error {

	parentPath := os.Getenv("STORAGE_PATH")

	if parentPath == "" {
		return errors.New("STORAGE_PATH is not set in .env")
	}

	link := filepath.Join(parentPath, data.FileName)

	query := `
		INSERT INTO dokumen
		(link, judul, addId, ADDTIME)
		VALUES (?, ?, ? ,NOW())
	`

	_, err := r.db.Exec(query,
		link,
		data.Judul,
		data.AddId,
	)

	return err
}

func (r *dokumenManagementRepository) EditDokumenGet(id int) (*entities.Dokumen, error) {
	dokumen := entities.Dokumen{}
	query := "SELECT * FROM dokumen WHERE id = ?"

	err := r.db.Get(&dokumen, query, id)
	if err != nil {
		return nil, err
	}

	return &dokumen, nil
}

func (r *dokumenManagementRepository) EditDokumen(data *dto.DokumenManagement_EditDokumen_Request) error {
	query := `
		UPDATE dokumen
		SET  email = ?, noTelp = ?, nama = ?, roles = ?, updId = ?, updTime = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		data.Email,
		data.NoTelp,
		data.Nama,
		data.RoleId,
		data.UpdId,
		data.Id,
	)

	return err
}

func (r *dokumenManagementRepository) DeleteDokumen(id int) error {
	query := "DELETE FROM dokumen WHERE id = ?"

	_, err := r.db.Exec(query, id)
	return err
}

func (r *dokumenManagementRepository) DownloadDokumen(id int) (*entities.Dokumen, []byte, error) {
	dok := entities.Dokumen{}
	query := "SELECT * FROM dokumen WHERE id = ?"

	err := r.db.Get(&dok, query, id)
	if err != nil {
		return nil, nil, err
	}

	if dok.Link == "" {
		return nil, nil, errors.New("file link is empty")
	}

	fileBytes, err := os.ReadFile(dok.Link)
	if err != nil {
		return nil, nil, err
	}

	return &dok, fileBytes, nil
}
