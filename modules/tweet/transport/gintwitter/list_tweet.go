package gintwitter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	tweetbusiness "twitter/modules/tweet/business"
	tweetmodel "twitter/modules/tweet/model"
	tweetstorage2 "twitter/modules/tweet/storage/elastic"
	tweetstorage "twitter/modules/tweet/storage/gorm"
)

func ListTweet(ctx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := ctx.GetMyDBConnection()
		es := ctx.GetMyESConnection()
		var pagingData common.Paging
		if err := c.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		pagingData.Fulfill()

		var filter tweetmodel.Filter
		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if filter.FakeUserID != "" {
			uid, err := common.FromBase58(filter.FakeUserID)
			if err != nil {
				panic(common.ErrInvalidRequest(err))
			}
			filter.UserID = int(uid.GetLocalID())
		}
		filter.Status = []int{1}

		var store interface{}

		if filter.Search == "" {
			store = tweetstorage.NewSQLStore(db)
		} else {
			store = tweetstorage2.NewESStore(es)
		}

		//store := tweetstorage2.NewESStore(es)
		biz := tweetbusiness.NewListTweetBusiness(store.(tweetbusiness.ListTwitterStore))

		result, err := biz.ListTweet(c.Request.Context(), &filter, &pagingData)
		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))
	}
}
