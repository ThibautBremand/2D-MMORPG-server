package db

// gamemapRepository contains all the methods used to directly interact with the Storage.
// It allows to fetch Gamemaps.

import (
	"fmt"
	"server/entity"
	"server/utils"
	"strings"
)

// FindGamemapByName returns a Gamemap object from the db using the given name.
func FindGamemapByName(name string) entity.Gamemap {
	var gamemap entity.Gamemap
	DB.Where("name = ?", name).Find(&gamemap)

	return gamemap
}

// PersistNewGamemap create a new gamemap into the storage, using the data sent by the user.
func PersistNewGamemap(name string, jsonMap string) error {
	var count int64
	DB.Model(&entity.Character{}).Where("name = ?", strings.ToLower(name)).Count(&count)
	if count > 0 {
		return &utils.NameAlreadyTaken{Err: fmt.Errorf("the mapname is already taken")}
	}
	character := entity.Gamemap{
		Name: strings.ToLower(name),
		Raw:  jsonMap,
	}
	DB.Model(&entity.Character{}).Create(&character)
	return nil
}
