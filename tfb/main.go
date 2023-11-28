package main

import (
	"log"
	"net/http"
	"tfgo/tfb/tfb"
	"time"
)

func onlyForV2() tfb.HandlerFunc {
	return func(c *tfb.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := tfb.New()
	r.Use(tfb.Logger()) // global middleware
	r.GET("/", func(c *tfb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Tfb</h1>")
	})

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *tfb.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *tfb.Context) {
			c.JSON(http.StatusOK, tfb.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")

}
