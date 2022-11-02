package models

import "github.com/go-playground/validator/v10"

type Content struct {
	ContentId      string `json:"contentId" validate:"required"`
	BatchId        string `json:"batchId" validate:"required"`
	Status         int    `json:"status" validate:"required"`
	CourseId       string `json:"courseId" validate:"required"`
	LastAccessTime string `json:"lastAccessTime" validate:"required"`
}

type request struct {
	UserId   string    `json:"userId" validate:"required"`
	Contents []Content `json:"contents" validate:"required"`
}

type Request struct {
	Request request `json:"request" validate:"required"`
}

func Validate(request Request) error {
	v := validator.New()
	err := v.Struct(request)
	return err

}
