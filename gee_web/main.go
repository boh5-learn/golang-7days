package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			_, err := fmt.Fprintf(w, "Header[%s] = %s\n", k, v)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	log.Fatal(r.Run(":9999"))
}
