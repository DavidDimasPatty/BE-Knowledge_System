package dto

import "be-knowledge/internal/entities"

type DokumenManagement_GetAllDokumen_Response struct {
	Data []entities.Dokumen `json:"data"`
}

type DokumenManagement_GetAllDokumen_Request struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	RoleId int    `json:"roleId"`
	Email  string `json:"email"`
	NoTelp string `json:"noTelp"`
}
