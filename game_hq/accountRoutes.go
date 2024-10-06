package hq

import (
	"accounts"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	Username, Password string
}
type RegisterRequestBody struct {
	Username, Password string
}
type LogoutRequestBody struct {
	UserHash string
}
type ResponseBody struct {
	Success     bool   `json:"Success"`
	AccountHash string `json:"AccountHash"`
}

func accountRoutes(r *gin.Engine, agh *accounts.AccountsGamesHandler) {

	r.GET("register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.go.tmpl", gin.H{"Title": "Register"})
	})
	r.GET("login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.go.tmpl", gin.H{"Title": "Login"})
	})

	r.POST("account/register/", func(c *gin.Context) {
		registerBody := &RegisterRequestBody{}
		c.Bind(registerBody)
		err := agh.Register(registerBody.Username, registerBody.Password)
		if err != nil {
			c.AbortWithStatus(401)
		} else {
			c.JSON(http.StatusOK, ResponseBody{Success: true})
		}

	})
	r.POST("account/login/", func(c *gin.Context) {
		loginBody := &LoginRequestBody{}
		c.Bind(loginBody)
		fmt.Println(loginBody)

		returnHash, err := agh.Login(loginBody.Username, loginBody.Password)
		if err != nil {
			//handle
			c.AbortWithStatus(401)
		} else {
			c.JSON(http.StatusOK, ResponseBody{Success: true, AccountHash: returnHash})
		}
	})
	r.POST("account/logout/:hash", func(c *gin.Context) {

	})
}
