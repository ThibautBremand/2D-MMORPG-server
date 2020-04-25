package db

// gamemapRepository contains all the methods used to directly interact with the Storage.
// It allows to fetch Gamemaps.

import "server/entity"

// FindGamemapByName returns a Gamemap object from the db using the given name.
func FindGamemapByName(name string) entity.Gamemap {
	var gamemap entity.Gamemap
	DB.Where("name = ?", name).Find(&gamemap)

	return gamemap
}
