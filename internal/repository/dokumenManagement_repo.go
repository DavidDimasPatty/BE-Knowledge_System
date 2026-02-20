package repository

import (
	dto "be-knowledge/internal/delivery/dto/dokumenManagement"
	"be-knowledge/internal/entities"
	Tracelog "be-knowledge/internal/tracelog"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
)

type DokumenManagementRepository interface {
	GetAllDokumen() (data *dto.DokumenManagement_GetAllDokumen_Response, er error)
	AddDokumen(data *dto.DokumenManagement_AddDokumen_Request, filePath string) error
	EditDokumenGet(id int) (data *entities.Dokumen, er error)
	EditDokumen(data *dto.DokumenManagement_EditDokumen_Request, filePath *string) error
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
	namaEndpoint := "GetAllDokumen"

	res := dto.DokumenManagement_GetAllDokumen_Response{}

	dokumens := []entities.Dokumen{}

	query := "SELECT * FROM dokumen"
	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s", query),
		namaEndpoint)

	err := r.db.Select(&dokumens, query)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return nil, err
	}

	res.Data = dokumens

	return &res, nil
}

func (r *dokumenManagementRepository) AddDokumen(data *dto.DokumenManagement_AddDokumen_Request, filePath string) error {
	namaEndpoint := "AddDokumen"
	query := `
        INSERT INTO dokumen (link, judul, addId, ADDTIME)
        VALUES (?, ?, ?, NOW())
    `
	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s | Params: link=%v, judul=%v, addId=%v", query, filePath, data.Judul, data.AddId),
		namaEndpoint)
	_, err := r.db.Exec(query, filePath, data.Judul, data.AddId)

	payload := map[string]string{
		"file_path": path.Base(filePath),
	}

	jsonData, _ := json.Marshal(payload)
	baseURL := os.Getenv("URL_PYTHON")
	Tracelog.DokumenManagementLog(fmt.Sprintf("Melakukan request ke : %v dengan data : %v", baseURL, payload), namaEndpoint)
	req, err := http.NewRequest(
		"POST",
		baseURL+"/insertDocument",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", resp.Body),
			namaEndpoint,
		)
		return fmt.Errorf("python api failed: %s", resp.Status)
	}

	return err
}

func (r *dokumenManagementRepository) EditDokumenGet(id int) (*entities.Dokumen, error) {
	namaEndpoint := "EditDokumenGet"
	dokumen := entities.Dokumen{}
	query := "SELECT * FROM dokumen WHERE id = ?"

	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint)
	err := r.db.Get(&dokumen, query, id)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return nil, err
	}

	return &dokumen, nil
}

func (r *dokumenManagementRepository) EditDokumen(data *dto.DokumenManagement_EditDokumen_Request, filePath *string) error {
	var (
		query string
		err   error
	)
	namaEndpoint := "EditDokumen"
	if filePath != nil {
		//delete
		dokumenOld := entities.Dokumen{}
		selectQuery := "SELECT * FROM dokumen WHERE id = ?"
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("SQL: %s | Params: id=%v,", selectQuery, data.Id),
			namaEndpoint)

		if err := r.db.Get(&dokumenOld, selectQuery, data.Id); err != nil {
			Tracelog.DokumenManagementLog(
				fmt.Sprintf("Error : %v", err.Error()),
				namaEndpoint,
			)
			return err
		}

		payload := map[string]string{
			"file_path_old": path.Base(dokumenOld.Link),
			"file_path_new": path.Base(*filePath),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		baseURL := os.Getenv("URL_PYTHON")
		Tracelog.DokumenManagementLog(fmt.Sprintf("Melakukan request ke : %v dengan data : %v", baseURL, payload), namaEndpoint)
		req, err := http.NewRequest(
			"POST",
			baseURL+"/editDocument",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			Tracelog.DokumenManagementLog(
				fmt.Sprintf("Error : %v", err.Error()),
				namaEndpoint,
			)
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			Tracelog.DokumenManagementLog(
				fmt.Sprintf("Error : %v", err.Error()),
				namaEndpoint,
			)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			Tracelog.DokumenManagementLog(
				fmt.Sprintf("Error : %v", resp.Body),
				namaEndpoint,
			)
			return fmt.Errorf("python api failed: %s", resp.Status)
		}

		query = `
		UPDATE dokumen
		SET  link = ?, judul = ?, updId = ?, updTime = NOW()
		WHERE id = ?
	`
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("SQL: %s | Params: link=%v, judul=%v, updId=%v, id=%v", query, filePath, data.Judul, data.UpdId, data.Id),
			namaEndpoint)
		_, err = r.db.Exec(query,
			*filePath,
			data.Judul,
			data.UpdId,
			data.Id,
		)
	} else {
		query = `
		UPDATE dokumen
		SET  judul = ?, updId = ?, updTime = NOW()
		WHERE id = ?
	`
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("SQL: %s | Params: judul=%v, updId=%v, id=%v", query, data.Judul, data.UpdId, data.Id),
			namaEndpoint)
		_, err = r.db.Exec(query,
			data.Judul,
			data.UpdId,
			data.Id,
		)
	}

	return err
}

func (r *dokumenManagementRepository) DeleteDokumen(id int) error {
	dokumen := entities.Dokumen{}
	namaEndpoint := "DeleteDokumen"
	selectQuery := "SELECT * FROM dokumen WHERE id = ?"
	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", selectQuery, id),
		namaEndpoint)

	if err := r.db.Get(&dokumen, selectQuery, id); err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}

	deleteQuery := "DELETE FROM dokumen WHERE id = ?"
	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", deleteQuery, id),
		namaEndpoint)
	if _, err := r.db.Exec(deleteQuery, id); err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}

	payload := map[string]string{
		"filename": path.Base(dokumen.Link),
	}

	jsonData, _ := json.Marshal(payload)

	baseURL := os.Getenv("URL_PYTHON")
	Tracelog.DokumenManagementLog(fmt.Sprintf("Melakukan request ke : %v dengan data : %v", baseURL, payload), namaEndpoint)
	req, err := http.NewRequest(
		"POST",
		baseURL+"/deleteDocument",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", resp.Body),
			namaEndpoint,
		)
		return fmt.Errorf("python api failed: %s", resp.Status)
	}

	return err
}

func (r *dokumenManagementRepository) DownloadDokumen(id int) (*entities.Dokumen, []byte, error) {
	namaEndpoint := "DownloadDokumen"
	dok := entities.Dokumen{}
	query := "SELECT * FROM dokumen WHERE id = ?"
	Tracelog.DokumenManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint)
	err := r.db.Get(&dok, query, id)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return nil, nil, err
	}

	if dok.Link == "" {
		Tracelog.DokumenManagementLog(
			"Error : file link is empty",
			namaEndpoint,
		)
		return nil, nil, errors.New("file link is empty")
	}

	fileBytes, err := os.ReadFile(dok.Link)
	if err != nil {
		Tracelog.DokumenManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return nil, nil, err
	}

	return &dok, fileBytes, nil
}
