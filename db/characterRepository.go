package db

// characterRepository contains all the methods used to directly interact with the DBs (Storage and Redis).
// It allows to fetch, update, and persist Characters.

import (
	"fmt"
	"server/entity"
	"strconv"
)

//
// RDBS (PostGreSQL)
//

// FindCharacterByName returns a Character object from the db using the given name.
func FindCharacterByName(name string) entity.Character {
	var character entity.Character
	DB.Preload("Gamemap").Where("name = ?", name).Find(&character)

	return character
}

// PersistCharacterByName retrieves a character's data from the Redis, using the given name,
// and updates the corresponding Character entry in the DBs.
func PersistCharacterByName(name string) {
	keys, _ := ScanKeys("", fmt.Sprintf("-%s", name))
	if len(keys) != 1 {
		return
	}
	key := keys[0]

	character := FindCharacterRedis(key)
	var model entity.Character
	DB.Model(&model).Where("name = ?", name).Update(map[string]interface{}{"x": character.X, "y": character.Y, "tileFormula": character.TileFormula, "gamemap_id": character.GamemapID})
}

//
// Redis
//

// PersistCharacterRedis stores the given Character object into Redis.
func PersistCharacterRedis(character entity.Character) {
	key := fmt.Sprintf("%s-%s", character.Gamemap.Name, character.Name)
	values :=  map[string]interface{}{
		"x": strconv.Itoa(character.X),
		"y": strconv.Itoa(character.Y),
		"tileFormula": character.TileFormula,
		"name": character.Name,
		"gamemap": character.Gamemap.ID,
	}

	Redis.HMSet(key, values)
}

func UpdateCharacterRedis(key string, values map[string]interface{}) {
	Redis.HMSet(key, values)
}

func DeleteCharacterRedis(name string) {
	keys, _ := ScanKeys("", fmt.Sprintf("-%s", name))
	if len(keys) != 1 {
		return
	}
	Redis.Del(keys[0])
	fmt.Println(fmt.Sprintf("%s deleted from redis !", keys[0]))
}

// KeyByNameRedis returns a Redis key using the given character name.
func KeyByNameRedis(characterName string) *string {
	keys, _ := ScanKeys("", fmt.Sprintf("-%s", characterName))
	if len(keys) < 1 {
		return nil
	}

	return &keys[0]
}

// FindCharacterRedis returns a character from Redis using the given key
func FindCharacterRedis(key string) *entity.CharacterView {
	value, _ := Redis.HGetAll(key).Result()
	x, _ := strconv.Atoi(value["x"])
	y, _ := strconv.Atoi(value["y"])
	gamemapID, _ := strconv.Atoi(value["gamemap"])

	return &entity.CharacterView{
		Name: value["name"],
		X: x,
		Y: y,
		TileFormula: value["tileFormula"],
		GamemapID: gamemapID,
	}
}

// FindCharactersByMapRedis returns a Characters list from Redis using the given map name.
func FindCharactersByMapRedis(mapName string) []entity.CharacterView {
	keys, _ := ScanKeys(mapName, "")
	if len(keys) < 1 {
		return make([]entity.CharacterView, 0)
	}

	connectedCharacters := make([]entity.CharacterView, len(keys))
	for i, key := range keys {
		character := FindCharacterRedis(key)
		connectedCharacters[i] = *character
	}

	return connectedCharacters
}