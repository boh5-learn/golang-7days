package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=boh5
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=boh5
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

		v1.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(middlewareForV2())
	{
		v2.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=boh5
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	log.Fatal(r.Run(":9999"))
}

func middlewareForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(http.StatusInternalServerError, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
