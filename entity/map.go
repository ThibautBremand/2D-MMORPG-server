package entity

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Gamemap struct {
	gorm.Model
	Name string `json:"name"`
	Raw string `json:"raw"`
}

type Content struct {
	Height int `json:"height"`
	Width int `json:"width"`
}

const edgeMargin = 1

func (g *Gamemap) Height() int {
	var content Content
	err := json.Unmarshal([]byte(g.Raw), &content)
	if err != nil {
		fmt.Println(fmt.Sprintf("an error occured when trying to get Gamemap %s 's height! %s", g.ID, err))
		return 10	// arbitrary default value
	}

	return content.Height - edgeMargin
}

func (g *Gamemap) Width() int {
	var content Content
	//fmt.Println(g.Raw)
	err := json.Unmarshal([]byte(g.Raw), &content)
	if err != nil {
		fmt.Println(fmt.Sprintf("an error occured when trying to get Gamemap %s 's width! %s", g.ID, err))
		return 10	// arbitrary default value
	}

	fmt.Println(fmt.Sprintf("test %v", content))

	return content.Width - edgeMargin
}

func (g *Gamemap) EdgeMargin() int {
	return edgeMargin
}