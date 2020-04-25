package communication

import (
	"encoding/json"
	"fmt"
	"server/controller"
	"server/entity"
	"server/game"
)

const (
	// Keys sent to the clients
	KeyOutLaunch     = "c_LAUNCH"
	KeyOutMap        = "c_MAP"
	KeyOutComing     = "c_COMING"
	KeyOutDisconnect = "c_DISCONNECT"
	KeyOutTP         = "c_TP"
	KeyOutMove       = "c_MOVE"

	// Keys received from the clients
	KeyInMove = "s_MOVE"
	KeyInTP   = "s_TP"
)

type Action struct {
	Key string
	Value interface{}
}

type Message struct {
	Key string
	Value string
}

// InitGame gets the character of the user from the storage, persists it into Redis
// and sends back data to this user.
func InitGame(client *Client){
	character := controller.FindCharacterByName(client.name)
	controller.PersistCharacterRedis(character)

	sendCharacterData(character, client)
	sendConnectedCharactersData(character, client)
}

// CloseGame persists the character of the connected user into the storage
// and deletes its entry from the redis.
func CloseGame(client *Client) {
	controller.PersistCharacterByName(client.name)
	controller.DeleteCharacterRedis(client.name)

	sendDisconnection(client)
}

// HandleUserAction parses a message coming from a user, and executes an action accordingly.
func HandleUserAction(message []byte, c *Client) {
	var action Action
	if err := json.Unmarshal(message, &action); err != nil {
		panic(err)
	}

	switch key := action.Key; key {
	case KeyInMove:
		HandleMoveAction(action, c)
	case KeyInTP:
		HandleTPAction(action, c)
	default:
		fmt.Println(fmt.Sprintf("Unknown action: %s", key))
	}
}

// HandleMoveAction is triggered when a character moves.
func HandleMoveAction(action Action, c *Client) {
	val, ok := action.Value.(map[string]interface{})
	if !ok {
		return
	}

	// Parse parameters from the user action
	user := parseUser(val)
	direction := parseDirection(val)
	key := controller.KeyByNameRedis(user)
	if key == nil {
		return
	}

	// Calculate the new coordinates of the character
	character := controller.FindCharacterRedis(*key)
	x, y := game.MoveCoordinates(character, direction)
	controller.UpdateCharacterRedis(character.Name, x, y, character.TileFormula, *key)

	// Notify the users that a character has moved
	message, err := message(KeyOutMove, val)
	if err != nil {
		fmt.Printf("error while building message: %s", err)
		return
	}

	sendMove(message, c)
}

// HandleTPAction is triggered when a character is teleported to another map.
func HandleTPAction(action Action, c *Client) {
	val, ok := action.Value.(map[string]interface{})
	if !ok {
		return
	}

	// Parse parameters from the user action
	user := parseUser(val)
	direction := parseDirection(val)
	key := controller.KeyByNameRedis(user)
	if key == nil {
		return
	}

	newMap := parseMap(val)
	if newMap == "" {
		return
	}

	// Get character from Redis
	character := controller.FindCharacterRedis(*key)
	controller.DeleteCharacterRedis(user)

	// Calculate new coordinates
	newGamemap := controller.FindGamemapByName(newMap)
	x, y := game.TPCoordinates(character, &newGamemap, direction)

	// Get character data from storage and update the Redis entry
	persistedCharacter := controller.FindCharacterByName(user)
	updateCharacter(&persistedCharacter, x, y, newGamemap)
	controller.PersistCharacterRedis(persistedCharacter)

	// Persist the updated character from Redis into Storage
	controller.PersistCharacterByName(user)

	sendTP(persistedCharacter, c)
}

// sendCharacterInfo sends the given character data to the connected users.
func sendCharacterData(character entity.Character, client *Client) {
	strChar, _ := json.Marshal(character)
	m, _ := json.Marshal(Message{Key: KeyOutLaunch, Value: string(strChar)})
	client.Unicast(m)

	m, _ = json.Marshal(Message{Key: KeyOutComing, Value: string(strChar)})
	client.Broadcast(m)
}

// sendConnectedCharactersInfo fetches other connected characters from the storage
// and sends this data back to the user.
func sendConnectedCharactersData(character entity.Character, client *Client) {
	connectedCharacters := controller.FindCharactersByMapRedis(character.Gamemap.Name)
	strConnectedChars, _ := json.Marshal(connectedCharacters)
	m, _ := json.Marshal(Message{Key: KeyOutMap, Value: string(strConnectedChars)})
	client.Unicast(m)
}

// sendDisconnection notifies everyone that a character has disconnected.
func sendDisconnection(client *Client) {
	m, _ := json.Marshal(Message{Key: KeyOutDisconnect, Value: client.name})
	client.Broadcast(m)
}

// sendTP notifies everyone that a character has teleported to a new map.
func sendTP(character entity.Character, client *Client) {
	strChar, _ := json.Marshal(character)

	m, _ := json.Marshal(Message{Key: KeyOutTP, Value: string(strChar)})
	client.Broadcast(m)

	m, _ = json.Marshal(Message{Key: KeyOutComing, Value: string(strChar)})
	client.Broadcast(m)

	sendConnectedCharactersData(character, client)
}

// sendMove notifies everyone that a character has moved.
func sendMove(message Message, client *Client) {
	m, _ := json.Marshal(message)
	client.Broadcast(m)
}

func updateCharacter(character *entity.Character, x int, y int, gamemap entity.Gamemap) {
	character.X = x
	character.Y = y
	character.Gamemap = gamemap
	character.GamemapID = gamemap.ID
}

func message(key string, values map[string]interface{}) (Message, error) {
	var message Message
	message.Key = key

	jsonString, err := json.Marshal(values)
	if err != nil {
		return message, err
	}
	message.Value = string(jsonString)

	return message, nil
}

func parseUser(val map[string]interface{}) string {
	return val["User"].(string)
}

func parseDirection(val map[string]interface{}) float64 {
	return val["Dir"].(float64)
}

func parseMap(val map[string]interface{}) string {
	return val["newMap"].(string)
}