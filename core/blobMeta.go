package core

type BlobMeta struct {
	OriginalFilename string `json:"originalFilename"`
	ContentType      string `json:"contentType"`
	Hash             string `json:"hash"`
}
