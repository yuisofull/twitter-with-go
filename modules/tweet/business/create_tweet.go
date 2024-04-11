package tweetbusiness

import (
	"context"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

type CreateTweetStore interface {
	Create(context.Context, *tweetmodel.TweetCreate) error
}

type createTweetBusiness struct {
	store CreateTweetStore
}

func NewCreateTweetBusiness(store CreateTweetStore) *createTweetBusiness {
	return &createTweetBusiness{store: store}
}

func (business *createTweetBusiness) CreateTweet(context context.Context, data *tweetmodel.TweetCreate) error {
	if err := data.Validate(); err != nil {
		return common.ErrInvalidRequest(err)
	}

	if err := business.store.Create(context, data); err != nil {
		return common.ErrCannotCreateEntity(tweetmodel.EntityName, err)
	}
	return nil
}
