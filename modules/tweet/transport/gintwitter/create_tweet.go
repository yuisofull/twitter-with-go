package gintwitter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	tweetbusiness "twitter/modules/tweet/business"
	tweetmodel "twitter/modules/tweet/model"
	tweetstorage "twitter/modules/tweet/storage/gorm"
	"twitter/modules/upload/uploadbusiness"
	"twitter/modules/upload/uploadstorage"
)

func CreateTweet(ctx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := ctx.GetMyDBConnection()

		var data tweetmodel.TweetCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		imageIds := make([]int, 0)
		for _, uid := range data.ImageUIDs {
			UID, err := common.FromBase58(uid)
			if err != nil {
				panic(err)
			}
			imageIds = append(imageIds, int(UID.GetLocalID()))
		}

		uploadStore := uploadstorage.NewSQLStore(db)
		uploadBiz := uploadbusiness.NewListImageBiz(uploadStore)
		uploads, err := uploadBiz.ListImages(c.Request.Context(), imageIds)
		var images common.Images
		if err != nil {
			panic(err)
		}
		for _, upload := range uploads {
			images = append(images, *upload.ToImage())
		}
		data.Images = &images

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
