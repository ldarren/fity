package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/ldarren/fity/cfg"
	"github.com/ldarren/fity/tmpl"
	"github.com/ldarren/fity/pubsub"
)

var ps = pubsub.New(0)

func debugHandler(w http.ResponseWriter, req bunrouter.Request) error {
	log.Println("Receiving event")

	return bunrouter.JSON(w, bunrouter.H{
		"route":  req.Route(),
		"params": req.Params().Map(),
	})
}

type Message struct {
	Text string    `json:"text"`
}

func server(addr string) {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		ably := cfg.Get([]string{"Ably"})
		return tmpl.RenderHTML(w, "sse", ably)
	})

	router.WithGroup("/api", func(g *bunrouter.Group) {
		g.GET("/users/:id/*path", debugHandler)

		g.POST("/topics/:topic", func(w http.ResponseWriter, req bunrouter.Request) error {
			topic := req.Params().ByName("topic")
			var body Message
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return err
			}
			ps.Pub(body.Text, topic)

			return bunrouter.JSON(w, bunrouter.H{
				"code":  req.Route(),
				"params": req.Params().Map(),
			})
		})

		g.GET("/topics/:topic", func(w http.ResponseWriter, req bunrouter.Request) error {

			// Make sure that the writer supports flushing.
			flusher, ok := w.(http.Flusher)

			if !ok {
				http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
				return nil
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			topic := req.Params().ByName("topic")
			ch := ps.Sub(topic)

			// Remove this client from the map of connected clients
			// when this handler exits.
			defer func() {
				go ps.Unsub(ch, topic)
			}()

			// Listen to connection close and un-register messageChan
			// notify := w.(http.CloseNotifier).CloseNotify()
			notify := req.Context().Done()

			go func() {
				<-notify
				go ps.Unsub(ch, topic)
			}()

			for {

				// Write to the ResponseWriter
				// Server Sent Events compatible
				fmt.Fprintf(w, "data: %s\n\n", <-ch)

				// Flush the data immediatly instead of buffering it for later.
				flusher.Flush()
			}
		})
	})

	log.Printf("listening on %v", addr)
	log.Println(http.ListenAndServe(addr, router))
}

func main() {
	cfg.LoadJSON()
	tmpl.Load(cfg.GetStr([]string{"Path", "Asset"}))
	server(cfg.GetStr([]string{"Server", "Addr"}))
}
