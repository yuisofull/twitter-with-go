package ginuser

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	"twitter/component/hasher"
	"twitter/component/tokenprovider/jwt"
	userbiz "twitter/modules/user/business"
	"twitter/modules/user/usermodel"
	"twitter/modules/user/userstore"
)

// Login
// @Summary Login
// @Description Login
// @Tags users
// @ID login
// @Accept  json
// @Produce  json
// @Param cinema body usermodel.UserLogin true "User"
// @Success 200 {object} common.simpleSuccessRes{data=tokenprovider.Token}
// @Router /login [post]
func Login(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMyDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.GetSecretKey()) //appctx.SecretKey()

		store := userstore.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		biz := userbiz.NewLoginBusiness(appCtx, store, 60*60*24*30, tokenProvider, md5)
		account, err := biz.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(account))
	}
}
