package uploadstorage

import (
	"context"
	"twitter/common"
	"twitter/modules/upload/uploadmodel"
)

func (store *sqlStore) CreateImage(_ context.Context, data *uploadmodel.Upload) error {
	db := store.db.Omit("Status", "UpdatedAt")

	if err := db.Table(uploadmodel.TableName).
		Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
