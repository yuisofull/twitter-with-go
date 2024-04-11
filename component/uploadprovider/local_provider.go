package uploadprovider

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"twitter/common"
)

type localProvider struct {
	domain string
	folder string
}

func NewLocalProvider(domain, folder string) *localProvider {
	return &localProvider{
		domain: domain,
		folder: folder,
	}
}

func (*localProvider) String() string {
	return "Local Provider"
}

func (provider *localProvider) SaveFileUploaded(_ context.Context, data []byte, dst string) (*common.Image, error) {
	filePath := fmt.Sprintf("%s/%s", provider.folder, dst)
	// Create the directory if it does not exist
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	err := ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		return nil, err
	}
	return &common.Image{
		Url:       fmt.Sprintf("%s/%s", provider.domain, dst),
		CloudName: "Local",
	}, nil
}
