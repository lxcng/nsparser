package server

import (
	"fmt"
	"log"
	"nsparser/config"
	"nsparser/page"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func StartServer() {
	router := fasthttprouter.New()

	router.GET("/", homeHandler)
	router.GET("/api/load/:title", showHandler)

	server := fasthttp.Server{
		Handler: router.Handler,
	}

	log.Fatal(fasthttp.ListenAndServe(":8080", server.Handler))
}

func homeHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("req: %v\n", ctx.Time())
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Response.SetBody(page.GetHome())
}

func showHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(302)
	ctx.Response.Header.Set("Location", "/")
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	titleRaw := ctx.UserValue("title")
	title, ok := titleRaw.(string)
	if !ok {
		return
	}
	fmt.Println(title)

	err := config.Start(title)
	if err != nil {
		log.Println(err)
		return
	}
}

func runHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

}
