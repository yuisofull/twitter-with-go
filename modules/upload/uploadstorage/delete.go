package uploadstorage

import (
	"context"
	"twitter/modules/upload/uploadmodel"
)

func (store *sqlStore) DeleteImages(ctx context.Context, ids []int) error {
	db := store.db

	if err := db.Table(uploadmodel.TableName).
		Where("id in (?)", ids).
		Delete(nil).
		Error; err != nil {
		return err
	}

	return nil
}
