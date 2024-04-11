package tweetstorage

import (
	"context"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

func (s *sqlStore) Create(_ context.Context, data *tweetmodel.TweetCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
