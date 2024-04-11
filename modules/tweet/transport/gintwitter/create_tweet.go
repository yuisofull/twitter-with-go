package gintwitter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	tweetbusiness "twitter/modules/tweet/business"
	tweetmodel "twitter/modules/tweet/model"
	tweetstorage "twitter/modules/tweet/storage/gorm"
)

func CreateTweet(ctx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := ctx.GetMyDBConnection()

		var data tweetmodel.TweetCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}
		data.UserID = requester.GetUserId()
		store := tweetstorage.NewSQLStore(db)
		biz := tweetbusiness.NewCreateTweetBusiness(store)

		if err := biz.CreateTweet(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		data.Mask(false)

		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(data.FakeID.String()))
	}
}
