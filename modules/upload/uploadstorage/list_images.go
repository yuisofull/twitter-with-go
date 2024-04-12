package uploadstorage

import (
	"context"
	"twitter/common"
	"twitter/modules/upload/uploadmodel"
)

func (store *sqlStore) ListImages(
	ctx context.Context,
	ids []int,
	moreKeys ...string,
) ([]uploadmodel.Upload, error) {
	db := store.db

	var result []uploadmodel.Upload

	db = db.Table(uploadmodel.TableName)

	if err := db.Where("id in (?)", ids).
		Find(&result).
		Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
