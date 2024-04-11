package ginuser

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	"twitter/component/hasher"
	userbiz "twitter/modules/user/business"
	"twitter/modules/user/usermodel"
	"twitter/modules/user/userstore"
)

func Register(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMyDBConnection()
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		}

		store := userstore.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()
		repo := userbiz.NewRegisterBusiness(store, md5)

		if err := repo.Register(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		data.Mask(false)
		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(data.FakeID.String()))
	}
}
