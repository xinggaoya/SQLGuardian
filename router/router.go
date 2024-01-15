package router

import (
	"SQLGuardian/app"
	"SQLGuardian/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
  @author: XingGao
  @date: 2023/8/13
**/

func InitRouters(router *gin.Engine) {
	view := router.Group("/")
	{
		view.GET("/favicon.ico", func(c *gin.Context) {
			// 返回静态文件
			c.File(consts.StaticDir + "favicon.ico")
		})
		view.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

	}

	api := router.Group("/api")
	{
		api.POST("/download", app.DownloadFile)
		api.POST("/delete", app.DeleteFile)
		api.GET("/config", app.SetConfig)
		api.GET("/file/all", app.GetBackupFilesList)
		api.POST("/config", app.SetConfig)
	}
}
