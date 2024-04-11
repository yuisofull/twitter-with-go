package ginupload

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"twitter/common"
	"twitter/component/appctx"
	"twitter/modules/upload/uploadbusiness"
	"twitter/modules/upload/uploadstorage"
)

// Upload file to S3
// 1. Get image/file from request header
// 2. Check file is real image
// 3. Save image
// 1. Save to local machine
// 2. Save to cloud storage (S3)
// 3. Improve security

func Upload(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := appCtx.GetMyDBConnection()
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				panic(common.ErrInvalidRequest(err))
			}
		}(file)

		folder := c.DefaultPostForm("folder", "images")

		// create a slice have length equal to lenth of file size
		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		imageStore := uploadstorage.NewSQLStore(db)
		biz := uploadbusiness.NewUploadBiz(appCtx.UploadProvider(), imageStore)
		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}
		img.Mask(false)
		c.JSON(http.StatusOK, common.SimpleNewSuccessResponse(img.FakeID.String()))

		//c.SaveUploadedFile(fileHeader, fmt.Sprintf("./static/%s", fileHeader.Filename))
		//c.JSON(200, common.SimpleSucessResponse(true))
	}
}
