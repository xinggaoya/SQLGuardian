package router

import (
	"SQLGuardian/app/service"
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
		view.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"itme": []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			})
		})

	}

	api := router.Group("/api")
	{
		api.POST("/download", service.DownloadFile)
		api.POST("/delete", service.DeleteFile)
		api.GET("/config", service.SetConfig)
		api.GET("/file/all", service.GetBackupFilesList)
		api.POST("/config", service.SetConfig)
	}
}
