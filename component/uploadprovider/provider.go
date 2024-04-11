package uploadprovider

import (
	"context"
	"fmt"
	"twitter/common"
)

type UploadProvider interface {
	fmt.Stringer
	SaveFileUploaded(context context.Context, data []byte, dst string) (*common.Image, error)
}
