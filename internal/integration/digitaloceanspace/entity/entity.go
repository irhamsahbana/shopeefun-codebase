package entity

import "mime/multipart"

type XxxRequest struct {
}

type XxxResponse struct {
}

type UploadFileRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}

type UploadFileResponse struct {
	FileName string `json:"filename" validate:"required"`
	Url      string `json:"url" validate:"required"`
}

type DeleteFileRequest struct {
	FileName string `json:"filename" validate:"required"`
}
