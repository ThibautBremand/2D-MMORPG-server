package controller

import (
	"fmt"
	"server/db"
	"server/entity"
)

//
// Persisted Characters
//

func FindCharacterByName(name string) entity.Character {
	return db.FindCharacterByName(name)
}

func PersistCharacterByName(name string) {
	db.PersistCharacterByName(name)
}

//
// Cached Characters
//

func FindCharactersByMapRedis(mapName string) []entity.CharacterView {
	return db.FindCharactersByMapRedis(mapName)
}

func FindCharacterRedis(key string) *entity.CharacterView {
	return db.FindCharacterRedis(key)
}

func PersistCharacterRedis(character entity.Character) {
	db.PersistCharacterRedis(character)
	fmt.Println("Character " + character.Name + " stored in Redis")
}

func UpdateCharacterRedis(name string, x int, y int, tileFormula string, key string) {
	values :=  map[string]interface{}{
		"name": name,
		"x": x,
		"y": y,
		"tileFormula": tileFormula,
	}

	db.UpdateCharacterRedis(key, values)
}

func DeleteCharacterRedis(name string) {
	db.DeleteCharacterRedis(name)
}

func FindGamemapByName(name string) entity.Gamemap {
	return db.FindGamemapByName(name)
}

func KeyByNameRedis(characterName string) *string {
	return db.KeyByNameRedis(characterName)
}