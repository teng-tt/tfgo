package main

import (
	"net/http"
	"tfgo/tfb/tfb"
)

func main() {
	r := tfb.New()
	r.GET("/", func(c *tfb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Tfb</h1>")
	})

	r.GET("/hello", func(c *tfb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *tfb.Context) {
		c.JSON(http.StatusOK, tfb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")

}