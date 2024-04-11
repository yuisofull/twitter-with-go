package appctx

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"twitter/component/uploadprovider"
)

type AppContext interface {
	GetMyDBConnection() *gorm.DB
	GetSecretKey() string
	GetMyESConnection() *elastic.Client
	UploadProvider() uploadprovider.UploadProvider
}

type appCtx struct {
	db             *gorm.DB
	secretKey      string
	uploadProvider uploadprovider.UploadProvider
	es             *elastic.Client
}

func NewAppContext(db *gorm.DB, secretKey string, uploadprovider uploadprovider.UploadProvider, es *elastic.Client) *appCtx {
	return &appCtx{
		db:             db,
		secretKey:      secretKey,
		uploadProvider: uploadprovider,
		es:             es,
	}
}

func (ctx *appCtx) GetMyDBConnection() *gorm.DB {
	return ctx.db
}
func (ctx *appCtx) GetSecretKey() string {
	return ctx.secretKey
}
func (ctx *appCtx) UploadProvider() uploadprovider.UploadProvider {
	return ctx.uploadProvider
}
func (ctx *appCtx) GetMyESConnection() *elastic.Client {
	return ctx.es
}
