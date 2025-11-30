package dto

type DokumenManagement_AddDokumen_Request struct {
	Judul    string `json:"judul"`
	AddId    string `json:"addId"`
	FileName string `json:"fileName"`
	FileData []byte `json:"fileData"`
}
