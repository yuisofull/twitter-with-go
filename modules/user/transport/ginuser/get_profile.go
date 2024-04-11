package ginuser

import (
	"github.com/gin-gonic/gin"
	"twitter/common"
	"twitter/component/appctx"

	"net/http"
)

func GetProfile(_ appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		data := c.MustGet(common.CurrentUser).(common.Requester)
		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(data))
	}
}
