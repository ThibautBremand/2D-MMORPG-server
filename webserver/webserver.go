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
	"strings"

	"github.com/gorilla/mux"
)

var r *mux.Router

func Launch() {
	clientPath := config.EnvVar("CLIENT_PATH")
	adminPath := config.EnvVar("ADMIN_PATH")

	r = mux.NewRouter()
	r.PathPrefix("/client/").Handler(http.FileServer(http.Dir(clientPath)))

	// For Character and Gamemap administration pages
	r.PathPrefix("/admin/").Handler(http.FileServer(http.Dir(adminPath)))

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

	// For Character and Gamemap administration pages
	r.HandleFunc("/newcharacter", newCharacter).Methods("POST")
	r.HandleFunc("/newgamemap", newGamemap).Methods("POST")
	r.HandleFunc("/character", serveCharacterPage)
	r.HandleFunc("/gamemap", serveGamemapPage)
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

// For Character administration page
func serveCharacterPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/character" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adminPath := config.EnvVar("ADMIN_PATH")
	if adminPath[len(adminPath)-1:] != "/" {
		adminPath = fmt.Sprintf("%s/", adminPath)
	}
	http.ServeFile(w, r, fmt.Sprintf("%sadmin/index.html", adminPath))
}

// For Gamemap administration page
func serveGamemapPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/gamemap" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adminPath := config.EnvVar("ADMIN_PATH")
	if adminPath[len(adminPath)-1:] != "/" {
		adminPath = fmt.Sprintf("%s/", adminPath)
	}
	http.ServeFile(w, r, fmt.Sprintf("%sadmin/loadmap.html", adminPath))
}

// newEntityData contains the name and the JSON properties of a new entity
// that a user has just created using the /newcharacter and /newgamemap endpoints.
// This entity is meant to be persisted into the storage.
type newEntityData struct {
	Name  string
	Props string
}

func newCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	characterData, err := decode(w, r)
	if err != nil {
		return
	}
	fmt.Printf("New character data (name: %s) has been received from a user", characterData.Name)
	err = controller.PersistNewCharacter(strings.ToLower(characterData.Name), characterData.Props)
	respond(w, err)
}

func newGamemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	gamemapData, err := decode(w, r)
	if err != nil {
		return
	}
	fmt.Printf("New gamemap data (name: %s) has been received from a user", gamemapData.Name)
	err = controller.PersistNewGamemap(strings.ToLower(gamemapData.Name), gamemapData.Props)
	respond(w, err)
}

func decode(w http.ResponseWriter, r *http.Request) (newEntityData, error) {
	var entityData newEntityData
	err := json.NewDecoder(r.Body).Decode(&entityData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		w.Write([]byte("KO"))
		return entityData, err
	}

	return entityData, nil
}

func respond(w http.ResponseWriter, err error) {
	var e *utils.NameAlreadyTaken
	if errors.As(err, &e) {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte("KO"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
