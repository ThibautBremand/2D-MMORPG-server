package webserver

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"server/communication"
	"server/config"
)

var addr = flag.String("addr", ":8080", "http service address")

// prepareConfig stets up mandatory settings.
func PrepareConfig() {
	mime.AddExtensionType(".js", "application/javascript; charset=utf-8")

	// Serve the deployed client at "/client/" path
	clientPath := config.EnvVar("CLIENT_PATH")
	http.Handle("/client/", http.FileServer(http.Dir(clientPath)))
}

// startHub launches the Hub which stores all the WebSocket clients.
func StartHub() {
	hub := communication.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		communication.ServeWs(hub, w, r, params["name"][0])
	})
}

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

func StartWebServer() {
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
