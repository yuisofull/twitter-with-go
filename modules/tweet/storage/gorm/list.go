package tweetstorage

import (
	"context"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

func (s *sqlStore) ListTweetWithCondition(
	_ context.Context,
	filter *tweetmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]tweetmodel.Tweet, error) {
	var result []tweetmodel.Tweet

	db := s.db.Table(tweetmodel.TableName)

	if f := filter; f != nil {
		if f.UserID > 0 {
			db = db.Where("user_id = ?", f.UserID)
		}
		if len(f.Status) > 0 {
			db = db.Where("status in (?)", f.Status)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)
		if err != nil {
			return nil, common.ErrDB(err)
		}
		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).
		Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if len(result) > 0 {
		last := result[len(result)-1]
		last.Mask(false)
		paging.NextCursor = last.FakeID.String()
	}

	return result, nil
}
