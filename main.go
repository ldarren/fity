package main

import (
	"fmt"
	"log"
	"net/http"
)

func postMsgByChannel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fit and able")
}

func server() {
	http.HandleFunc("/channels/channelId/messages", postMsgByChannel)
	log.Fatalln(http.ListenAndServe(":4751", nil))
}

func main() {
	server()
}
