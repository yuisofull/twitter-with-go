package ginuser

import (
	"github.com/gin-gonic/gin"
	"twitter/common"
	"twitter/component/appctx"

	"net/http"
)

// GetProfile
// @Summary Get profile
// @Description Get profile
// @Tags users
// @ID get-profile
// @Accept  json
// @Produce  json
// @Success 200 {object} common.simpleSuccessRes{data=usermodel.User}
// @Security ApiKeyAuth
// @Router /profile [get]
func GetProfile(_ appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		data := c.MustGet(common.CurrentUser).(common.Requester)
		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(data))
	}
}
