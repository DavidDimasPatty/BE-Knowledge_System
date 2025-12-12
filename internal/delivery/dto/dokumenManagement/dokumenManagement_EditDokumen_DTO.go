package dto

type DokumenManagement_EditDokumen_Request struct {
	Id       int    `json:"id"`
	Judul    string `json:"judul"`
	UpdId    string `json:"updId"`
	FileName string `json:"fileName"`
	FileData []byte `json:"fileData"`
}
