package uploadmodel

import (
	"errors"
	"twitter/common"
)

const EntityName = "Upload"
const TableName = "images"

type Upload struct {
	common.SQLModel `json:",inline" gorm:"embedded" swaggerignore:"true"`
	Url             string `json:"url" gorm:"column:url"`
	Width           int    `json:"width" gorm:"column:width"`
	Height          int    `json:"height" gorm:"column:height"`
}

func (Upload) TableName() string {
	return TableName
}

func (u *Upload) Mask(isAdmin bool) {
	u.GenUID(common.DBTypeUpload)
}

func (u *Upload) ToImage() *common.Image {
	return &common.Image{
		Id:     u.Id,
		Url:    u.Url,
		Width:  u.Width,
		Height: u.Height,
	}

}

var (
	ErrFileTooLarge = common.NewCustomError(
		errors.New("file too large"),
		"file too large",
		"ErrFileTooLarge",
	)
)

func ErrCannotSaveFile(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"cannot save uploaded file",
		"ErrCannotSaveFile",
	)
}

func ErrFileIsNotImage(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"file is not image",
		"ErrFileIsNotImage",
	)
}
