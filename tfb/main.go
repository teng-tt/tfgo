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
		// expect /hello?name=tf
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *tfb.Context) {
		// expect /hello/tfb
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *tfb.Context) {
		c.JSON(http.StatusOK, tfb.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")

}
