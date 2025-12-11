package dto

import "be-knowledge/internal/entities"

type DokumenManagement_DownloadDokumen_Response struct {
	Data entities.User `json:"data"`
}

type DokumenManagement_DownloadDokumen_Request struct {
	Id int `json:"id"`
}
