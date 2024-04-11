package tweetstorage

import (
	"context"
	"gorm.io/gorm"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

func (s *sqlStore) IncreaseLikeCount(
	_ context.Context,
	id int,
) error {
	db := s.db

	if err := db.Table(tweetmodel.TableName).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + ?", 1)).
		Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) DecreaseLikeCount(
	_ context.Context,
	id int,
) error {
	db := s.db

	if err := db.Table(tweetmodel.TableName).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count - ?", 1)).
		Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
