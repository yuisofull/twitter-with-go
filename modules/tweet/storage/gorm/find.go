package tweetstorage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

func (s *sqlStore) FindTweetWithCondition(
	context context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*tweetmodel.Tweet, error) {
	var data tweetmodel.Tweet
	if err := s.db.Where(condition).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
