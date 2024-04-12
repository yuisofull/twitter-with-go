package uploadbusiness

import (
	"context"
	"twitter/modules/upload/uploadmodel"
)

type ListImageStorage interface {
	ListImages(
		ctx context.Context,
		ids []int,
		moreKeys ...string,
	) ([]uploadmodel.Upload, error)
}

type listImageBiz struct {
	store ListImageStorage
}

func NewListImageBiz(store ListImageStorage) *listImageBiz {
	return &listImageBiz{store: store}
}

func (biz *listImageBiz) ListImages(ctx context.Context, ids []int, moreKeys ...string) ([]uploadmodel.Upload, error) {
	result, err := biz.store.ListImages(ctx, ids, moreKeys...)

	if err != nil {
		return nil, err
	}

	return result, nil
}
