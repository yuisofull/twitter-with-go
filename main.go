package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"twitter/component/appctx"
	"twitter/component/uploadprovider"
	"twitter/middleware"
)

func main() {
	//dsn := os.Getenv("MYSQL_CONN_STRING")
	dsn := "system:admin123@tcp(127.0.0.1:3306)/temp_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db = db.Debug()
	uploadProvider := uploadprovider.NewLocalProvider("localhost:8080/static", "static")
	es, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetBasicAuth("elastic", "changeme"))
	if err != nil {
		log.Fatalln(err)
	}

	appCtx := appctx.NewAppContext(db, "secretKey", uploadProvider, es)

	r := gin.Default()

	r.Use(middleware.Recover(appCtx))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	setupRoute(appCtx, r.Group("/v1"))
	err = r.Run()
	if err != nil {
		log.Println(err)
	}
}
