package uploadbusiness

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"twitter/component/uploadprovider"
	"twitter/modules/upload/uploadmodel"
	// _ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type CreateImageStorage interface {
	CreateImage(context context.Context, data *uploadmodel.Upload) error
}

type uploadBiz struct {
	provider   uploadprovider.UploadProvider
	imageStore CreateImageStorage
}

func NewUploadBiz(provider uploadprovider.UploadProvider, imageStore CreateImageStorage) *uploadBiz {
	return &uploadBiz{provider: provider, imageStore: imageStore}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*uploadmodel.Upload, error) {
	fileBytes := bytes.NewBuffer(data)

	w, h, err := getImageDimension(fileBytes)

	if err != nil {
		return nil, uploadmodel.ErrFileIsNotImage(err)
	}

	if strings.TrimSpace(folder) == "" {
		folder = "images"
	}

	fileExt := filepath.Ext(fileName)
	fileName = fmt.Sprintf("%s-%v%s", fileNameWithoutExtSliceNotation(fileName), time.Now().UnixNano(), fileExt)

	img, err := biz.provider.SaveFileUploaded(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))

	if err != nil {
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	img.Width = w
	img.Height = h
	img.CloudName = biz.provider.String()
	img.Extension = fileExt

	img2 := &uploadmodel.Upload{
		Url:    img.Url,
		Width:  img.Width,
		Height: img.Height,
	}
	if err := biz.imageStore.CreateImage(ctx, img2); err != nil {
		// delete img on S3
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	return img2, nil
}

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

// other way
func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(reader)

	if err != nil {
		log.Println("Err====>", err)
		return 0, 0, err
	}

	return img.Width, img.Height, err
}
