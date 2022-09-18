package main

import (
	"log"
	"html/template"
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/ldarren/fity/cfg"
)

type Ably struct {
	Key string
}

func indexTemplate(fname string) *template.Template {
	return template.Must(template.New(fname).ParseFiles("asset/" + fname))
}

func indexHandler(w http.ResponseWriter, req bunrouter.Request) error {
	key := cfg.GetStr([]string{"Ably", "Key"})
	ably := Ably{key}
	return indexTemplate("sse.html").Execute(w, ably)
}

func debugHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"route":  req.Route(),
		"params": req.Params().Map(),
	})
}

func server() {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/", indexHandler)

	router.WithGroup("/api", func(g *bunrouter.Group) {
		g.GET("/users/:id", debugHandler)
		g.GET("/users/current", debugHandler)
		g.GET("/users/*path", debugHandler)
	})

	addr := cfg.GetStr([]string{"Server", "Addr"})
	log.Printf("listening on %v", addr)
	log.Println(http.ListenAndServe(addr, router))
}

func main() {
	cfg.LoadJSON()
	server()
}
