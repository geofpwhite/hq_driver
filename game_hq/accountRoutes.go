package hq

import (
	"accounts"
	"net/http"

	"github.com/gin-gonic/gin"
)

func accountRoutes(r *gin.Engine, agh *accounts.AccountsGamesHandler) {

	r.GET("account/register/:username/:password", func(c *gin.Context) {
		usr, flag := c.Params.Get("username")
		if !flag {
			//error
			c.AbortWithStatus(401)
		}
		passwd, flag := c.Params.Get("password")
		if !flag {
			//error
			c.AbortWithStatus(401)
		}
		agh.Register(usr, passwd)
	})
	r.GET("account/login/:username/:password", func(c *gin.Context) {
		usr, flag := c.Params.Get("username")
		if !flag {
			//error
			c.AbortWithStatus(401)
		}
		passwd, flag := c.Params.Get("password")
		if !flag {
			//error
			c.AbortWithStatus(401)
		}
		returnHash, err := agh.Login(usr, passwd)
		if err != nil {
			//handle
			c.AbortWithStatus(401)
		}
		c.String(http.StatusOK, "%s", returnHash)
	})
	r.GET("account/logout/:hash", func(c *gin.Context) {

	})
}
