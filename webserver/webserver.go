package webserver

import (
	"fmt"
	"log"
	"net/http"
	"server/communication"
	"server/config"

	"github.com/gorilla/mux"
)

var r *mux.Router

func Launch() {
	clientPath := config.EnvVar("CLIENT_PATH")

	r = mux.NewRouter()
	r.PathPrefix("/client/").Handler(http.FileServer(http.Dir(clientPath)))

	startHub()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}

// startHub launches the Hub which stores all the WebSocket clients.
func startHub() {
	hub := communication.NewHub()
	go hub.Run()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		communication.ServeWs(hub, w, r, params["name"][0])
	})
}

// serveHome serves the html frontpage.
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientPath := config.EnvVar("CLIENT_PATH")
	if clientPath[len(clientPath)-1:] != "/" {
		clientPath = fmt.Sprintf("%s/", clientPath)
	}
	http.ServeFile(w, r, fmt.Sprintf("%sclient/home.html", clientPath))
}
