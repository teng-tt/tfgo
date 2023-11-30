package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"tfgo/tfb/tfb"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatDate(t time.Time) string {
	year, month, daya := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, daya)
}

func onlyForV2() tfb.HandlerFunc {
	return func(c *tfb.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := tfb.Default()
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{
		Name: "David",
		Age:  22,
	}
	stu2 := &student{
		Name: "Jack",
		Age:  30,
	}

	r.GET("/", func(c *tfb.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/student", func(c *tfb.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", tfb.H{
			"title":  "tfb",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *tfb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", tfb.H{
			"title": "tfb",
			"now":   time.Date(2022, 8, 12, 0, 0, 0, 0, time.UTC),
		})
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *tfb.Context) {
		names := []string{"tfbxxxas"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")

}
