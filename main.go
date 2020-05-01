package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"server/communication"
	"server/config"
	"server/db"

	"github.com/joho/godotenv"
)

var addr = flag.String("addr", ":8080", "http service address")

// prepareConfig stets up mandatory settings.
func prepareConfig() {
	mime.AddExtensionType(".js", "application/javascript; charset=utf-8")

	// Serve the deployed client at "/client/" path
	clientPath := config.EnvVar("CLIENT_PATH")
	http.Handle("/client/", http.FileServer(http.Dir(clientPath)))
}

// startDatabases connects to the storage and the redis.
func startDatabases() {
	if err := db.Open(); err != nil {
		fmt.Printf("error %v", err)
	}
	db.Start()
}

// startHub launches the Hub which stores all the WebSocket clients.
func startHub() {
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

func startWebServer() {
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// init is invoked before main(), and loads values from .env as env variables.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	prepareConfig()
	startDatabases()
	defer db.Close()

	startHub()
	startWebServer()
}
