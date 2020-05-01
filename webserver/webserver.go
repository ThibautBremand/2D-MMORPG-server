package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/communication"
	"server/config"
	"server/controller"
	"server/utils"

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
	r.HandleFunc("/newcharacter", newCharacter).Methods("POST")
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

type NewCharacterData struct {
	Name  string
	Tiles string
}

type responseServer struct {
	Message string
}

func newCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var characterData NewCharacterData
	err := json.NewDecoder(r.Body).Decode(&characterData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		w.Write([]byte("KO"))
		return
	}
	fmt.Printf("New character data (name: %s) has been received from a user", characterData.Name)
	err = controller.PersistNewCharacter(characterData.Name, characterData.Tiles)
	var e *utils.UsernameTaken
	if errors.As(err, &e) {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte("KO"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
