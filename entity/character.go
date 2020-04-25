package entity

import "github.com/jinzhu/gorm"

type Character struct {
	gorm.Model
	Name string `json:"name"`
	X int `json:"x"`
	Y int `json:"y"`
	TileFormula string `json:"tileFormula"`
	GamemapID uint
	Gamemap Gamemap
}

type CharacterView struct {
	Name string `json:"name"`
	X int `json:"x"`
	Y int `json:"y"`
	TileFormula string `json:"tileFormula"`
	GamemapID int
}