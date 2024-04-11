package tweetbusiness

import (
	"context"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

type ListTwitterStore interface {
	ListTweetWithCondition(
		ctx context.Context,
		filter *tweetmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]tweetmodel.Tweet, error)
}

type listTweetsBusiness struct {
	store ListTwitterStore
}

func NewListTweetBusiness(store ListTwitterStore) *listTweetsBusiness {
	return &listTweetsBusiness{store: store}
}

func (biz *listTweetsBusiness) ListTweet(
	ctx context.Context,
	filter *tweetmodel.Filter,
	paging *common.Paging,
	moreKeys ...string) ([]tweetmodel.Tweet, error) {

	//ctx1, span := trace.StartSpan(ctx, "List Restaurant Business")

	result, err := biz.store.ListTweetWithCondition(ctx, filter, paging, "User")

	//span.End()

	if err != nil {
		return nil, common.ErrCannotListEntity(tweetmodel.EntityName, err)
	}

	return result, nil
}
